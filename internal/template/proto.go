package template

var (
	ProtoFNC = `syntax = "proto3";

package {{.FQDN}};

service Example {
	rpc Call(Request) returns (Response) {}
}

message Message {
	string say = 1;
}

message Request {
	string name = 1;
}

message Response {
	string msg = 1;
}
`

	ProtoSRV = `syntax = "proto3";

package {{.FQDN}};

service Example {
	rpc Call(Request) returns (Response) {}
	rpc Stream(StreamingRequest) returns (stream StreamingResponse) {}
	rpc PingPong(stream Ping) returns (stream Pong) {}
}

message Message {
	string say = 1;
}

message Request {
	string name = 1;
}

message Response {
	string msg = 1;
}

message StreamingRequest {
	int64 count = 1;
}

message StreamingResponse {
	int64 count = 1;
}

message Ping {
	int64 stroke = 1;
}

message Pong {
	int64 stroke = 1;
}
`

	ProtoAPI = `syntax = "proto3";

package {{.FQDN}};

import "github.com/micro/go-api/proto/api.proto";

service Example {
	rpc Call(api.Request) returns (api.Response) {}
}
`
)
