package gorm

import (
	"fmt"
	"strings"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	gogorm "gorm.io/gorm"
)

func (c *Client) RegisterTraceCallbacks() {
	c.traceOnce.Do(func() {
		// for create
		c.Callback().Create().Before("gorm:before_create").Register("ggs:trace_before_create", func(scope *gogorm.DB) {
			iface, ok := scope.Get(OpentracingContextKey)
			if !ok {
				return
			}
			ctx, ok := iface.(TraceContext)
			if !ok {
				return
			}
			spanner := ctx.StartSpan("CREATE")
			ext.SpanKindRPCClient.Set(spanner)
			ext.PeerService.Set(spanner, "mysql")
			ext.PeerHostname.Set(spanner, c.config.mycfg.Addr)
			ext.DBType.Set(spanner, "mysql")
			ext.DBUser.Set(spanner, c.config.mycfg.User)
			ext.DBInstance.Set(spanner, scope.Statement.Table)

			scope.InstanceSet(OpentracingSpanContextKey, spanner)
		})
		c.Callback().Create().After("gorm:after_create").Register("ggs:trace_after_create", func(scope *gogorm.DB) {
			iface, ok := scope.InstanceGet(OpentracingSpanContextKey)
			if !ok {
				return
			}

			spanner, ok := iface.(opentracing.Span)
			if !ok {
				return
			}
			defer spanner.Finish()

			if scope.Error != nil {
				ext.DBStatement.Set(spanner, fmt.Sprintf("%v", scope.Statement.SQL))
				ext.Error.Set(spanner, true)
				spanner.LogKV("event", "error", "message", scope.Error.Error())
			} else {
				ext.DBStatement.Set(spanner, scope.Statement.SQL.String())
			}
		})

		// for update
		c.Callback().Update().Before("gorm:before_update").Register("ggs:trace_before_update", func(scope *gogorm.DB) {
			iface, ok := scope.Get(OpentracingContextKey)
			if !ok {
				return
			}

			ctx, ok := iface.(TraceContext)
			if !ok {
				return
			}

			spanner := ctx.StartSpan("UPDATE")
			ext.SpanKindRPCClient.Set(spanner)
			ext.PeerService.Set(spanner, "mysql")
			ext.PeerHostname.Set(spanner, c.config.mycfg.Addr)
			ext.DBType.Set(spanner, "mysql")
			ext.DBUser.Set(spanner, c.config.mycfg.User)
			ext.DBInstance.Set(spanner, scope.Statement.Table)

			scope.InstanceSet(OpentracingSpanContextKey, spanner)
		})
		c.Callback().Update().After("gorm:after_update").Register("ggs:trace_after_update", func(scope *gogorm.DB) {
			iface, ok := scope.InstanceGet(OpentracingSpanContextKey)
			if !ok {
				return
			}

			spanner, ok := iface.(opentracing.Span)
			if !ok {
				return
			}
			defer spanner.Finish()

			if scope.Error != nil {
				ext.DBStatement.Set(spanner, fmt.Sprintf("%v", scope.Statement.SQL))
				ext.Error.Set(spanner, true)

				spanner.LogKV("event", "error", "message", scope.Error.Error())
			} else {
				ext.DBStatement.Set(spanner, scope.Statement.SQL.String())
			}
		})

		// for query
		c.Callback().Query().Before("gorm:query").Register("ggs:trace_before_query", func(scope *gogorm.DB) {
			iface, ok := scope.Get(OpentracingContextKey)
			if !ok {
				return
			}

			ctx, ok := iface.(TraceContext)
			if !ok {
				return
			}

			spanner := ctx.StartSpan("QUERY")
			ext.SpanKindRPCClient.Set(spanner)
			ext.PeerService.Set(spanner, "mysql")
			ext.PeerHostname.Set(spanner, c.config.mycfg.Addr)
			ext.DBType.Set(spanner, "mysql")
			ext.DBUser.Set(spanner, c.config.mycfg.User)
			ext.DBInstance.Set(spanner, scope.Statement.SQL.String())

			scope.InstanceSet(OpentracingSpanContextKey, spanner)
		})
		c.Callback().Query().After("gorm:after_query").Register("ggs:trace_after_query", func(scope *gogorm.DB) {
			iface, ok := scope.InstanceGet(OpentracingSpanContextKey)
			if !ok {
				return
			}

			spanner, ok := iface.(opentracing.Span)
			if !ok {
				return
			}
			defer spanner.Finish()

			if scope.Error != nil && (c.config.TraceIncludeNotFound || scope.Error != gogorm.ErrRecordNotFound) {
				ext.DBStatement.Set(spanner, fmt.Sprintf("%v", scope.Statement.SQL))
				ext.Error.Set(spanner, true)

				spanner.LogKV("event", "error", "message", scope.Error.Error())
			} else {
				ext.DBStatement.Set(spanner, scope.Statement.SQL.String())
			}
		})

		// for delete
		c.Callback().Delete().Before("gorm:before_delete").Register("ggs:trace_before_delete", func(scope *gogorm.DB) {
			iface, ok := scope.Get(OpentracingContextKey)
			if !ok {
				return
			}

			ctx, ok := iface.(TraceContext)
			if !ok {
				return
			}

			spanner := ctx.StartSpan("DELETE")
			ext.SpanKindRPCClient.Set(spanner)
			ext.PeerService.Set(spanner, "mysql")
			ext.PeerHostname.Set(spanner, c.config.mycfg.Addr)
			ext.DBType.Set(spanner, "mysql")
			ext.DBUser.Set(spanner, c.config.mycfg.User)
			ext.DBInstance.Set(spanner, scope.Statement.Table)

			scope.InstanceSet(OpentracingSpanContextKey, spanner)
		})
		c.Callback().Delete().After("gorm:after_delete").Register("ggs:trace_after_delete", func(scope *gogorm.DB) {
			iface, ok := scope.InstanceGet(OpentracingSpanContextKey)
			if !ok {
				return
			}

			spanner, ok := iface.(opentracing.Span)
			if !ok {
				return
			}
			defer spanner.Finish()

			if scope.Error != nil {
				ext.DBStatement.Set(spanner, fmt.Sprintf("%v", scope.Statement.SQL))
				ext.Error.Set(spanner, true)
				spanner.LogKV("event", "error", "message", scope.Error.Error())
			} else {
				ext.DBStatement.Set(spanner, scope.Statement.SQL.String())
			}
		})

		// for raw
		c.Callback().Query().Before("publish:update_table_nam").Register("ggs:trace_before_sql", func(scope *gogorm.DB) {
			iface, ok := scope.Get(OpentracingContextKey)
			if !ok {
				return
			}

			ctx, ok := iface.(TraceContext)
			if !ok {
				return
			}

			name := "SQL"
			switch strings.SplitN(scope.Statement.SQL.String(), " ", 2)[0] {
			case "select", "SELECT", "Select":
				name = "QUERY"
			case "insert", "INSERT", "Insert", "replace", "REPLACE", "Replace":
				name = "CREATE"
			case "update", "UPDATE", "Update":
				name = "UPDATE"
			case "delete", "DELETE", "Delete":
				name = "DELETE"
			}

			name = "mysql_" + name

			spanner := ctx.StartSpan(name)
			ext.SpanKindRPCClient.Set(spanner)
			ext.PeerService.Set(spanner, "mysql")
			ext.PeerHostname.Set(spanner, c.config.mycfg.Addr)
			ext.DBType.Set(spanner, "mysql")
			ext.DBUser.Set(spanner, c.config.mycfg.User)
			ext.DBInstance.Set(spanner, scope.Statement.Table)

			scope.InstanceSet(OpentracingSpanContextKey, spanner)
		})
		c.Callback().Query().After("gorm:row_query").Register("ggs:trace_after_sql", func(scope *gogorm.DB) {
			iface, ok := scope.InstanceGet(OpentracingSpanContextKey)
			if !ok {
				return
			}

			spanner, ok := iface.(opentracing.Span)
			if !ok {
				return
			}
			defer spanner.Finish()

			if scope.Error != nil {
				ext.DBStatement.Set(spanner, scope.Statement.SQL.String())
				ext.Error.Set(spanner, true)
				spanner.LogKV("event", "error", "message", scope.Error.Error())
			} else {
				ext.DBStatement.Set(spanner, scope.Statement.SQL.String())
			}
		})
	})
}
