package sql

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type sqlOpts struct {
	DbSystem string
	Server   string
	Port     int
	Username string
	Database string
	Table    string
}

type sqlOptFunc func(*sqlOpts)

type SqlEnricher struct {
	Opts *sqlOpts
}

// Create a SQL database instance
func NewSqlEnricher(
	optFuncs ...sqlOptFunc,
) *SqlEnricher {

	// Apply external options
	opts := &sqlOpts{
		DbSystem: "",
		Server:   "",
		Port:     0,
		Username: "",
		Database: "",
		Table:    "",
	}
	for _, f := range optFuncs {
		f(opts)
	}

	return &SqlEnricher{
		Opts: opts,
	}
}

// Configure DB system
func WithDbSystem(dbSystem string) sqlOptFunc {
	return func(opts *sqlOpts) {
		opts.DbSystem = dbSystem
	}
}

// Configure SQL server
func WithServer(server string) sqlOptFunc {
	return func(opts *sqlOpts) {
		opts.Server = server
	}
}

// Configure SQL port
func WithPort(port int) sqlOptFunc {
	return func(opts *sqlOpts) {
		opts.Port = port
	}
}

// Configure SQL username
func WithUsername(username string) sqlOptFunc {
	return func(opts *sqlOpts) {
		opts.Username = username
	}
}

// Configure SQL database
func WithDatabase(database string) sqlOptFunc {
	return func(opts *sqlOpts) {
		opts.Database = database
	}
}

// Configure SQL table
func WithTable(table string) sqlOptFunc {
	return func(opts *sqlOpts) {
		opts.Table = table
	}
}

func (e *SqlEnricher) CreateSpan(
	ctx context.Context,
	operation string,
	statement string,
) (
	context.Context,
	trace.Span,
) {
	// Create database span
	parentSpan := trace.SpanFromContext(ctx)
	ctx, dbSpan := parentSpan.TracerProvider().
		Tracer("sql-enricher").
		Start(
			ctx,
			operation+" "+e.Opts.Database+"."+e.Opts.Table,
			trace.WithSpanKind(trace.SpanKindClient),
		)

	// Set additional span attributes
	dbSpanAttrs := e.getCommonAttributes()
	dbSpanAttrs = append(dbSpanAttrs, attribute.Key("db.operation").String(operation))
	dbSpanAttrs = append(dbSpanAttrs, attribute.Key("db.statement").String(statement))
	dbSpan.SetAttributes(dbSpanAttrs...)

	return ctx, dbSpan
}

func (e *SqlEnricher) getCommonAttributes() []attribute.KeyValue {
	return []attribute.KeyValue{
		attribute.Key("server.address").String(e.Opts.Server),
		attribute.Key("server.port").Int(e.Opts.Port),
		attribute.Key("db.system").String(e.Opts.DbSystem),
		attribute.Key("db.user").String(e.Opts.Username),
		attribute.Key("db.name").String(e.Opts.Database),
		attribute.Key("db.table").String(e.Opts.Table),
	}
}
