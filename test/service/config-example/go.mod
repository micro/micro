module config

go 1.13

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

require (
	github.com/micro/go-micro/v3 v3.0.0-beta.2.0.20200918112555-9168c7c61064 // indirect
	github.com/micro/micro/v3 v3.0.0-20200728090928-ad22505562c9
)

replace github.com/micro/micro/v3 => ../../..
