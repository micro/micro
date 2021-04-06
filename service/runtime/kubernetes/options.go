package kubernetes

import (
	"context"

	"github.com/micro/micro/v3/service/runtime"
)

type runtimeClassNameKey struct{}

// RuntimeClassName sets the runtimeClassName pods will be started with, e.g. kata-fc
func RuntimeClassName(rcn string) runtime.Option {
	return func(o *runtime.Options) {
		if o.Context == nil {
			o.Context = context.WithValue(context.TODO(), runtimeClassNameKey{}, rcn)
		} else {
			o.Context = context.WithValue(o.Context, runtimeClassNameKey{}, rcn)
		}
	}
}

func getRuntimeClassName(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	return ctx.Value(runtimeClassNameKey{}).(string)
}
