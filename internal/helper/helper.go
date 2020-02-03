package helper

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2/metadata"
)

func ACMEHosts(ctx *cli.Context) []string {
	var hosts []string
	for _, host := range strings.Split(ctx.String("acme_hosts"), ",") {
		if len(host) > 0 {
			hosts = append(hosts, host)
		}
	}
	return hosts
}

func RequestToContext(r *http.Request) context.Context {
	ctx := context.Background()
	md := make(metadata.Metadata)
	for k, v := range r.Header {
		md[k] = strings.Join(v, ",")
	}
	return metadata.NewContext(ctx, md)
}

func TLSConfig(ctx *cli.Context) (*tls.Config, error) {
	cert := ctx.String("tls_cert_file")
	key := ctx.String("tls_key_file")
	ca := ctx.String("tls_client_ca_file")

	if len(cert) > 0 && len(key) > 0 {
		certs, err := tls.LoadX509KeyPair(cert, key)
		if err != nil {
			return nil, err
		}

		if len(ca) > 0 {
			caCert, err := ioutil.ReadFile(ca)
			if err != nil {
				return nil, err
			}

			caCertPool := x509.NewCertPool()
			caCertPool.AppendCertsFromPEM(caCert)

			return &tls.Config{
				Certificates: []tls.Certificate{certs},
				ClientCAs:    caCertPool,
				ClientAuth:   tls.RequireAndVerifyClientCert,
				NextProtos:   []string{"h2", "http/1.1"},
			}, nil
		}

		return &tls.Config{
			Certificates: []tls.Certificate{certs}, NextProtos: []string{"h2", "http/1.1"},
		}, nil
	}

	return nil, errors.New("TLS certificate and key files not specified")
}

func ServeCORS(w http.ResponseWriter, r *http.Request) {
	set := func(w http.ResponseWriter, k, v string) {
		if v := w.Header().Get(k); len(v) > 0 {
			return
		}
		w.Header().Set(k, v)
	}

	if origin := r.Header.Get("Origin"); len(origin) > 0 {
		set(w, "Access-Control-Allow-Origin", origin)
	} else {
		set(w, "Access-Control-Allow-Origin", "*")
	}

	set(w, "Access-Control-Allow-Methods", "POST, PATCH, GET, OPTIONS, PUT, DELETE")
	set(w, "Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}
