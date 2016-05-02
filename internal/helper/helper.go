package helper

import (
	"crypto/tls"
	"errors"
	"net/http"
	"strings"

	"github.com/micro/cli"
	"github.com/micro/go-micro/metadata"

	"golang.org/x/net/context"
)

func RequestToContext(r *http.Request) context.Context {
	ctx := context.Background()
	md := make(metadata.Metadata)
	for k, v := range r.Header {
		md[k] = strings.Join(v, ",")
	}
	return metadata.NewContext(ctx, md)
}

func TLSConfig(ctx *cli.Context) (*tls.Config, error) {
	cert := ctx.GlobalString("tls_cert_file")
	key := ctx.GlobalString("tls_key_file")

	if len(cert) > 0 && len(key) > 0 {
		certs, err := tls.LoadX509KeyPair(cert, key)
		if err != nil {
			return nil, err
		}

		return &tls.Config{
			Certificates: []tls.Certificate{certs},
		}, nil
	}

	return nil, errors.New("TLS certificate and key files not specified")
}
