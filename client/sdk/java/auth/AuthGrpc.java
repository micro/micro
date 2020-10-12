package auth;

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
    comments = "Source: auth/auth.proto")
public final class AuthGrpc {

  private AuthGrpc() {}

  public static final String SERVICE_NAME = "auth.Auth";

  // Static method descriptors that strictly reflect the proto.
  private static volatile io.grpc.MethodDescriptor<auth.AuthOuterClass.GenerateRequest,
      auth.AuthOuterClass.GenerateResponse> getGenerateMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "Generate",
      requestType = auth.AuthOuterClass.GenerateRequest.class,
      responseType = auth.AuthOuterClass.GenerateResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<auth.AuthOuterClass.GenerateRequest,
      auth.AuthOuterClass.GenerateResponse> getGenerateMethod() {
    io.grpc.MethodDescriptor<auth.AuthOuterClass.GenerateRequest, auth.AuthOuterClass.GenerateResponse> getGenerateMethod;
    if ((getGenerateMethod = AuthGrpc.getGenerateMethod) == null) {
      synchronized (AuthGrpc.class) {
        if ((getGenerateMethod = AuthGrpc.getGenerateMethod) == null) {
          AuthGrpc.getGenerateMethod = getGenerateMethod =
              io.grpc.MethodDescriptor.<auth.AuthOuterClass.GenerateRequest, auth.AuthOuterClass.GenerateResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "Generate"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  auth.AuthOuterClass.GenerateRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  auth.AuthOuterClass.GenerateResponse.getDefaultInstance()))
              .setSchemaDescriptor(new AuthMethodDescriptorSupplier("Generate"))
              .build();
        }
      }
    }
    return getGenerateMethod;
  }

  private static volatile io.grpc.MethodDescriptor<auth.AuthOuterClass.InspectRequest,
      auth.AuthOuterClass.InspectResponse> getInspectMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "Inspect",
      requestType = auth.AuthOuterClass.InspectRequest.class,
      responseType = auth.AuthOuterClass.InspectResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<auth.AuthOuterClass.InspectRequest,
      auth.AuthOuterClass.InspectResponse> getInspectMethod() {
    io.grpc.MethodDescriptor<auth.AuthOuterClass.InspectRequest, auth.AuthOuterClass.InspectResponse> getInspectMethod;
    if ((getInspectMethod = AuthGrpc.getInspectMethod) == null) {
      synchronized (AuthGrpc.class) {
        if ((getInspectMethod = AuthGrpc.getInspectMethod) == null) {
          AuthGrpc.getInspectMethod = getInspectMethod =
              io.grpc.MethodDescriptor.<auth.AuthOuterClass.InspectRequest, auth.AuthOuterClass.InspectResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "Inspect"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  auth.AuthOuterClass.InspectRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  auth.AuthOuterClass.InspectResponse.getDefaultInstance()))
              .setSchemaDescriptor(new AuthMethodDescriptorSupplier("Inspect"))
              .build();
        }
      }
    }
    return getInspectMethod;
  }

  private static volatile io.grpc.MethodDescriptor<auth.AuthOuterClass.TokenRequest,
      auth.AuthOuterClass.TokenResponse> getTokenMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "Token",
      requestType = auth.AuthOuterClass.TokenRequest.class,
      responseType = auth.AuthOuterClass.TokenResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<auth.AuthOuterClass.TokenRequest,
      auth.AuthOuterClass.TokenResponse> getTokenMethod() {
    io.grpc.MethodDescriptor<auth.AuthOuterClass.TokenRequest, auth.AuthOuterClass.TokenResponse> getTokenMethod;
    if ((getTokenMethod = AuthGrpc.getTokenMethod) == null) {
      synchronized (AuthGrpc.class) {
        if ((getTokenMethod = AuthGrpc.getTokenMethod) == null) {
          AuthGrpc.getTokenMethod = getTokenMethod =
              io.grpc.MethodDescriptor.<auth.AuthOuterClass.TokenRequest, auth.AuthOuterClass.TokenResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "Token"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  auth.AuthOuterClass.TokenRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  auth.AuthOuterClass.TokenResponse.getDefaultInstance()))
              .setSchemaDescriptor(new AuthMethodDescriptorSupplier("Token"))
              .build();
        }
      }
    }
    return getTokenMethod;
  }

  /**
   * Creates a new async stub that supports all call types for the service
   */
  public static AuthStub newStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<AuthStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<AuthStub>() {
        @java.lang.Override
        public AuthStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new AuthStub(channel, callOptions);
        }
      };
    return AuthStub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static AuthBlockingStub newBlockingStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<AuthBlockingStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<AuthBlockingStub>() {
        @java.lang.Override
        public AuthBlockingStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new AuthBlockingStub(channel, callOptions);
        }
      };
    return AuthBlockingStub.newStub(factory, channel);
  }

  /**
   * Creates a new ListenableFuture-style stub that supports unary calls on the service
   */
  public static AuthFutureStub newFutureStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<AuthFutureStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<AuthFutureStub>() {
        @java.lang.Override
        public AuthFutureStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new AuthFutureStub(channel, callOptions);
        }
      };
    return AuthFutureStub.newStub(factory, channel);
  }

  /**
   */
  public static abstract class AuthImplBase implements io.grpc.BindableService {

    /**
     */
    public void generate(auth.AuthOuterClass.GenerateRequest request,
        io.grpc.stub.StreamObserver<auth.AuthOuterClass.GenerateResponse> responseObserver) {
      asyncUnimplementedUnaryCall(getGenerateMethod(), responseObserver);
    }

    /**
     */
    public void inspect(auth.AuthOuterClass.InspectRequest request,
        io.grpc.stub.StreamObserver<auth.AuthOuterClass.InspectResponse> responseObserver) {
      asyncUnimplementedUnaryCall(getInspectMethod(), responseObserver);
    }

    /**
     */
    public void token(auth.AuthOuterClass.TokenRequest request,
        io.grpc.stub.StreamObserver<auth.AuthOuterClass.TokenResponse> responseObserver) {
      asyncUnimplementedUnaryCall(getTokenMethod(), responseObserver);
    }

    @java.lang.Override public final io.grpc.ServerServiceDefinition bindService() {
      return io.grpc.ServerServiceDefinition.builder(getServiceDescriptor())
          .addMethod(
            getGenerateMethod(),
            asyncUnaryCall(
              new MethodHandlers<
                auth.AuthOuterClass.GenerateRequest,
                auth.AuthOuterClass.GenerateResponse>(
                  this, METHODID_GENERATE)))
          .addMethod(
            getInspectMethod(),
            asyncUnaryCall(
              new MethodHandlers<
                auth.AuthOuterClass.InspectRequest,
                auth.AuthOuterClass.InspectResponse>(
                  this, METHODID_INSPECT)))
          .addMethod(
            getTokenMethod(),
            asyncUnaryCall(
              new MethodHandlers<
                auth.AuthOuterClass.TokenRequest,
                auth.AuthOuterClass.TokenResponse>(
                  this, METHODID_TOKEN)))
          .build();
    }
  }

  /**
   */
  public static final class AuthStub extends io.grpc.stub.AbstractAsyncStub<AuthStub> {
    private AuthStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected AuthStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new AuthStub(channel, callOptions);
    }

    /**
     */
    public void generate(auth.AuthOuterClass.GenerateRequest request,
        io.grpc.stub.StreamObserver<auth.AuthOuterClass.GenerateResponse> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(getGenerateMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void inspect(auth.AuthOuterClass.InspectRequest request,
        io.grpc.stub.StreamObserver<auth.AuthOuterClass.InspectResponse> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(getInspectMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void token(auth.AuthOuterClass.TokenRequest request,
        io.grpc.stub.StreamObserver<auth.AuthOuterClass.TokenResponse> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(getTokenMethod(), getCallOptions()), request, responseObserver);
    }
  }

  /**
   */
  public static final class AuthBlockingStub extends io.grpc.stub.AbstractBlockingStub<AuthBlockingStub> {
    private AuthBlockingStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected AuthBlockingStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new AuthBlockingStub(channel, callOptions);
    }

    /**
     */
    public auth.AuthOuterClass.GenerateResponse generate(auth.AuthOuterClass.GenerateRequest request) {
      return blockingUnaryCall(
          getChannel(), getGenerateMethod(), getCallOptions(), request);
    }

    /**
     */
    public auth.AuthOuterClass.InspectResponse inspect(auth.AuthOuterClass.InspectRequest request) {
      return blockingUnaryCall(
          getChannel(), getInspectMethod(), getCallOptions(), request);
    }

    /**
     */
    public auth.AuthOuterClass.TokenResponse token(auth.AuthOuterClass.TokenRequest request) {
      return blockingUnaryCall(
          getChannel(), getTokenMethod(), getCallOptions(), request);
    }
  }

  /**
   */
  public static final class AuthFutureStub extends io.grpc.stub.AbstractFutureStub<AuthFutureStub> {
    private AuthFutureStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected AuthFutureStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new AuthFutureStub(channel, callOptions);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<auth.AuthOuterClass.GenerateResponse> generate(
        auth.AuthOuterClass.GenerateRequest request) {
      return futureUnaryCall(
          getChannel().newCall(getGenerateMethod(), getCallOptions()), request);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<auth.AuthOuterClass.InspectResponse> inspect(
        auth.AuthOuterClass.InspectRequest request) {
      return futureUnaryCall(
          getChannel().newCall(getInspectMethod(), getCallOptions()), request);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<auth.AuthOuterClass.TokenResponse> token(
        auth.AuthOuterClass.TokenRequest request) {
      return futureUnaryCall(
          getChannel().newCall(getTokenMethod(), getCallOptions()), request);
    }
  }

  private static final int METHODID_GENERATE = 0;
  private static final int METHODID_INSPECT = 1;
  private static final int METHODID_TOKEN = 2;

  private static final class MethodHandlers<Req, Resp> implements
      io.grpc.stub.ServerCalls.UnaryMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ServerStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ClientStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.BidiStreamingMethod<Req, Resp> {
    private final AuthImplBase serviceImpl;
    private final int methodId;

    MethodHandlers(AuthImplBase serviceImpl, int methodId) {
      this.serviceImpl = serviceImpl;
      this.methodId = methodId;
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public void invoke(Req request, io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        case METHODID_GENERATE:
          serviceImpl.generate((auth.AuthOuterClass.GenerateRequest) request,
              (io.grpc.stub.StreamObserver<auth.AuthOuterClass.GenerateResponse>) responseObserver);
          break;
        case METHODID_INSPECT:
          serviceImpl.inspect((auth.AuthOuterClass.InspectRequest) request,
              (io.grpc.stub.StreamObserver<auth.AuthOuterClass.InspectResponse>) responseObserver);
          break;
        case METHODID_TOKEN:
          serviceImpl.token((auth.AuthOuterClass.TokenRequest) request,
              (io.grpc.stub.StreamObserver<auth.AuthOuterClass.TokenResponse>) responseObserver);
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

  private static abstract class AuthBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier, io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    AuthBaseDescriptorSupplier() {}

    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return auth.AuthOuterClass.getDescriptor();
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.ServiceDescriptor getServiceDescriptor() {
      return getFileDescriptor().findServiceByName("Auth");
    }
  }

  private static final class AuthFileDescriptorSupplier
      extends AuthBaseDescriptorSupplier {
    AuthFileDescriptorSupplier() {}
  }

  private static final class AuthMethodDescriptorSupplier
      extends AuthBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoMethodDescriptorSupplier {
    private final String methodName;

    AuthMethodDescriptorSupplier(String methodName) {
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
      synchronized (AuthGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor = result = io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
              .setSchemaDescriptor(new AuthFileDescriptorSupplier())
              .addMethod(getGenerateMethod())
              .addMethod(getInspectMethod())
              .addMethod(getTokenMethod())
              .build();
        }
      }
    }
    return result;
  }
}
