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

const createDbOperation = "INSERT"

type StoreRequest struct {
	Task string
}

type StoreResult struct {
	Message string
	Body    *data.Task
}

type IStoreService interface {
	Store(ctx context.Context, req *StoreRequest) (*StoreResult, error)
}

type StoreService struct {
	logger        loggers.ILogger
	entityService data.IEntityService
	sqlEnricher   *otelsql.SqlEnricher
}

func NewStoreService(logger loggers.ILogger, createDbNotReachableError bool) *StoreService {
	sqlEnricher := otelsql.NewSqlEnricher(
		otelsql.WithDbSystem("postgres"),
		otelsql.WithServer("postgres.default.svc.cluster.local"),
		otelsql.WithDatabase("mydb"),
		otelsql.WithTable("tasks"),
		otelsql.WithPort(5432),
	)

	return &StoreService{
		logger:        logger,
		entityService: postgres.NewDatabase(createDbNotReachableError),
		sqlEnricher:   sqlEnricher,
	}
}

func (s *StoreService) Store(ctx context.Context, req *StoreRequest) (*StoreResult, error) {

	dbStatement := createDbOperation + " INTO tasks (id, message) VALUES (?, ?)"
	ctx, dbSpan := s.sqlEnricher.CreateSpan(ctx, createDbOperation, dbStatement)
	defer dbSpan.End()

	res := s.entityService.Create(ctx, dbStatement)
	if !res.Success {
		err := errors.New(res.Message)

		msg := "Storing task failed."
		dbSpan.SetStatus(codes.Error, msg)
		dbSpan.RecordError(err)

		s.logger.Log(ctx, loggers.Error, msg, map[string]interface{}{
			"error.message": res.Message,
		})
		return nil, err
	}

	return &StoreResult{
		Message: res.Message,
		Body:    res.Body,
	}, nil
}
