package server

import "github.com/utr1903/monitoring-applications-with-opentelemetry/apps/commons/pkg/data"

type DeleteRequest struct {
}

type DeleteResult struct {
	Result  bool
	Message string
}

type IDeleteService interface {
	Delete(*DeleteRequest) *DeleteResult
}

type DeleteService struct {
	entityService data.IEntityService
}

func NewDeleteService() *DeleteService {
	return &DeleteService{
		entityService: data.NewUnimplementedEntityService(),
	}
}

func (s *DeleteService) Delete(req *DeleteRequest) *DeleteResult {
	res := s.entityService.Delete()

	return &DeleteResult{
		Result:  res.Result,
		Message: res.Message,
	}
}
