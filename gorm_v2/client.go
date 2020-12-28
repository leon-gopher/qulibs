package gorm

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/leon-gopher/qulibs"
	"github.com/opentracing/opentracing-go"
	"gorm.io/driver/mysql"
	gogorm "gorm.io/gorm"
)

// A Client wrap *gorm.DB with best practices for development
type Client struct {
	*gogorm.DB

	mux       sync.RWMutex
	log       qulibs.Logger
	config    *Config
	traceOnce sync.Once
}

// New creates mysql client with config given and a dummy logger.
func New(config *Config) (*Client, error) {
	return NewWithLogger(config, qulibs.NewDummyLogger())
}

// NewWithLogger creates mysql client with config and logger given.
func NewWithLogger(config *Config, log qulibs.Logger) (client *Client, err error) {
	config.FillWithDefaults()

	mycfg, err := config.NewMycfg()
	if err != nil {
		return
	}
	fmt.Println(mysql.Open(mycfg.FormatDSN()))

	db, err := gogorm.Open(mysql.Open(mycfg.FormatDSN()), &gogorm.Config{})
	if err != nil {
		return
	}

	sqlDb, err := db.DB()
	if err != nil {
		return
	}

	if config.MaxOpenConns > 0 {
		sqlDb.SetMaxOpenConns(config.MaxOpenConns)
	}
	if config.MaxIdleConns > 0 {
		sqlDb.SetMaxIdleConns(config.MaxIdleConns)
	}
	if config.MaxLifeConns > 0 {
		sqlDb.SetConnMaxLifetime(time.Duration(config.MaxLifeConns) * time.Second)
	}

	err = sqlDb.Ping()
	if err != nil {
		return
	}

	if config.DebugSQL {
		db = db.Debug()
	}

	config.mycfg = mycfg

	client = &Client{
		DB:     db,
		log:    log,
		config: config,
	}
	return
}

// Select switches to a new database of dbname given by creating a new gorm instance.
func (c *Client) Select(dbname string) (client *Client, err error) {
	c.mux.RLock()
	if c.config.IsEqualDB(dbname) {
		c.mux.RUnlock()

		return c, nil
	}

	config, err := c.config.NewWithDB(dbname)
	if err != nil {
		c.mux.RUnlock()
		return
	}

	name := config.Name()

	// first, try loading a client from default manager
	client, err = DefaultMgr.NewClientWithLogger(name, c.log)
	if err == nil {
		c.mux.RUnlock()

		return client, nil
	}

	c.mux.RUnlock()

	// second, register new client for default manager
	c.mux.Lock()
	defer c.mux.Unlock()

	DefaultMgr.Add(name, config)

	return DefaultMgr.NewClientWithLogger(name, c.log)
}

func (c *Client) Trace(ctx context.Context, tracers ...opentracing.Tracer) *TraceClient {
	if ctx == nil {
		return c.TraceWithSpanContext(nil, tracers...)
	}

	span := opentracing.SpanFromContext(ctx)
	if span != nil {
		return c.TraceWithSpanContext(span.Context(), tracers...)
	}
	return c.TraceWithSpanContext(nil, tracers...)
}

func (c *Client) TraceWithSpanContext(ctx opentracing.SpanContext, tracers ...opentracing.Tracer) *TraceClient {

	var tracer opentracing.Tracer
	if len(tracers) > 0 {
		tracer = tracers[0]
	} else {
		tracer = opentracing.GlobalTracer()
	}

	trace := TraceContext{
		tracer:  tracer,
		spanCtx: ctx,
	}

	return &TraceClient{
		DB:    c.Set(OpentracingContextKey, trace),
		trace: &trace,
	}
}
