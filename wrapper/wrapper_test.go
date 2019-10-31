package wrapper

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/mocktracer"
	"reflect"
	"testing"
)

func init() {
	opentracing.SetGlobalTracer(mocktracer.New())
}

func TestTracerWrapper_WrapExecContext(t *testing.T) {
	type fields struct {
		tracer *tracer
	}
	type args struct {
		fn    ExecContextFunc
		query string
		args  []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   ExecContextFunc
	}{
		// TODO: Add test cases.
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t1 *testing.T) {
			t := &TracerWrapper{
				tracer: tt.fields.tracer,
			}
			if got := t.WrapExecContext(tt.args.fn, tt.args.query, tt.args.args...); !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("WrapExecContext() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTracerWrapper_WrapQueryContext(t *testing.T) {
	type args struct {
		ctx   context.Context
		fn    QueryContextFunc
		query string
		args  []interface{}
	}
	tests := []struct {
		name   string
		fields *TracerWrapper
		args   args
		want   string
	}{
		// TODO: Add test cases.
		{
			name: "TestTracerWrapper_WrapQueryContext_MySQL",
			args: args{
				ctx: context.TODO(),
				fn: QueryContextFunc(func(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
					db, _, err := sqlmock.New()
					if err != nil {
						t.Errorf("mock sql conn failed:%v", err.Error())
					}
					return db.QueryContext(ctx, query, args...)
				}),
				query: "SELECT a FROM b WHERE c = ?",
				args:  []interface{}{"d"},
			},
			fields: NewMySQLTracerWrapper(),
			want:   "SELECT ... FROM b WHERE c = ?",
		},
		{
			name: "TestTracerWrapper_WrapQueryContext_MySQL_Customized",
			args: args{
				ctx: context.TODO(),
				fn: QueryContextFunc(func(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
					db, _, err := sqlmock.New()
					if err != nil {
						t.Errorf("mock sql conn failed:%v", err.Error())
					}
					return db.QueryContext(ctx, query, args...)
				}),
				query: "SELECT a FROM b WHERE c = ?",
				args:  []interface{}{"d"},
			},
			fields: NewMySQLTracerWrapperWithOpts(RawQueryOption),
			want:   "SELECT a FROM b WHERE c = d",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields.WrapQueryContext(tt.args.fn, tt.args.query, tt.args.args...)(tt.args.ctx, tt.args.query, tt.args.args...)
			dt := tt.fields.tracer
			if ins := dt.span.(*mocktracer.MockSpan).Tag(string(ext.DBInstance)); ins != dt.dbInstance {
				t.Errorf("ext.DBInstance = %v, want %v", ins, dt.dbInstance)
			}

			if st := dt.span.(*mocktracer.MockSpan).Tag(string(ext.DBStatement)); st != tt.want {
				t.Errorf("ext.DBStatement = %v, want %v", st, tt.want)
			}

			if tp := dt.span.(*mocktracer.MockSpan).Tag(string(ext.DBType)); tp != dt.dbType {
				t.Errorf("ext.DBType = %v, want %v", tp, dt.dbType)
			}
		})
	}
}
