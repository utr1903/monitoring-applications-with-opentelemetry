package server

import "github.com/utr1903/monitoring-applications-with-opentelemetry/apps/commons/pkg/data"

type StoreRequest struct {
	Task string
}

type StoreResult struct {
	Result  bool
	Message string
	Body    *data.Task
}

type IStoreService interface {
	Store(*StoreRequest) *StoreResult
}

type StoreService struct {
	entityService data.IEntityService
}

func NewStoreService() *StoreService {
	return &StoreService{
		entityService: data.NewUnimplementedEntityService(),
	}
}

func (s *StoreService) Store(req *StoreRequest) *StoreResult {
	res := s.entityService.Create(req.Task)

	return &StoreResult{
		Result:  res.Result,
		Message: res.Message,
		Body:    res.Body,
	}
}
