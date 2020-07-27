module example-service

go 1.13

require (
	github.com/gogo/protobuf v1.2.2-0.20190723190241-65acae22fc9d // indirect
	github.com/golang/protobuf v1.4.2
	github.com/micro/go-micro/v3 v2.9.1-0.20200724132314-1f697400e61a
	github.com/micro/micro/v2 v2.9.2-0.20200724125730-944073fbedd0
	google.golang.org/grpc v1.27.0
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
