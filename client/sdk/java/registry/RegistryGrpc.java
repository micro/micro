package registry;

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
    comments = "Source: registry/registry.proto")
public final class RegistryGrpc {

  private RegistryGrpc() {}

  public static final String SERVICE_NAME = "registry.Registry";

  // Static method descriptors that strictly reflect the proto.
  private static volatile io.grpc.MethodDescriptor<registry.RegistryOuterClass.GetRequest,
      registry.RegistryOuterClass.GetResponse> getGetServiceMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetService",
      requestType = registry.RegistryOuterClass.GetRequest.class,
      responseType = registry.RegistryOuterClass.GetResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<registry.RegistryOuterClass.GetRequest,
      registry.RegistryOuterClass.GetResponse> getGetServiceMethod() {
    io.grpc.MethodDescriptor<registry.RegistryOuterClass.GetRequest, registry.RegistryOuterClass.GetResponse> getGetServiceMethod;
    if ((getGetServiceMethod = RegistryGrpc.getGetServiceMethod) == null) {
      synchronized (RegistryGrpc.class) {
        if ((getGetServiceMethod = RegistryGrpc.getGetServiceMethod) == null) {
          RegistryGrpc.getGetServiceMethod = getGetServiceMethod =
              io.grpc.MethodDescriptor.<registry.RegistryOuterClass.GetRequest, registry.RegistryOuterClass.GetResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetService"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  registry.RegistryOuterClass.GetRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  registry.RegistryOuterClass.GetResponse.getDefaultInstance()))
              .setSchemaDescriptor(new RegistryMethodDescriptorSupplier("GetService"))
              .build();
        }
      }
    }
    return getGetServiceMethod;
  }

  private static volatile io.grpc.MethodDescriptor<registry.RegistryOuterClass.Service,
      registry.RegistryOuterClass.EmptyResponse> getRegisterMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "Register",
      requestType = registry.RegistryOuterClass.Service.class,
      responseType = registry.RegistryOuterClass.EmptyResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<registry.RegistryOuterClass.Service,
      registry.RegistryOuterClass.EmptyResponse> getRegisterMethod() {
    io.grpc.MethodDescriptor<registry.RegistryOuterClass.Service, registry.RegistryOuterClass.EmptyResponse> getRegisterMethod;
    if ((getRegisterMethod = RegistryGrpc.getRegisterMethod) == null) {
      synchronized (RegistryGrpc.class) {
        if ((getRegisterMethod = RegistryGrpc.getRegisterMethod) == null) {
          RegistryGrpc.getRegisterMethod = getRegisterMethod =
              io.grpc.MethodDescriptor.<registry.RegistryOuterClass.Service, registry.RegistryOuterClass.EmptyResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "Register"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  registry.RegistryOuterClass.Service.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  registry.RegistryOuterClass.EmptyResponse.getDefaultInstance()))
              .setSchemaDescriptor(new RegistryMethodDescriptorSupplier("Register"))
              .build();
        }
      }
    }
    return getRegisterMethod;
  }

  private static volatile io.grpc.MethodDescriptor<registry.RegistryOuterClass.Service,
      registry.RegistryOuterClass.EmptyResponse> getDeregisterMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "Deregister",
      requestType = registry.RegistryOuterClass.Service.class,
      responseType = registry.RegistryOuterClass.EmptyResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<registry.RegistryOuterClass.Service,
      registry.RegistryOuterClass.EmptyResponse> getDeregisterMethod() {
    io.grpc.MethodDescriptor<registry.RegistryOuterClass.Service, registry.RegistryOuterClass.EmptyResponse> getDeregisterMethod;
    if ((getDeregisterMethod = RegistryGrpc.getDeregisterMethod) == null) {
      synchronized (RegistryGrpc.class) {
        if ((getDeregisterMethod = RegistryGrpc.getDeregisterMethod) == null) {
          RegistryGrpc.getDeregisterMethod = getDeregisterMethod =
              io.grpc.MethodDescriptor.<registry.RegistryOuterClass.Service, registry.RegistryOuterClass.EmptyResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "Deregister"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  registry.RegistryOuterClass.Service.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  registry.RegistryOuterClass.EmptyResponse.getDefaultInstance()))
              .setSchemaDescriptor(new RegistryMethodDescriptorSupplier("Deregister"))
              .build();
        }
      }
    }
    return getDeregisterMethod;
  }

  private static volatile io.grpc.MethodDescriptor<registry.RegistryOuterClass.ListRequest,
      registry.RegistryOuterClass.ListResponse> getListServicesMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "ListServices",
      requestType = registry.RegistryOuterClass.ListRequest.class,
      responseType = registry.RegistryOuterClass.ListResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<registry.RegistryOuterClass.ListRequest,
      registry.RegistryOuterClass.ListResponse> getListServicesMethod() {
    io.grpc.MethodDescriptor<registry.RegistryOuterClass.ListRequest, registry.RegistryOuterClass.ListResponse> getListServicesMethod;
    if ((getListServicesMethod = RegistryGrpc.getListServicesMethod) == null) {
      synchronized (RegistryGrpc.class) {
        if ((getListServicesMethod = RegistryGrpc.getListServicesMethod) == null) {
          RegistryGrpc.getListServicesMethod = getListServicesMethod =
              io.grpc.MethodDescriptor.<registry.RegistryOuterClass.ListRequest, registry.RegistryOuterClass.ListResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "ListServices"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  registry.RegistryOuterClass.ListRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  registry.RegistryOuterClass.ListResponse.getDefaultInstance()))
              .setSchemaDescriptor(new RegistryMethodDescriptorSupplier("ListServices"))
              .build();
        }
      }
    }
    return getListServicesMethod;
  }

  private static volatile io.grpc.MethodDescriptor<registry.RegistryOuterClass.WatchRequest,
      registry.RegistryOuterClass.Result> getWatchMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "Watch",
      requestType = registry.RegistryOuterClass.WatchRequest.class,
      responseType = registry.RegistryOuterClass.Result.class,
      methodType = io.grpc.MethodDescriptor.MethodType.SERVER_STREAMING)
  public static io.grpc.MethodDescriptor<registry.RegistryOuterClass.WatchRequest,
      registry.RegistryOuterClass.Result> getWatchMethod() {
    io.grpc.MethodDescriptor<registry.RegistryOuterClass.WatchRequest, registry.RegistryOuterClass.Result> getWatchMethod;
    if ((getWatchMethod = RegistryGrpc.getWatchMethod) == null) {
      synchronized (RegistryGrpc.class) {
        if ((getWatchMethod = RegistryGrpc.getWatchMethod) == null) {
          RegistryGrpc.getWatchMethod = getWatchMethod =
              io.grpc.MethodDescriptor.<registry.RegistryOuterClass.WatchRequest, registry.RegistryOuterClass.Result>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.SERVER_STREAMING)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "Watch"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  registry.RegistryOuterClass.WatchRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  registry.RegistryOuterClass.Result.getDefaultInstance()))
              .setSchemaDescriptor(new RegistryMethodDescriptorSupplier("Watch"))
              .build();
        }
      }
    }
    return getWatchMethod;
  }

  /**
   * Creates a new async stub that supports all call types for the service
   */
  public static RegistryStub newStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<RegistryStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<RegistryStub>() {
        @java.lang.Override
        public RegistryStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new RegistryStub(channel, callOptions);
        }
      };
    return RegistryStub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static RegistryBlockingStub newBlockingStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<RegistryBlockingStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<RegistryBlockingStub>() {
        @java.lang.Override
        public RegistryBlockingStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new RegistryBlockingStub(channel, callOptions);
        }
      };
    return RegistryBlockingStub.newStub(factory, channel);
  }

  /**
   * Creates a new ListenableFuture-style stub that supports unary calls on the service
   */
  public static RegistryFutureStub newFutureStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<RegistryFutureStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<RegistryFutureStub>() {
        @java.lang.Override
        public RegistryFutureStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new RegistryFutureStub(channel, callOptions);
        }
      };
    return RegistryFutureStub.newStub(factory, channel);
  }

  /**
   */
  public static abstract class RegistryImplBase implements io.grpc.BindableService {

    /**
     */
    public void getService(registry.RegistryOuterClass.GetRequest request,
        io.grpc.stub.StreamObserver<registry.RegistryOuterClass.GetResponse> responseObserver) {
      asyncUnimplementedUnaryCall(getGetServiceMethod(), responseObserver);
    }

    /**
     */
    public void register(registry.RegistryOuterClass.Service request,
        io.grpc.stub.StreamObserver<registry.RegistryOuterClass.EmptyResponse> responseObserver) {
      asyncUnimplementedUnaryCall(getRegisterMethod(), responseObserver);
    }

    /**
     */
    public void deregister(registry.RegistryOuterClass.Service request,
        io.grpc.stub.StreamObserver<registry.RegistryOuterClass.EmptyResponse> responseObserver) {
      asyncUnimplementedUnaryCall(getDeregisterMethod(), responseObserver);
    }

    /**
     */
    public void listServices(registry.RegistryOuterClass.ListRequest request,
        io.grpc.stub.StreamObserver<registry.RegistryOuterClass.ListResponse> responseObserver) {
      asyncUnimplementedUnaryCall(getListServicesMethod(), responseObserver);
    }

    /**
     */
    public void watch(registry.RegistryOuterClass.WatchRequest request,
        io.grpc.stub.StreamObserver<registry.RegistryOuterClass.Result> responseObserver) {
      asyncUnimplementedUnaryCall(getWatchMethod(), responseObserver);
    }

    @java.lang.Override public final io.grpc.ServerServiceDefinition bindService() {
      return io.grpc.ServerServiceDefinition.builder(getServiceDescriptor())
          .addMethod(
            getGetServiceMethod(),
            asyncUnaryCall(
              new MethodHandlers<
                registry.RegistryOuterClass.GetRequest,
                registry.RegistryOuterClass.GetResponse>(
                  this, METHODID_GET_SERVICE)))
          .addMethod(
            getRegisterMethod(),
            asyncUnaryCall(
              new MethodHandlers<
                registry.RegistryOuterClass.Service,
                registry.RegistryOuterClass.EmptyResponse>(
                  this, METHODID_REGISTER)))
          .addMethod(
            getDeregisterMethod(),
            asyncUnaryCall(
              new MethodHandlers<
                registry.RegistryOuterClass.Service,
                registry.RegistryOuterClass.EmptyResponse>(
                  this, METHODID_DEREGISTER)))
          .addMethod(
            getListServicesMethod(),
            asyncUnaryCall(
              new MethodHandlers<
                registry.RegistryOuterClass.ListRequest,
                registry.RegistryOuterClass.ListResponse>(
                  this, METHODID_LIST_SERVICES)))
          .addMethod(
            getWatchMethod(),
            asyncServerStreamingCall(
              new MethodHandlers<
                registry.RegistryOuterClass.WatchRequest,
                registry.RegistryOuterClass.Result>(
                  this, METHODID_WATCH)))
          .build();
    }
  }

  /**
   */
  public static final class RegistryStub extends io.grpc.stub.AbstractAsyncStub<RegistryStub> {
    private RegistryStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected RegistryStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new RegistryStub(channel, callOptions);
    }

    /**
     */
    public void getService(registry.RegistryOuterClass.GetRequest request,
        io.grpc.stub.StreamObserver<registry.RegistryOuterClass.GetResponse> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(getGetServiceMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void register(registry.RegistryOuterClass.Service request,
        io.grpc.stub.StreamObserver<registry.RegistryOuterClass.EmptyResponse> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(getRegisterMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void deregister(registry.RegistryOuterClass.Service request,
        io.grpc.stub.StreamObserver<registry.RegistryOuterClass.EmptyResponse> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(getDeregisterMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void listServices(registry.RegistryOuterClass.ListRequest request,
        io.grpc.stub.StreamObserver<registry.RegistryOuterClass.ListResponse> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(getListServicesMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void watch(registry.RegistryOuterClass.WatchRequest request,
        io.grpc.stub.StreamObserver<registry.RegistryOuterClass.Result> responseObserver) {
      asyncServerStreamingCall(
          getChannel().newCall(getWatchMethod(), getCallOptions()), request, responseObserver);
    }
  }

  /**
   */
  public static final class RegistryBlockingStub extends io.grpc.stub.AbstractBlockingStub<RegistryBlockingStub> {
    private RegistryBlockingStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected RegistryBlockingStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new RegistryBlockingStub(channel, callOptions);
    }

    /**
     */
    public registry.RegistryOuterClass.GetResponse getService(registry.RegistryOuterClass.GetRequest request) {
      return blockingUnaryCall(
          getChannel(), getGetServiceMethod(), getCallOptions(), request);
    }

    /**
     */
    public registry.RegistryOuterClass.EmptyResponse register(registry.RegistryOuterClass.Service request) {
      return blockingUnaryCall(
          getChannel(), getRegisterMethod(), getCallOptions(), request);
    }

    /**
     */
    public registry.RegistryOuterClass.EmptyResponse deregister(registry.RegistryOuterClass.Service request) {
      return blockingUnaryCall(
          getChannel(), getDeregisterMethod(), getCallOptions(), request);
    }

    /**
     */
    public registry.RegistryOuterClass.ListResponse listServices(registry.RegistryOuterClass.ListRequest request) {
      return blockingUnaryCall(
          getChannel(), getListServicesMethod(), getCallOptions(), request);
    }

    /**
     */
    public java.util.Iterator<registry.RegistryOuterClass.Result> watch(
        registry.RegistryOuterClass.WatchRequest request) {
      return blockingServerStreamingCall(
          getChannel(), getWatchMethod(), getCallOptions(), request);
    }
  }

  /**
   */
  public static final class RegistryFutureStub extends io.grpc.stub.AbstractFutureStub<RegistryFutureStub> {
    private RegistryFutureStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected RegistryFutureStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new RegistryFutureStub(channel, callOptions);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<registry.RegistryOuterClass.GetResponse> getService(
        registry.RegistryOuterClass.GetRequest request) {
      return futureUnaryCall(
          getChannel().newCall(getGetServiceMethod(), getCallOptions()), request);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<registry.RegistryOuterClass.EmptyResponse> register(
        registry.RegistryOuterClass.Service request) {
      return futureUnaryCall(
          getChannel().newCall(getRegisterMethod(), getCallOptions()), request);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<registry.RegistryOuterClass.EmptyResponse> deregister(
        registry.RegistryOuterClass.Service request) {
      return futureUnaryCall(
          getChannel().newCall(getDeregisterMethod(), getCallOptions()), request);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<registry.RegistryOuterClass.ListResponse> listServices(
        registry.RegistryOuterClass.ListRequest request) {
      return futureUnaryCall(
          getChannel().newCall(getListServicesMethod(), getCallOptions()), request);
    }
  }

  private static final int METHODID_GET_SERVICE = 0;
  private static final int METHODID_REGISTER = 1;
  private static final int METHODID_DEREGISTER = 2;
  private static final int METHODID_LIST_SERVICES = 3;
  private static final int METHODID_WATCH = 4;

  private static final class MethodHandlers<Req, Resp> implements
      io.grpc.stub.ServerCalls.UnaryMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ServerStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ClientStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.BidiStreamingMethod<Req, Resp> {
    private final RegistryImplBase serviceImpl;
    private final int methodId;

    MethodHandlers(RegistryImplBase serviceImpl, int methodId) {
      this.serviceImpl = serviceImpl;
      this.methodId = methodId;
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public void invoke(Req request, io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        case METHODID_GET_SERVICE:
          serviceImpl.getService((registry.RegistryOuterClass.GetRequest) request,
              (io.grpc.stub.StreamObserver<registry.RegistryOuterClass.GetResponse>) responseObserver);
          break;
        case METHODID_REGISTER:
          serviceImpl.register((registry.RegistryOuterClass.Service) request,
              (io.grpc.stub.StreamObserver<registry.RegistryOuterClass.EmptyResponse>) responseObserver);
          break;
        case METHODID_DEREGISTER:
          serviceImpl.deregister((registry.RegistryOuterClass.Service) request,
              (io.grpc.stub.StreamObserver<registry.RegistryOuterClass.EmptyResponse>) responseObserver);
          break;
        case METHODID_LIST_SERVICES:
          serviceImpl.listServices((registry.RegistryOuterClass.ListRequest) request,
              (io.grpc.stub.StreamObserver<registry.RegistryOuterClass.ListResponse>) responseObserver);
          break;
        case METHODID_WATCH:
          serviceImpl.watch((registry.RegistryOuterClass.WatchRequest) request,
              (io.grpc.stub.StreamObserver<registry.RegistryOuterClass.Result>) responseObserver);
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

  private static abstract class RegistryBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier, io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    RegistryBaseDescriptorSupplier() {}

    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return registry.RegistryOuterClass.getDescriptor();
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.ServiceDescriptor getServiceDescriptor() {
      return getFileDescriptor().findServiceByName("Registry");
    }
  }

  private static final class RegistryFileDescriptorSupplier
      extends RegistryBaseDescriptorSupplier {
    RegistryFileDescriptorSupplier() {}
  }

  private static final class RegistryMethodDescriptorSupplier
      extends RegistryBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoMethodDescriptorSupplier {
    private final String methodName;

    RegistryMethodDescriptorSupplier(String methodName) {
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
      synchronized (RegistryGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor = result = io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
              .setSchemaDescriptor(new RegistryFileDescriptorSupplier())
              .addMethod(getGetServiceMethod())
              .addMethod(getRegisterMethod())
              .addMethod(getDeregisterMethod())
              .addMethod(getListServicesMethod())
              .addMethod(getWatchMethod())
              .build();
        }
      }
    }
    return result;
  }
}
