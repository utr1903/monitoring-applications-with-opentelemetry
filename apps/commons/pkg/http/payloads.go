package http

type Task struct {
	Id      string `json:"id"`
	Message string `json:"message"`
}

type StoreTaskRequest struct {
	Message string `json:"message"`
}

type StoreTaskResponse struct {
	Message string `json:"message"`
	Body    *Task  `json:"body"`
}

type ListTasksResponse struct {
	Message string `json:"message"`
	Body    []Task `json:"body"`
}

type DeleteTasksResponse struct {
	Message string `json:"message"`
}
