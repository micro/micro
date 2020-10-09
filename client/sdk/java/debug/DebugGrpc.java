package debug;

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
 */
@javax.annotation.Generated(
    value = "by gRPC proto compiler (version 1.28.0)",
    comments = "Source: debug/debug.proto")
public final class DebugGrpc {

  private DebugGrpc() {}

  public static final String SERVICE_NAME = "debug.Debug";

  // Static method descriptors that strictly reflect the proto.
  private static volatile io.grpc.MethodDescriptor<debug.DebugOuterClass.LogRequest,
      debug.DebugOuterClass.LogResponse> getLogMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "Log",
      requestType = debug.DebugOuterClass.LogRequest.class,
      responseType = debug.DebugOuterClass.LogResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<debug.DebugOuterClass.LogRequest,
      debug.DebugOuterClass.LogResponse> getLogMethod() {
    io.grpc.MethodDescriptor<debug.DebugOuterClass.LogRequest, debug.DebugOuterClass.LogResponse> getLogMethod;
    if ((getLogMethod = DebugGrpc.getLogMethod) == null) {
      synchronized (DebugGrpc.class) {
        if ((getLogMethod = DebugGrpc.getLogMethod) == null) {
          DebugGrpc.getLogMethod = getLogMethod =
              io.grpc.MethodDescriptor.<debug.DebugOuterClass.LogRequest, debug.DebugOuterClass.LogResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "Log"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  debug.DebugOuterClass.LogRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  debug.DebugOuterClass.LogResponse.getDefaultInstance()))
              .setSchemaDescriptor(new DebugMethodDescriptorSupplier("Log"))
              .build();
        }
      }
    }
    return getLogMethod;
  }

  private static volatile io.grpc.MethodDescriptor<debug.DebugOuterClass.HealthRequest,
      debug.DebugOuterClass.HealthResponse> getHealthMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "Health",
      requestType = debug.DebugOuterClass.HealthRequest.class,
      responseType = debug.DebugOuterClass.HealthResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<debug.DebugOuterClass.HealthRequest,
      debug.DebugOuterClass.HealthResponse> getHealthMethod() {
    io.grpc.MethodDescriptor<debug.DebugOuterClass.HealthRequest, debug.DebugOuterClass.HealthResponse> getHealthMethod;
    if ((getHealthMethod = DebugGrpc.getHealthMethod) == null) {
      synchronized (DebugGrpc.class) {
        if ((getHealthMethod = DebugGrpc.getHealthMethod) == null) {
          DebugGrpc.getHealthMethod = getHealthMethod =
              io.grpc.MethodDescriptor.<debug.DebugOuterClass.HealthRequest, debug.DebugOuterClass.HealthResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "Health"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  debug.DebugOuterClass.HealthRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  debug.DebugOuterClass.HealthResponse.getDefaultInstance()))
              .setSchemaDescriptor(new DebugMethodDescriptorSupplier("Health"))
              .build();
        }
      }
    }
    return getHealthMethod;
  }

  private static volatile io.grpc.MethodDescriptor<debug.DebugOuterClass.StatsRequest,
      debug.DebugOuterClass.StatsResponse> getStatsMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "Stats",
      requestType = debug.DebugOuterClass.StatsRequest.class,
      responseType = debug.DebugOuterClass.StatsResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<debug.DebugOuterClass.StatsRequest,
      debug.DebugOuterClass.StatsResponse> getStatsMethod() {
    io.grpc.MethodDescriptor<debug.DebugOuterClass.StatsRequest, debug.DebugOuterClass.StatsResponse> getStatsMethod;
    if ((getStatsMethod = DebugGrpc.getStatsMethod) == null) {
      synchronized (DebugGrpc.class) {
        if ((getStatsMethod = DebugGrpc.getStatsMethod) == null) {
          DebugGrpc.getStatsMethod = getStatsMethod =
              io.grpc.MethodDescriptor.<debug.DebugOuterClass.StatsRequest, debug.DebugOuterClass.StatsResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "Stats"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  debug.DebugOuterClass.StatsRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  debug.DebugOuterClass.StatsResponse.getDefaultInstance()))
              .setSchemaDescriptor(new DebugMethodDescriptorSupplier("Stats"))
              .build();
        }
      }
    }
    return getStatsMethod;
  }

  private static volatile io.grpc.MethodDescriptor<debug.DebugOuterClass.TraceRequest,
      debug.DebugOuterClass.TraceResponse> getTraceMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "Trace",
      requestType = debug.DebugOuterClass.TraceRequest.class,
      responseType = debug.DebugOuterClass.TraceResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<debug.DebugOuterClass.TraceRequest,
      debug.DebugOuterClass.TraceResponse> getTraceMethod() {
    io.grpc.MethodDescriptor<debug.DebugOuterClass.TraceRequest, debug.DebugOuterClass.TraceResponse> getTraceMethod;
    if ((getTraceMethod = DebugGrpc.getTraceMethod) == null) {
      synchronized (DebugGrpc.class) {
        if ((getTraceMethod = DebugGrpc.getTraceMethod) == null) {
          DebugGrpc.getTraceMethod = getTraceMethod =
              io.grpc.MethodDescriptor.<debug.DebugOuterClass.TraceRequest, debug.DebugOuterClass.TraceResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "Trace"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  debug.DebugOuterClass.TraceRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  debug.DebugOuterClass.TraceResponse.getDefaultInstance()))
              .setSchemaDescriptor(new DebugMethodDescriptorSupplier("Trace"))
              .build();
        }
      }
    }
    return getTraceMethod;
  }

  /**
   * Creates a new async stub that supports all call types for the service
   */
  public static DebugStub newStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<DebugStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<DebugStub>() {
        @java.lang.Override
        public DebugStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new DebugStub(channel, callOptions);
        }
      };
    return DebugStub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static DebugBlockingStub newBlockingStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<DebugBlockingStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<DebugBlockingStub>() {
        @java.lang.Override
        public DebugBlockingStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new DebugBlockingStub(channel, callOptions);
        }
      };
    return DebugBlockingStub.newStub(factory, channel);
  }

  /**
   * Creates a new ListenableFuture-style stub that supports unary calls on the service
   */
  public static DebugFutureStub newFutureStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<DebugFutureStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<DebugFutureStub>() {
        @java.lang.Override
        public DebugFutureStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new DebugFutureStub(channel, callOptions);
        }
      };
    return DebugFutureStub.newStub(factory, channel);
  }

  /**
   */
  public static abstract class DebugImplBase implements io.grpc.BindableService {

    /**
     */
    public void log(debug.DebugOuterClass.LogRequest request,
        io.grpc.stub.StreamObserver<debug.DebugOuterClass.LogResponse> responseObserver) {
      asyncUnimplementedUnaryCall(getLogMethod(), responseObserver);
    }

    /**
     */
    public void health(debug.DebugOuterClass.HealthRequest request,
        io.grpc.stub.StreamObserver<debug.DebugOuterClass.HealthResponse> responseObserver) {
      asyncUnimplementedUnaryCall(getHealthMethod(), responseObserver);
    }

    /**
     */
    public void stats(debug.DebugOuterClass.StatsRequest request,
        io.grpc.stub.StreamObserver<debug.DebugOuterClass.StatsResponse> responseObserver) {
      asyncUnimplementedUnaryCall(getStatsMethod(), responseObserver);
    }

    /**
     */
    public void trace(debug.DebugOuterClass.TraceRequest request,
        io.grpc.stub.StreamObserver<debug.DebugOuterClass.TraceResponse> responseObserver) {
      asyncUnimplementedUnaryCall(getTraceMethod(), responseObserver);
    }

    @java.lang.Override public final io.grpc.ServerServiceDefinition bindService() {
      return io.grpc.ServerServiceDefinition.builder(getServiceDescriptor())
          .addMethod(
            getLogMethod(),
            asyncUnaryCall(
              new MethodHandlers<
                debug.DebugOuterClass.LogRequest,
                debug.DebugOuterClass.LogResponse>(
                  this, METHODID_LOG)))
          .addMethod(
            getHealthMethod(),
            asyncUnaryCall(
              new MethodHandlers<
                debug.DebugOuterClass.HealthRequest,
                debug.DebugOuterClass.HealthResponse>(
                  this, METHODID_HEALTH)))
          .addMethod(
            getStatsMethod(),
            asyncUnaryCall(
              new MethodHandlers<
                debug.DebugOuterClass.StatsRequest,
                debug.DebugOuterClass.StatsResponse>(
                  this, METHODID_STATS)))
          .addMethod(
            getTraceMethod(),
            asyncUnaryCall(
              new MethodHandlers<
                debug.DebugOuterClass.TraceRequest,
                debug.DebugOuterClass.TraceResponse>(
                  this, METHODID_TRACE)))
          .build();
    }
  }

  /**
   */
  public static final class DebugStub extends io.grpc.stub.AbstractAsyncStub<DebugStub> {
    private DebugStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected DebugStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new DebugStub(channel, callOptions);
    }

    /**
     */
    public void log(debug.DebugOuterClass.LogRequest request,
        io.grpc.stub.StreamObserver<debug.DebugOuterClass.LogResponse> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(getLogMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void health(debug.DebugOuterClass.HealthRequest request,
        io.grpc.stub.StreamObserver<debug.DebugOuterClass.HealthResponse> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(getHealthMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void stats(debug.DebugOuterClass.StatsRequest request,
        io.grpc.stub.StreamObserver<debug.DebugOuterClass.StatsResponse> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(getStatsMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void trace(debug.DebugOuterClass.TraceRequest request,
        io.grpc.stub.StreamObserver<debug.DebugOuterClass.TraceResponse> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(getTraceMethod(), getCallOptions()), request, responseObserver);
    }
  }

  /**
   */
  public static final class DebugBlockingStub extends io.grpc.stub.AbstractBlockingStub<DebugBlockingStub> {
    private DebugBlockingStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected DebugBlockingStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new DebugBlockingStub(channel, callOptions);
    }

    /**
     */
    public debug.DebugOuterClass.LogResponse log(debug.DebugOuterClass.LogRequest request) {
      return blockingUnaryCall(
          getChannel(), getLogMethod(), getCallOptions(), request);
    }

    /**
     */
    public debug.DebugOuterClass.HealthResponse health(debug.DebugOuterClass.HealthRequest request) {
      return blockingUnaryCall(
          getChannel(), getHealthMethod(), getCallOptions(), request);
    }

    /**
     */
    public debug.DebugOuterClass.StatsResponse stats(debug.DebugOuterClass.StatsRequest request) {
      return blockingUnaryCall(
          getChannel(), getStatsMethod(), getCallOptions(), request);
    }

    /**
     */
    public debug.DebugOuterClass.TraceResponse trace(debug.DebugOuterClass.TraceRequest request) {
      return blockingUnaryCall(
          getChannel(), getTraceMethod(), getCallOptions(), request);
    }
  }

  /**
   */
  public static final class DebugFutureStub extends io.grpc.stub.AbstractFutureStub<DebugFutureStub> {
    private DebugFutureStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected DebugFutureStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new DebugFutureStub(channel, callOptions);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<debug.DebugOuterClass.LogResponse> log(
        debug.DebugOuterClass.LogRequest request) {
      return futureUnaryCall(
          getChannel().newCall(getLogMethod(), getCallOptions()), request);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<debug.DebugOuterClass.HealthResponse> health(
        debug.DebugOuterClass.HealthRequest request) {
      return futureUnaryCall(
          getChannel().newCall(getHealthMethod(), getCallOptions()), request);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<debug.DebugOuterClass.StatsResponse> stats(
        debug.DebugOuterClass.StatsRequest request) {
      return futureUnaryCall(
          getChannel().newCall(getStatsMethod(), getCallOptions()), request);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<debug.DebugOuterClass.TraceResponse> trace(
        debug.DebugOuterClass.TraceRequest request) {
      return futureUnaryCall(
          getChannel().newCall(getTraceMethod(), getCallOptions()), request);
    }
  }

  private static final int METHODID_LOG = 0;
  private static final int METHODID_HEALTH = 1;
  private static final int METHODID_STATS = 2;
  private static final int METHODID_TRACE = 3;

  private static final class MethodHandlers<Req, Resp> implements
      io.grpc.stub.ServerCalls.UnaryMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ServerStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ClientStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.BidiStreamingMethod<Req, Resp> {
    private final DebugImplBase serviceImpl;
    private final int methodId;

    MethodHandlers(DebugImplBase serviceImpl, int methodId) {
      this.serviceImpl = serviceImpl;
      this.methodId = methodId;
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public void invoke(Req request, io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        case METHODID_LOG:
          serviceImpl.log((debug.DebugOuterClass.LogRequest) request,
              (io.grpc.stub.StreamObserver<debug.DebugOuterClass.LogResponse>) responseObserver);
          break;
        case METHODID_HEALTH:
          serviceImpl.health((debug.DebugOuterClass.HealthRequest) request,
              (io.grpc.stub.StreamObserver<debug.DebugOuterClass.HealthResponse>) responseObserver);
          break;
        case METHODID_STATS:
          serviceImpl.stats((debug.DebugOuterClass.StatsRequest) request,
              (io.grpc.stub.StreamObserver<debug.DebugOuterClass.StatsResponse>) responseObserver);
          break;
        case METHODID_TRACE:
          serviceImpl.trace((debug.DebugOuterClass.TraceRequest) request,
              (io.grpc.stub.StreamObserver<debug.DebugOuterClass.TraceResponse>) responseObserver);
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

  private static abstract class DebugBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier, io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    DebugBaseDescriptorSupplier() {}

    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return debug.DebugOuterClass.getDescriptor();
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.ServiceDescriptor getServiceDescriptor() {
      return getFileDescriptor().findServiceByName("Debug");
    }
  }

  private static final class DebugFileDescriptorSupplier
      extends DebugBaseDescriptorSupplier {
    DebugFileDescriptorSupplier() {}
  }

  private static final class DebugMethodDescriptorSupplier
      extends DebugBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoMethodDescriptorSupplier {
    private final String methodName;

    DebugMethodDescriptorSupplier(String methodName) {
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
      synchronized (DebugGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor = result = io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
              .setSchemaDescriptor(new DebugFileDescriptorSupplier())
              .addMethod(getLogMethod())
              .addMethod(getHealthMethod())
              .addMethod(getStatsMethod())
              .addMethod(getTraceMethod())
              .build();
        }
      }
    }
    return result;
  }
}
