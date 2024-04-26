package config

import (
	"os"
	"strconv"
)

type Config struct {
	ServiceName string
	Port        string

	StoreDelay  int
	ListDelay   int
	DeleteDelay int
}

func NewConfig() *Config {
	serviceName := os.Getenv("OTEL_SERVICE_NAME")
	port := os.Getenv("PORT")

	storeDelay, err := strconv.Atoi(os.Getenv("STORE_DELAY"))
	if err != nil {
		panic("STORE_DELAY could not be parsed into an integer.")
	}
	listDelay, err := strconv.Atoi(os.Getenv("LIST_DELAY"))
	if err != nil {
		panic("LIST_DELAY could not be parsed into an integer.")
	}
	deleteDelay, err := strconv.Atoi(os.Getenv("DELETE_DELAY"))
	if err != nil {
		panic("DELETE_DELAY could not be parsed into an integer.")
	}

	return &Config{
		ServiceName: serviceName,
		Port:        port,

		StoreDelay:  storeDelay,
		ListDelay:   listDelay,
		DeleteDelay: deleteDelay,
	}
}