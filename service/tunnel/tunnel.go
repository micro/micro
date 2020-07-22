package tunnel

import (
	"github.com/micro/go-micro/v2/transport"
	thttp "github.com/micro/go-micro/v2/transport/http"
)

var DefaultTransport transport.Transport = thttp.NewTransport()
