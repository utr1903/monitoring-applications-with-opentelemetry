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

type ListResponse struct {
	Result bool
	Tasks  []Task
}

type IEntityService interface {
	Create(task string) *CreateResponse
	List() *ListResponse
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

func (s *UnimplementedEntityService) List() *ListResponse {
	tasks := make([]Task, 5)

	for i := 1; i <= 5; i++ {
		id, _ := uuid.NewUUID()
		tasks = append(tasks, Task{
			Id:      id,
			Message: "Some task.",
		})
	}

	return &ListResponse{
		Result: true,
		Tasks:  tasks,
	}
}
