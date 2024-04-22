package data

import "github.com/google/uuid"

type Task struct {
	Id      uuid.UUID
	Message string
}

type CreateResponse struct {
	Result bool
	Task   *Task
}

type IEntityService interface {
	Create(task string) *CreateResponse
}

type UnimplementedEntityService struct {
}

func NewUnimplementedEntityService() *UnimplementedEntityService {
	return &UnimplementedEntityService{}
}

func (s *UnimplementedEntityService) Create(message string) *CreateResponse {
	id, _ := uuid.NewUUID()
	return &CreateResponse{
		Result: true,
		Task: &Task{
			Id:      id,
			Message: message,
		},
	}
}
