package loggers

import "context"

type Level int

const (
	Info Level = iota
	Warning
	Error
)

type ILogger interface {
	Log(ctx context.Context, lvl Level, message string, attrs map[string]string)
}
