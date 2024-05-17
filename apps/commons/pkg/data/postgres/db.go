package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/utr1903/monitoring-applications-with-opentelemetry/apps/commons/pkg/data"
)

type Database struct {
	createDbNotReachableError bool
}

func NewDatabase(createDbNotReachableError bool) *Database {
	return &Database{
		createDbNotReachableError: createDbNotReachableError,
	}
}

func (db *Database) Create(ctx context.Context, query string) *data.CreateResponse {
	if db.createDbNotReachableError {
		return &data.CreateResponse{
			Success: false,
			Message: "Creating task failed. Database is not reachable.",
			Body:    nil,
		}
	}

	id, _ := uuid.NewUUID()
	return &data.CreateResponse{
		Success: true,
		Message: "Creating task succeeded.",
		Body: &data.Task{
			Id:      id,
			Message: query,
		},
	}
}

func (db *Database) List(ctx context.Context, query string) *data.ListResponse {
	if db.createDbNotReachableError {
		return &data.ListResponse{
			Success: false,
			Message: "Listing tasks failed. Database is not reachable.",
			Body:    nil,
		}
	}

	tasks := make([]data.Task, 5)
	for i := 1; i <= 5; i++ {
		id, _ := uuid.NewUUID()
		tasks = append(tasks, data.Task{
			Id:      id,
			Message: "Some task.",
		})
	}

	return &data.ListResponse{
		Success: true,
		Message: "Listing tasks succeeded.",
		Body:    tasks,
	}
}

func (db *Database) Delete(ctx context.Context, query string) *data.DeleteResponse {
	if db.createDbNotReachableError {
		return &data.DeleteResponse{
			Success: false,
			Message: "Deleting tasks failed. Database is not reachable.",
		}
	}

	return &data.DeleteResponse{
		Success: true,
		Message: "Deleting tasks succeeded.",
	}
}
