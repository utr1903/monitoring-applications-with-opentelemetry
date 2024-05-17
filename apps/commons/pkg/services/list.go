package server

import (
	"context"
	"errors"

	"github.com/utr1903/monitoring-applications-with-opentelemetry/apps/commons/pkg/data"
	"github.com/utr1903/monitoring-applications-with-opentelemetry/apps/commons/pkg/data/postgres"
	"github.com/utr1903/monitoring-applications-with-opentelemetry/apps/commons/pkg/loggers"
	otelsql "github.com/utr1903/monitoring-applications-with-opentelemetry/apps/commons/pkg/opentelemetry/sql"
	"go.opentelemetry.io/otel/codes"
)

const listDbOperation = "SELECT"

type ListRequest struct {
}

type ListResult struct {
	Result  bool
	Message string
	Body    []data.Task
}

type IListService interface {
	List(ctx context.Context, req *ListRequest) (*ListResult, error)
}

type ListService struct {
	logger        loggers.ILogger
	entityService data.IEntityService
	sqlEnricher   *otelsql.SqlEnricher
}

func NewListService(logger loggers.ILogger, createDbNotReachableError bool) *ListService {
	sqlEnricher := otelsql.NewSqlEnricher(
		otelsql.WithDbSystem("postgres"),
		otelsql.WithServer("postgres.default.svc.cluster.local"),
		otelsql.WithDatabase("mydb"),
		otelsql.WithTable("tasks"),
		otelsql.WithPort(5432),
	)

	return &ListService{
		logger:        logger,
		entityService: postgres.NewDatabase(createDbNotReachableError),
		sqlEnricher:   sqlEnricher,
	}
}

func (s *ListService) List(ctx context.Context, req *ListRequest) (*ListResult, error) {

	dbStatement := listDbOperation + " id, message FROM tasks"
	ctx, dbSpan := s.sqlEnricher.CreateSpan(ctx, listDbOperation, dbStatement)
	defer dbSpan.End()

	res := s.entityService.List(ctx, dbStatement)
	if !res.Success {
		err := errors.New(res.Message)

		msg := "Listing tasks failed."
		dbSpan.SetStatus(codes.Error, msg)
		dbSpan.RecordError(err)

		s.logger.Log(ctx, loggers.Error, msg, map[string]interface{}{
			"error.message": res.Message,
		})
		return nil, err
	}

	return &ListResult{
		Message: res.Message,
		Body:    res.Body,
	}, nil
}
