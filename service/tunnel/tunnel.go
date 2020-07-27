package tunnel

import (
	"github.com/micro/go-micro/v3/transport"
	thttp "github.com/micro/go-micro/v3/transport/http"
)

var DefaultTransport transport.Transport = thttp.NewTransport()
