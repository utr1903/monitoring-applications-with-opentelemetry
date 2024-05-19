package data

import (
	"context"

	"github.com/google/uuid"
)

type Task struct {
	Id      uuid.UUID
	Message string
}

type CreateResponse struct {
	Success bool
	Message string
	Body    *Task
}

type ListResponse struct {
	Success bool
	Message string
	Body    []Task
}

type DeleteResponse struct {
	Success bool
	Message string
}

type IEntityService interface {
	Create(ctx context.Context, taskMessage string) *CreateResponse
	List(ctx context.Context, query string) *ListResponse
	Delete(ctx context.Context, query string) *DeleteResponse
}
