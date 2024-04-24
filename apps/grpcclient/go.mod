module github.com/utr1903/monitoring-applications-with-opentelemetry/apps/grpcclient

go 1.21

replace github.com/utr1903/monitoring-applications-with-opentelemetry/apps/commons => ../commons

require (
	github.com/utr1903/monitoring-applications-with-opentelemetry/apps/commons v0.0.0
	google.golang.org/grpc v1.63.2
	google.golang.org/protobuf v1.33.0
)

require (
	github.com/sirupsen/logrus v1.9.3 // indirect
	golang.org/x/net v0.21.0 // indirect
	golang.org/x/sys v0.19.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240227224415-6ceb2ff114de // indirect
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.3.0 // indirect
)
