package wrapper

import (
	"context"
	"database/sql"
)

// define the query func with context.
type QueryContextFunc func(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)

// define the exec func with context.
type ExecContextFunc func(ctx context.Context, query string, args ...interface{}) (sql.Result, error)

// define the database operation.
type Wrapper interface {
	WrapQueryContext(fn QueryContextFunc, sql string, args ...interface{}) QueryContextFunc
	WrapExecContext(fn ExecContextFunc, sql string, args ...interface{}) ExecContextFunc
}

// TracerOption define the wrapper's option
type TracerOption interface {
	QueryBuilder() func(query string, args ...interface{}) string
}
