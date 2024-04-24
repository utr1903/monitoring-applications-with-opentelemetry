package data

import "github.com/google/uuid"

type Task struct {
	Id      uuid.UUID
	Message string
}

type CreateResponse struct {
	Result  bool
	Message string
	Body    *Task
}

type ListResponse struct {
	Result  bool
	Message string
	Body    []Task
}

type DeleteResponse struct {
	Result  bool
	Message string
}

type IEntityService interface {
	Create(task string) *CreateResponse
	List() *ListResponse
	Delete() *DeleteResponse
}

type UnimplementedEntityService struct {
}

func NewUnimplementedEntityService() *UnimplementedEntityService {
	return &UnimplementedEntityService{}
}

func (s *UnimplementedEntityService) Create(message string) *CreateResponse {
	id, _ := uuid.NewUUID()
	return &CreateResponse{
		Result:  true,
		Message: "Creating task succeeded.",
		Body: &Task{
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
		Result:  true,
		Message: "Listing tasks succeeded.",
		Body:    tasks,
	}
}

func (s *UnimplementedEntityService) Delete() *DeleteResponse {
	return &DeleteResponse{
		Result:  true,
		Message: "Deleting tasks succeeded.",
	}
}
