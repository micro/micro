package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/golang/protobuf/proto"
	hello "github.com/micro/micro/examples/greeter/srv/proto/hello"
)

func main() {
	req, err := proto.Marshal(&hello.Request{Name: "John"})
	if err != nil {
		fmt.Println(err)
		return
	}

	r, err := http.Post("http://localhost:8081/greeter/say/hello", "application/protobuf", bytes.NewReader(req))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer r.Body.Close()

	b, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	rsp := &hello.Response{}
	if err := proto.Unmarshal(b, rsp); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(rsp.Msg)
}
