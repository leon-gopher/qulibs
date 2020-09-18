package gorm

import (
	"fmt"
	"strings"

	gogorm "github.com/jinzhu/gorm"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

func (c *Client) RegisterTraceCallbacks() {
	c.traceOnce.Do(func() {
		// for create
		c.Callback().Create().Before("gorm:before_create").Register("pedestal:trace_before_create", func(scope *gogorm.Scope) {
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
			ext.DBInstance.Set(spanner, scope.TableName())

			scope.InstanceSet(OpentracingSpanContextKey, spanner)
		})
		c.Callback().Create().After("gorm:after_create").Register("pedestal:trace_after_create", func(scope *gogorm.Scope) {
			iface, ok := scope.InstanceGet(OpentracingSpanContextKey)
			if !ok {
				return
			}

			spanner, ok := iface.(opentracing.Span)
			if !ok {
				return
			}
			defer spanner.Finish()

			if scope.HasError() {
				ext.DBStatement.Set(spanner, fmt.Sprintf("%v", scope.DB().QueryExpr()))
				ext.Error.Set(spanner, true)
				spanner.LogKV("event", "error", "message", scope.DB().Error.Error())
			} else {
				ext.DBStatement.Set(spanner, scope.SQL)
			}
		})

		// for update
		c.Callback().Update().Before("gorm:before_update").Register("pedestal:trace_before_update", func(scope *gogorm.Scope) {
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
			ext.DBInstance.Set(spanner, scope.TableName())

			scope.InstanceSet(OpentracingSpanContextKey, spanner)
		})
		c.Callback().Update().After("gorm:after_update").Register("pedestal:trace_after_update", func(scope *gogorm.Scope) {
			iface, ok := scope.InstanceGet(OpentracingSpanContextKey)
			if !ok {
				return
			}

			spanner, ok := iface.(opentracing.Span)
			if !ok {
				return
			}
			defer spanner.Finish()

			if scope.HasError() {
				ext.DBStatement.Set(spanner, fmt.Sprintf("%v", scope.DB().QueryExpr()))
				ext.Error.Set(spanner, true)

				spanner.LogKV("event", "error", "message", scope.DB().Error.Error())
			} else {
				ext.DBStatement.Set(spanner, scope.SQL)
			}
		})

		// for query
		c.Callback().Query().Before("gorm:query").Register("pedestal:trace_before_query", func(scope *gogorm.Scope) {
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
			ext.DBInstance.Set(spanner, scope.TableName())

			scope.InstanceSet(OpentracingSpanContextKey, spanner)
		})
		c.Callback().Query().After("gorm:after_query").Register("pedestal:trace_after_query", func(scope *gogorm.Scope) {
			iface, ok := scope.InstanceGet(OpentracingSpanContextKey)
			if !ok {
				return
			}

			spanner, ok := iface.(opentracing.Span)
			if !ok {
				return
			}
			defer spanner.Finish()

			if scope.HasError() && (c.config.TraceIncludeNotFound || scope.DB().Error != gogorm.ErrRecordNotFound) {
				ext.DBStatement.Set(spanner, fmt.Sprintf("%v", scope.DB().QueryExpr()))
				ext.Error.Set(spanner, true)

				spanner.LogKV("event", "error", "message", scope.DB().Error.Error())
			} else {
				ext.DBStatement.Set(spanner, scope.SQL)
			}
		})

		// for delete
		c.Callback().Delete().Before("gorm:before_delete").Register("pedestal:trace_before_delete", func(scope *gogorm.Scope) {
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
			ext.DBInstance.Set(spanner, scope.TableName())

			scope.InstanceSet(OpentracingSpanContextKey, spanner)
		})
		c.Callback().Delete().After("gorm:after_delete").Register("pedestal:trace_after_delete", func(scope *gogorm.Scope) {
			iface, ok := scope.InstanceGet(OpentracingSpanContextKey)
			if !ok {
				return
			}

			spanner, ok := iface.(opentracing.Span)
			if !ok {
				return
			}
			defer spanner.Finish()

			if scope.HasError() {
				ext.DBStatement.Set(spanner, fmt.Sprintf("%v", scope.DB().QueryExpr()))
				ext.Error.Set(spanner, true)
				spanner.LogKV("event", "error", "message", scope.DB().Error.Error())
			} else {
				ext.DBStatement.Set(spanner, scope.SQL)
			}
		})

		// for raw
		c.Callback().RowQuery().Before("publish:update_table_nam").Register("pedestal:trace_before_sql", func(scope *gogorm.Scope) {
			iface, ok := scope.Get(OpentracingContextKey)
			if !ok {
				return
			}

			ctx, ok := iface.(TraceContext)
			if !ok {
				return
			}

			name := "SQL"
			switch strings.SplitN(scope.SQL, " ", 2)[0] {
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
			ext.DBInstance.Set(spanner, scope.TableName())

			scope.InstanceSet(OpentracingSpanContextKey, spanner)
		})
		c.Callback().RowQuery().After("gorm:row_query").Register("pedestal:trace_after_sql", func(scope *gogorm.Scope) {
			iface, ok := scope.InstanceGet(OpentracingSpanContextKey)
			if !ok {
				return
			}

			spanner, ok := iface.(opentracing.Span)
			if !ok {
				return
			}
			defer spanner.Finish()

			if scope.HasError() {
				ext.DBStatement.Set(spanner, scope.SQL)
				ext.Error.Set(spanner, true)
				spanner.LogKV("event", "error", "message", scope.DB().Error.Error())
			} else {
				ext.DBStatement.Set(spanner, scope.SQL)
			}
		})
	})
}
