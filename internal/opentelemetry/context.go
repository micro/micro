package opentelemetry

// This file borrows heavily from https://github.com/grpc-ecosystem/grpc-opentracing/tree/master/go/otgrpc

// // ExtractSpanContext uses the default tracer to extract context:
// func ExtractSpanContext(ctx context.Context) (opentracing.SpanContext, error) {
// 	md, ok := metadata.FromIncomingContext(ctx)
// 	if !ok {
// 		md = metadata.New(nil)
// 	}
// 	return DefaultOpenTracer.Extract(opentracing.HTTPHeaders, metadataReaderWriter{md})
// }

// // InjectSpanContext uses the default tracer to inject context:
// func InjectSpanContext(ctx context.Context, clientSpan opentracing.Span) context.Context {
// 	md, ok := metadata.FromOutgoingContext(ctx)
// 	if !ok {
// 		md = metadata.New(nil)
// 	} else {
// 		md = md.Copy()
// 	}
// 	mdWriter := metadataReaderWriter{md}
// 	err := DefaultOpenTracer.Inject(clientSpan.Context(), opentracing.HTTPHeaders, mdWriter)
// 	// We have no better place to record an error than the Span itself :-/
// 	if err != nil {
// 		clientSpan.LogFields(log.String("event", "opentelemetry.InjectSpanContext() failed"), log.Error(err))
// 	}
// 	return metadata.NewOutgoingContext(ctx, md)
// }

// type metadataReaderWriter struct {
// 	metadata.MD
// }

// func (w metadataReaderWriter) Set(key, val string) {
// 	// The GRPC HPACK implementation rejects any uppercase keys here.
// 	//
// 	// As such, since the HTTP_HEADERS format is case-insensitive anyway, we
// 	// blindly lowercase the key (which is guaranteed to work in the
// 	// Inject/Extract sense per the OpenTracing spec).
// 	key = strings.ToLower(key)
// 	w.MD[key] = append(w.MD[key], val)
// }

// func (w metadataReaderWriter) ForeachKey(handler func(key, val string) error) error {
// 	for k, vals := range w.MD {
// 		for _, v := range vals {
// 			if err := handler(k, v); err != nil {
// 				return err
// 			}
// 		}
// 	}

// 	return nil
// }
