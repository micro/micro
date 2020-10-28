package api;

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
    comments = "Source: api/api.proto")
public final class ApiGrpc {

  private ApiGrpc() {}

  public static final String SERVICE_NAME = "api.Api";

  // Static method descriptors that strictly reflect the proto.
  private static volatile io.grpc.MethodDescriptor<api.ApiOuterClass.Endpoint,
      api.ApiOuterClass.EmptyResponse> getRegisterMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "Register",
      requestType = api.ApiOuterClass.Endpoint.class,
      responseType = api.ApiOuterClass.EmptyResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<api.ApiOuterClass.Endpoint,
      api.ApiOuterClass.EmptyResponse> getRegisterMethod() {
    io.grpc.MethodDescriptor<api.ApiOuterClass.Endpoint, api.ApiOuterClass.EmptyResponse> getRegisterMethod;
    if ((getRegisterMethod = ApiGrpc.getRegisterMethod) == null) {
      synchronized (ApiGrpc.class) {
        if ((getRegisterMethod = ApiGrpc.getRegisterMethod) == null) {
          ApiGrpc.getRegisterMethod = getRegisterMethod =
              io.grpc.MethodDescriptor.<api.ApiOuterClass.Endpoint, api.ApiOuterClass.EmptyResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "Register"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  api.ApiOuterClass.Endpoint.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  api.ApiOuterClass.EmptyResponse.getDefaultInstance()))
              .setSchemaDescriptor(new ApiMethodDescriptorSupplier("Register"))
              .build();
        }
      }
    }
    return getRegisterMethod;
  }

  private static volatile io.grpc.MethodDescriptor<api.ApiOuterClass.Endpoint,
      api.ApiOuterClass.EmptyResponse> getDeregisterMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "Deregister",
      requestType = api.ApiOuterClass.Endpoint.class,
      responseType = api.ApiOuterClass.EmptyResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<api.ApiOuterClass.Endpoint,
      api.ApiOuterClass.EmptyResponse> getDeregisterMethod() {
    io.grpc.MethodDescriptor<api.ApiOuterClass.Endpoint, api.ApiOuterClass.EmptyResponse> getDeregisterMethod;
    if ((getDeregisterMethod = ApiGrpc.getDeregisterMethod) == null) {
      synchronized (ApiGrpc.class) {
        if ((getDeregisterMethod = ApiGrpc.getDeregisterMethod) == null) {
          ApiGrpc.getDeregisterMethod = getDeregisterMethod =
              io.grpc.MethodDescriptor.<api.ApiOuterClass.Endpoint, api.ApiOuterClass.EmptyResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "Deregister"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  api.ApiOuterClass.Endpoint.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  api.ApiOuterClass.EmptyResponse.getDefaultInstance()))
              .setSchemaDescriptor(new ApiMethodDescriptorSupplier("Deregister"))
              .build();
        }
      }
    }
    return getDeregisterMethod;
  }

  /**
   * Creates a new async stub that supports all call types for the service
   */
  public static ApiStub newStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<ApiStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<ApiStub>() {
        @java.lang.Override
        public ApiStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new ApiStub(channel, callOptions);
        }
      };
    return ApiStub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static ApiBlockingStub newBlockingStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<ApiBlockingStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<ApiBlockingStub>() {
        @java.lang.Override
        public ApiBlockingStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new ApiBlockingStub(channel, callOptions);
        }
      };
    return ApiBlockingStub.newStub(factory, channel);
  }

  /**
   * Creates a new ListenableFuture-style stub that supports unary calls on the service
   */
  public static ApiFutureStub newFutureStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<ApiFutureStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<ApiFutureStub>() {
        @java.lang.Override
        public ApiFutureStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new ApiFutureStub(channel, callOptions);
        }
      };
    return ApiFutureStub.newStub(factory, channel);
  }

  /**
   */
  public static abstract class ApiImplBase implements io.grpc.BindableService {

    /**
     */
    public void register(api.ApiOuterClass.Endpoint request,
        io.grpc.stub.StreamObserver<api.ApiOuterClass.EmptyResponse> responseObserver) {
      asyncUnimplementedUnaryCall(getRegisterMethod(), responseObserver);
    }

    /**
     */
    public void deregister(api.ApiOuterClass.Endpoint request,
        io.grpc.stub.StreamObserver<api.ApiOuterClass.EmptyResponse> responseObserver) {
      asyncUnimplementedUnaryCall(getDeregisterMethod(), responseObserver);
    }

    @java.lang.Override public final io.grpc.ServerServiceDefinition bindService() {
      return io.grpc.ServerServiceDefinition.builder(getServiceDescriptor())
          .addMethod(
            getRegisterMethod(),
            asyncUnaryCall(
              new MethodHandlers<
                api.ApiOuterClass.Endpoint,
                api.ApiOuterClass.EmptyResponse>(
                  this, METHODID_REGISTER)))
          .addMethod(
            getDeregisterMethod(),
            asyncUnaryCall(
              new MethodHandlers<
                api.ApiOuterClass.Endpoint,
                api.ApiOuterClass.EmptyResponse>(
                  this, METHODID_DEREGISTER)))
          .build();
    }
  }

  /**
   */
  public static final class ApiStub extends io.grpc.stub.AbstractAsyncStub<ApiStub> {
    private ApiStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected ApiStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new ApiStub(channel, callOptions);
    }

    /**
     */
    public void register(api.ApiOuterClass.Endpoint request,
        io.grpc.stub.StreamObserver<api.ApiOuterClass.EmptyResponse> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(getRegisterMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void deregister(api.ApiOuterClass.Endpoint request,
        io.grpc.stub.StreamObserver<api.ApiOuterClass.EmptyResponse> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(getDeregisterMethod(), getCallOptions()), request, responseObserver);
    }
  }

  /**
   */
  public static final class ApiBlockingStub extends io.grpc.stub.AbstractBlockingStub<ApiBlockingStub> {
    private ApiBlockingStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected ApiBlockingStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new ApiBlockingStub(channel, callOptions);
    }

    /**
     */
    public api.ApiOuterClass.EmptyResponse register(api.ApiOuterClass.Endpoint request) {
      return blockingUnaryCall(
          getChannel(), getRegisterMethod(), getCallOptions(), request);
    }

    /**
     */
    public api.ApiOuterClass.EmptyResponse deregister(api.ApiOuterClass.Endpoint request) {
      return blockingUnaryCall(
          getChannel(), getDeregisterMethod(), getCallOptions(), request);
    }
  }

  /**
   */
  public static final class ApiFutureStub extends io.grpc.stub.AbstractFutureStub<ApiFutureStub> {
    private ApiFutureStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected ApiFutureStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new ApiFutureStub(channel, callOptions);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<api.ApiOuterClass.EmptyResponse> register(
        api.ApiOuterClass.Endpoint request) {
      return futureUnaryCall(
          getChannel().newCall(getRegisterMethod(), getCallOptions()), request);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<api.ApiOuterClass.EmptyResponse> deregister(
        api.ApiOuterClass.Endpoint request) {
      return futureUnaryCall(
          getChannel().newCall(getDeregisterMethod(), getCallOptions()), request);
    }
  }

  private static final int METHODID_REGISTER = 0;
  private static final int METHODID_DEREGISTER = 1;

  private static final class MethodHandlers<Req, Resp> implements
      io.grpc.stub.ServerCalls.UnaryMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ServerStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ClientStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.BidiStreamingMethod<Req, Resp> {
    private final ApiImplBase serviceImpl;
    private final int methodId;

    MethodHandlers(ApiImplBase serviceImpl, int methodId) {
      this.serviceImpl = serviceImpl;
      this.methodId = methodId;
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public void invoke(Req request, io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        case METHODID_REGISTER:
          serviceImpl.register((api.ApiOuterClass.Endpoint) request,
              (io.grpc.stub.StreamObserver<api.ApiOuterClass.EmptyResponse>) responseObserver);
          break;
        case METHODID_DEREGISTER:
          serviceImpl.deregister((api.ApiOuterClass.Endpoint) request,
              (io.grpc.stub.StreamObserver<api.ApiOuterClass.EmptyResponse>) responseObserver);
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

  private static abstract class ApiBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier, io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    ApiBaseDescriptorSupplier() {}

    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return api.ApiOuterClass.getDescriptor();
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.ServiceDescriptor getServiceDescriptor() {
      return getFileDescriptor().findServiceByName("Api");
    }
  }

  private static final class ApiFileDescriptorSupplier
      extends ApiBaseDescriptorSupplier {
    ApiFileDescriptorSupplier() {}
  }

  private static final class ApiMethodDescriptorSupplier
      extends ApiBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoMethodDescriptorSupplier {
    private final String methodName;

    ApiMethodDescriptorSupplier(String methodName) {
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
      synchronized (ApiGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor = result = io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
              .setSchemaDescriptor(new ApiFileDescriptorSupplier())
              .addMethod(getRegisterMethod())
              .addMethod(getDeregisterMethod())
              .build();
        }
      }
    }
    return result;
  }
}
