package router;

import static io.grpc.MethodDescriptor.generateFullMethodName;
import static io.grpc.stub.ClientCalls.asyncBidiStreamingCall;
import static io.grpc.stub.ClientCalls.asyncClientStreamingCall;
import static io.grpc.stub.ClientCalls.asyncServerStreamingCall;
import static io.grpc.stub.ClientCalls.asyncUnaryCall;
import static io.grpc.stub.ClientCalls.blockingServerStreamingCall;
import static io.grpc.stub.ClientCalls.blockingUnaryCall;
import static io.grpc.stub.ClientCalls.futureUnaryCall;
import static io.grpc.stub.ServerCalls.asyncBidiStreamingCall;
import static io.grpc.stub.ServerCalls.asyncClientStreamingCall;
import static io.grpc.stub.ServerCalls.asyncServerStreamingCall;
import static io.grpc.stub.ServerCalls.asyncUnaryCall;
import static io.grpc.stub.ServerCalls.asyncUnimplementedStreamingCall;
import static io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall;

/**
 * <pre>
 * Router service is used by the proxy to lookup routes
 * </pre>
 */
@javax.annotation.Generated(
    value = "by gRPC proto compiler (version 1.28.0)",
    comments = "Source: router/router.proto")
public final class RouterGrpc {

  private RouterGrpc() {}

  public static final String SERVICE_NAME = "router.Router";

  // Static method descriptors that strictly reflect the proto.
  private static volatile io.grpc.MethodDescriptor<router.RouterOuterClass.LookupRequest,
      router.RouterOuterClass.LookupResponse> getLookupMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "Lookup",
      requestType = router.RouterOuterClass.LookupRequest.class,
      responseType = router.RouterOuterClass.LookupResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<router.RouterOuterClass.LookupRequest,
      router.RouterOuterClass.LookupResponse> getLookupMethod() {
    io.grpc.MethodDescriptor<router.RouterOuterClass.LookupRequest, router.RouterOuterClass.LookupResponse> getLookupMethod;
    if ((getLookupMethod = RouterGrpc.getLookupMethod) == null) {
      synchronized (RouterGrpc.class) {
        if ((getLookupMethod = RouterGrpc.getLookupMethod) == null) {
          RouterGrpc.getLookupMethod = getLookupMethod =
              io.grpc.MethodDescriptor.<router.RouterOuterClass.LookupRequest, router.RouterOuterClass.LookupResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "Lookup"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  router.RouterOuterClass.LookupRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  router.RouterOuterClass.LookupResponse.getDefaultInstance()))
              .setSchemaDescriptor(new RouterMethodDescriptorSupplier("Lookup"))
              .build();
        }
      }
    }
    return getLookupMethod;
  }

  private static volatile io.grpc.MethodDescriptor<router.RouterOuterClass.WatchRequest,
      router.RouterOuterClass.Event> getWatchMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "Watch",
      requestType = router.RouterOuterClass.WatchRequest.class,
      responseType = router.RouterOuterClass.Event.class,
      methodType = io.grpc.MethodDescriptor.MethodType.SERVER_STREAMING)
  public static io.grpc.MethodDescriptor<router.RouterOuterClass.WatchRequest,
      router.RouterOuterClass.Event> getWatchMethod() {
    io.grpc.MethodDescriptor<router.RouterOuterClass.WatchRequest, router.RouterOuterClass.Event> getWatchMethod;
    if ((getWatchMethod = RouterGrpc.getWatchMethod) == null) {
      synchronized (RouterGrpc.class) {
        if ((getWatchMethod = RouterGrpc.getWatchMethod) == null) {
          RouterGrpc.getWatchMethod = getWatchMethod =
              io.grpc.MethodDescriptor.<router.RouterOuterClass.WatchRequest, router.RouterOuterClass.Event>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.SERVER_STREAMING)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "Watch"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  router.RouterOuterClass.WatchRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  router.RouterOuterClass.Event.getDefaultInstance()))
              .setSchemaDescriptor(new RouterMethodDescriptorSupplier("Watch"))
              .build();
        }
      }
    }
    return getWatchMethod;
  }

  /**
   * Creates a new async stub that supports all call types for the service
   */
  public static RouterStub newStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<RouterStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<RouterStub>() {
        @java.lang.Override
        public RouterStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new RouterStub(channel, callOptions);
        }
      };
    return RouterStub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static RouterBlockingStub newBlockingStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<RouterBlockingStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<RouterBlockingStub>() {
        @java.lang.Override
        public RouterBlockingStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new RouterBlockingStub(channel, callOptions);
        }
      };
    return RouterBlockingStub.newStub(factory, channel);
  }

  /**
   * Creates a new ListenableFuture-style stub that supports unary calls on the service
   */
  public static RouterFutureStub newFutureStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<RouterFutureStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<RouterFutureStub>() {
        @java.lang.Override
        public RouterFutureStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new RouterFutureStub(channel, callOptions);
        }
      };
    return RouterFutureStub.newStub(factory, channel);
  }

  /**
   * <pre>
   * Router service is used by the proxy to lookup routes
   * </pre>
   */
  public static abstract class RouterImplBase implements io.grpc.BindableService {

    /**
     */
    public void lookup(router.RouterOuterClass.LookupRequest request,
        io.grpc.stub.StreamObserver<router.RouterOuterClass.LookupResponse> responseObserver) {
      asyncUnimplementedUnaryCall(getLookupMethod(), responseObserver);
    }

    /**
     */
    public void watch(router.RouterOuterClass.WatchRequest request,
        io.grpc.stub.StreamObserver<router.RouterOuterClass.Event> responseObserver) {
      asyncUnimplementedUnaryCall(getWatchMethod(), responseObserver);
    }

    @java.lang.Override public final io.grpc.ServerServiceDefinition bindService() {
      return io.grpc.ServerServiceDefinition.builder(getServiceDescriptor())
          .addMethod(
            getLookupMethod(),
            asyncUnaryCall(
              new MethodHandlers<
                router.RouterOuterClass.LookupRequest,
                router.RouterOuterClass.LookupResponse>(
                  this, METHODID_LOOKUP)))
          .addMethod(
            getWatchMethod(),
            asyncServerStreamingCall(
              new MethodHandlers<
                router.RouterOuterClass.WatchRequest,
                router.RouterOuterClass.Event>(
                  this, METHODID_WATCH)))
          .build();
    }
  }

  /**
   * <pre>
   * Router service is used by the proxy to lookup routes
   * </pre>
   */
  public static final class RouterStub extends io.grpc.stub.AbstractAsyncStub<RouterStub> {
    private RouterStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected RouterStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new RouterStub(channel, callOptions);
    }

    /**
     */
    public void lookup(router.RouterOuterClass.LookupRequest request,
        io.grpc.stub.StreamObserver<router.RouterOuterClass.LookupResponse> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(getLookupMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void watch(router.RouterOuterClass.WatchRequest request,
        io.grpc.stub.StreamObserver<router.RouterOuterClass.Event> responseObserver) {
      asyncServerStreamingCall(
          getChannel().newCall(getWatchMethod(), getCallOptions()), request, responseObserver);
    }
  }

  /**
   * <pre>
   * Router service is used by the proxy to lookup routes
   * </pre>
   */
  public static final class RouterBlockingStub extends io.grpc.stub.AbstractBlockingStub<RouterBlockingStub> {
    private RouterBlockingStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected RouterBlockingStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new RouterBlockingStub(channel, callOptions);
    }

    /**
     */
    public router.RouterOuterClass.LookupResponse lookup(router.RouterOuterClass.LookupRequest request) {
      return blockingUnaryCall(
          getChannel(), getLookupMethod(), getCallOptions(), request);
    }

    /**
     */
    public java.util.Iterator<router.RouterOuterClass.Event> watch(
        router.RouterOuterClass.WatchRequest request) {
      return blockingServerStreamingCall(
          getChannel(), getWatchMethod(), getCallOptions(), request);
    }
  }

  /**
   * <pre>
   * Router service is used by the proxy to lookup routes
   * </pre>
   */
  public static final class RouterFutureStub extends io.grpc.stub.AbstractFutureStub<RouterFutureStub> {
    private RouterFutureStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected RouterFutureStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new RouterFutureStub(channel, callOptions);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<router.RouterOuterClass.LookupResponse> lookup(
        router.RouterOuterClass.LookupRequest request) {
      return futureUnaryCall(
          getChannel().newCall(getLookupMethod(), getCallOptions()), request);
    }
  }

  private static final int METHODID_LOOKUP = 0;
  private static final int METHODID_WATCH = 1;

  private static final class MethodHandlers<Req, Resp> implements
      io.grpc.stub.ServerCalls.UnaryMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ServerStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ClientStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.BidiStreamingMethod<Req, Resp> {
    private final RouterImplBase serviceImpl;
    private final int methodId;

    MethodHandlers(RouterImplBase serviceImpl, int methodId) {
      this.serviceImpl = serviceImpl;
      this.methodId = methodId;
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public void invoke(Req request, io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        case METHODID_LOOKUP:
          serviceImpl.lookup((router.RouterOuterClass.LookupRequest) request,
              (io.grpc.stub.StreamObserver<router.RouterOuterClass.LookupResponse>) responseObserver);
          break;
        case METHODID_WATCH:
          serviceImpl.watch((router.RouterOuterClass.WatchRequest) request,
              (io.grpc.stub.StreamObserver<router.RouterOuterClass.Event>) responseObserver);
          break;
        default:
          throw new AssertionError();
      }
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public io.grpc.stub.StreamObserver<Req> invoke(
        io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        default:
          throw new AssertionError();
      }
    }
  }

  private static abstract class RouterBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier, io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    RouterBaseDescriptorSupplier() {}

    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return router.RouterOuterClass.getDescriptor();
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.ServiceDescriptor getServiceDescriptor() {
      return getFileDescriptor().findServiceByName("Router");
    }
  }

  private static final class RouterFileDescriptorSupplier
      extends RouterBaseDescriptorSupplier {
    RouterFileDescriptorSupplier() {}
  }

  private static final class RouterMethodDescriptorSupplier
      extends RouterBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoMethodDescriptorSupplier {
    private final String methodName;

    RouterMethodDescriptorSupplier(String methodName) {
      this.methodName = methodName;
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.MethodDescriptor getMethodDescriptor() {
      return getServiceDescriptor().findMethodByName(methodName);
    }
  }

  private static volatile io.grpc.ServiceDescriptor serviceDescriptor;

  public static io.grpc.ServiceDescriptor getServiceDescriptor() {
    io.grpc.ServiceDescriptor result = serviceDescriptor;
    if (result == null) {
      synchronized (RouterGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor = result = io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
              .setSchemaDescriptor(new RouterFileDescriptorSupplier())
              .addMethod(getLookupMethod())
              .addMethod(getWatchMethod())
              .build();
        }
      }
    }
    return result;
  }
}
