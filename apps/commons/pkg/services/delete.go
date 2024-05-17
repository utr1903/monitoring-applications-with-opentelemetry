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

const deleteDbOperation = "DELETE"

type DeleteRequest struct {
}

type DeleteResult struct {
	Result  bool
	Message string
}

type IDeleteService interface {
	Delete(ctx context.Context, req *DeleteRequest) (*DeleteResult, error)
}

type DeleteService struct {
	logger        loggers.ILogger
	entityService data.IEntityService
	sqlEnricher   *otelsql.SqlEnricher
}

func NewDeleteService(logger loggers.ILogger, createDbNotReachableError bool) *DeleteService {
	sqlEnricher := otelsql.NewSqlEnricher(
		otelsql.WithDbSystem("postgres"),
		otelsql.WithServer("postgres.default.svc.cluster.local"),
		otelsql.WithDatabase("mydb"),
		otelsql.WithTable("tasks"),
		otelsql.WithPort(5432),
	)

	return &DeleteService{
		logger:        logger,
		entityService: postgres.NewDatabase(createDbNotReachableError),
		sqlEnricher:   sqlEnricher,
	}
}

func (s *DeleteService) Delete(ctx context.Context, req *DeleteRequest) (*DeleteResult, error) {
	dbStatement := deleteDbOperation + " FROM tasks"
	ctx, dbSpan := s.sqlEnricher.CreateSpan(ctx, deleteDbOperation, dbStatement)
	defer dbSpan.End()

	res := s.entityService.Delete(ctx, dbStatement)
	if !res.Success {
		err := errors.New(res.Message)

		msg := "Deleting task failed."
		dbSpan.SetStatus(codes.Error, msg)
		dbSpan.RecordError(err)

		s.logger.Log(ctx, loggers.Error, msg, map[string]interface{}{
			"error.message": res.Message,
		})
		return nil, err
	}

	return &DeleteResult{
		Message: res.Message,
	}, nil
}
