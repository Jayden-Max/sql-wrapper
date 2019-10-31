package wrapper

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"regexp"
	"strings"
)

var (
	// RawQueryOption is the raw query option.
	// rawQueryOption convert the '?' place to real data.
	// Ex: "SELECT a FROM b WHERE c = ?" will be "SELECT a FROM b WHERE c = d ".
	RawQueryOption = rawQueryOption{}
	// IgnoreSelectColumnsOption enable the ignore select columns option,
	// Ex: "SELECT A,B FROM C WHERE D = ?" will be "SELECT ... FROM C WHERE D = ?".
	IgnoreSelectColumnsOption = ignoreSelectColumnsOption{}
)

// https://github.com/opentracing-contrib/opentracing-specification-zh/blob/master/semantic_conventions.md.
type tracer struct {
	// 数据库实例名称
	// 以Java为例，如果 jdbc.url="jdbc:mysql://127.0.0.1:3306/customers"，实例名为 "customers".
	dbInstance string

	// 一个针对给定数据库类型的数据库访问语句
	// 例如， 针对数据库类型 db.type="sql"，语句可能是 "SELECT * FROM wuser_table"; 针对数据库类型为
	// db.type="redis"，语句可能是 "SET mykey 'WuValue'".
	dbStatement string

	// 数据库类型。对于任何支持SQL的数据库，取值为 "sql". 否则，使用小写的数据类型名称，如 "cassandra", "hbase", or "redis".
	dbType string

	// 访问数据库的用户名。如 "readonly_user" 或 "reporting_user"
	dbUser string

	span          opentracing.Span
	queryBuilders []func(query string, args ...interface{}) string
}

func newTracer(dbType string, options ...TracerOption) *tracer {
	t := &tracer{
		dbType: dbType,
	}

	for _, op := range options {
		t.addQueryBuilder(op.QueryBuilder())
	}

	return t
}

// Do obtain openTracing's global tracer and add span tags.
// The tags follow https://github.com/opentracing/specification/blob/master/semantic_conventions.md.
func (t *tracer) do(ctx context.Context) {
	tracer := opentracing.GlobalTracer()
	span := opentracing.SpanFromContext(ctx)
	if span == nil {
		span = tracer.StartSpan(t.dbType)
	} else {
		span = tracer.StartSpan(t.dbType, opentracing.ChildOf(span.Context()))
	}
	// span set tag
	ext.DBInstance.Set(span, t.dbInstance)
	ext.DBStatement.Set(span, t.dbStatement)
	ext.DBType.Set(span, t.dbType)
	ext.DBUser.Set(span, t.dbUser)
	ctx = opentracing.ContextWithSpan(ctx, span)
	t.span = span
}

// Add queryBuilder to tracer.
func (t *tracer) addQueryBuilder(fn func(query string, args ...interface{}) string) {
	t.queryBuilders = append(t.queryBuilders, fn)
}

// Close span.
func (t *tracer) close() {
	if t.span != nil {
		t.span.Finish()
	}
}

// ***************** TracerOption *********************
// implement TracerOption
// ***************** TracerOption *********************

type ignoreSelectColumnsOption struct{}

func (opt ignoreSelectColumnsOption) QueryBuilder() func(query string, args ...interface{}) string {
	return ignoreSelectColumnQueryBuilder
}

func ignoreSelectColumnQueryBuilder(query string, args ...interface{}) string {
	query = strings.Replace(query, "select", "SELECT", -1)
	query = strings.Replace(query, "from", "FROM", -1)

	// 这里的正则不懂？？？
	r := regexp.MustCompile("(?s)SELECT (.*) FROM")
	return r.ReplaceAllString(query, "SELECT ... FROM")
}

type rawQueryOption struct{}

func (rq rawQueryOption) QueryBuilder() func(query string, args ...interface{}) string {
	return rawQueryBuilder
}

func rawQueryBuilder(query string, args ...interface{}) string {
	query = strings.Replace(query, "?", "%v", -1)
	return fmt.Sprintf(query, args...)
}

// ******************* TracerWrapper *******************
// implement Wrapper
// ******************* TracerWrapper *******************

func NewTracerWrapper(dbType string) *TracerWrapper {
	return newTracerWrapper(newTracer(dbType, IgnoreSelectColumnsOption))
}

type TracerWrapper struct {
	tracer *tracer
}

func newTracerWrapper(t *tracer) *TracerWrapper {
	return &TracerWrapper{tracer: t}
}

// For sql's select ...
func (t *TracerWrapper) WrapQueryContext(fn QueryContextFunc) QueryContextFunc {
	tracerFn := func(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
		t.tracer.dbStatement = t.QueryBuilder(query, args...)
		t.tracer.do(ctx)
		defer t.tracer.close()
		return fn(ctx, query, args...)
	}

	return tracerFn
}

// For sql's update、insert、delete ...
func (t *TracerWrapper) WrapExecContext(fn ExecContextFunc) ExecContextFunc {
	tracerFn := func(ctx context.Context, sql string, args ...interface{}) (sql.Result, error) {
		t.tracer.dbStatement = t.QueryBuilder(sql, args...)
		t.tracer.do(ctx)
		defer t.tracer.close()
		return fn(ctx, sql, args...)
	}

	return tracerFn
}

func (t *TracerWrapper) QueryBuilder(query string, args ...interface{}) string {
	for _, fn := range t.tracer.queryBuilders {
		query = fn(query, args...)
	}
	return query
}
