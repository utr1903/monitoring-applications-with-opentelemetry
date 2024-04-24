package server

import "github.com/utr1903/monitoring-applications-with-opentelemetry/apps/commons/pkg/data"

type ListRequest struct {
}

type ListResult struct {
	Result  bool
	Message string
	Body    []data.Task
}

type IListService interface {
	List(*ListRequest) *ListResult
}

type ListService struct {
	entityService data.IEntityService
}

func NewListService() *ListService {
	return &ListService{
		entityService: data.NewUnimplementedEntityService(),
	}
}

func (s *ListService) List(req *ListRequest) *ListResult {
	res := s.entityService.List()

	return &ListResult{
		Result:  res.Result,
		Message: res.Message,
		Body:    res.Body,
	}
}
