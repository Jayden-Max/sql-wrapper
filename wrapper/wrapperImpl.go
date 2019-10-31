package wrapper

import (
	"context"
	"database/sql"
)

// define the query func with context.
type QueryContextFunc func(ctx context.Context, sql string, args ...interface{}) (*sql.Rows, error)

// define the exec func with context.
type ExecContextFunc func(ctx context.Context, sql string, args ...interface{}) (sql.Result, error)

// define the database operation.
type Wrapper interface {
	WrapQueryContext(fn QueryContextFunc) QueryContextFunc
	WrapExecContext(fn ExecContextFunc) ExecContextFunc
}

// TracerOption define the wrapper's option
type TracerOption interface {
	QueryBuilder() func(query string, args ...interface{}) string
}
