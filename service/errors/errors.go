package errors

import "github.com/micro/go-micro/v3/errors"

var (
	BadRequest          = errors.BadRequest
	Unauthorized        = errors.Unauthorized
	Forbidden           = errors.Forbidden
	NotFound            = errors.NotFound
	MethodNotAllowed    = errors.MethodNotAllowed
	Timeout             = errors.Timeout
	Conflict            = errors.Conflict
	InternalServerError = errors.InternalServerError
	NotImplemented      = errors.NotImplemented
	BadGateway          = errors.BadGateway
	ServiceUnavailable  = errors.ServiceUnavailable
	GatewayTimeout      = errors.GatewayTimeout

	Equal = errors.Equal
)

// Parse an error into a go-micro error
func Parse(err error) *errors.Error {
	verr, _ := err.(*errors.Error)
	return verr
}
