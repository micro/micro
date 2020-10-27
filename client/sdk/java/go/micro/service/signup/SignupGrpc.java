package go.micro.service.signup;

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
    comments = "Source: signup/signup.proto")
public final class SignupGrpc {

  private SignupGrpc() {}

  public static final String SERVICE_NAME = "go.micro.service.signup.Signup";

  // Static method descriptors that strictly reflect the proto.
  private static volatile io.grpc.MethodDescriptor<go.micro.service.signup.SignupOuterClass.SendVerificationEmailRequest,
      go.micro.service.signup.SignupOuterClass.SendVerificationEmailResponse> getSendVerificationEmailMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "SendVerificationEmail",
      requestType = go.micro.service.signup.SignupOuterClass.SendVerificationEmailRequest.class,
      responseType = go.micro.service.signup.SignupOuterClass.SendVerificationEmailResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<go.micro.service.signup.SignupOuterClass.SendVerificationEmailRequest,
      go.micro.service.signup.SignupOuterClass.SendVerificationEmailResponse> getSendVerificationEmailMethod() {
    io.grpc.MethodDescriptor<go.micro.service.signup.SignupOuterClass.SendVerificationEmailRequest, go.micro.service.signup.SignupOuterClass.SendVerificationEmailResponse> getSendVerificationEmailMethod;
    if ((getSendVerificationEmailMethod = SignupGrpc.getSendVerificationEmailMethod) == null) {
      synchronized (SignupGrpc.class) {
        if ((getSendVerificationEmailMethod = SignupGrpc.getSendVerificationEmailMethod) == null) {
          SignupGrpc.getSendVerificationEmailMethod = getSendVerificationEmailMethod =
              io.grpc.MethodDescriptor.<go.micro.service.signup.SignupOuterClass.SendVerificationEmailRequest, go.micro.service.signup.SignupOuterClass.SendVerificationEmailResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "SendVerificationEmail"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  go.micro.service.signup.SignupOuterClass.SendVerificationEmailRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  go.micro.service.signup.SignupOuterClass.SendVerificationEmailResponse.getDefaultInstance()))
              .setSchemaDescriptor(new SignupMethodDescriptorSupplier("SendVerificationEmail"))
              .build();
        }
      }
    }
    return getSendVerificationEmailMethod;
  }

  private static volatile io.grpc.MethodDescriptor<go.micro.service.signup.SignupOuterClass.VerifyRequest,
      go.micro.service.signup.SignupOuterClass.VerifyResponse> getVerifyMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "Verify",
      requestType = go.micro.service.signup.SignupOuterClass.VerifyRequest.class,
      responseType = go.micro.service.signup.SignupOuterClass.VerifyResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<go.micro.service.signup.SignupOuterClass.VerifyRequest,
      go.micro.service.signup.SignupOuterClass.VerifyResponse> getVerifyMethod() {
    io.grpc.MethodDescriptor<go.micro.service.signup.SignupOuterClass.VerifyRequest, go.micro.service.signup.SignupOuterClass.VerifyResponse> getVerifyMethod;
    if ((getVerifyMethod = SignupGrpc.getVerifyMethod) == null) {
      synchronized (SignupGrpc.class) {
        if ((getVerifyMethod = SignupGrpc.getVerifyMethod) == null) {
          SignupGrpc.getVerifyMethod = getVerifyMethod =
              io.grpc.MethodDescriptor.<go.micro.service.signup.SignupOuterClass.VerifyRequest, go.micro.service.signup.SignupOuterClass.VerifyResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "Verify"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  go.micro.service.signup.SignupOuterClass.VerifyRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  go.micro.service.signup.SignupOuterClass.VerifyResponse.getDefaultInstance()))
              .setSchemaDescriptor(new SignupMethodDescriptorSupplier("Verify"))
              .build();
        }
      }
    }
    return getVerifyMethod;
  }

  private static volatile io.grpc.MethodDescriptor<go.micro.service.signup.SignupOuterClass.SetPaymentMethodRequest,
      go.micro.service.signup.SignupOuterClass.SetPaymentMethodResponse> getSetPaymentMethodMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "SetPaymentMethod",
      requestType = go.micro.service.signup.SignupOuterClass.SetPaymentMethodRequest.class,
      responseType = go.micro.service.signup.SignupOuterClass.SetPaymentMethodResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<go.micro.service.signup.SignupOuterClass.SetPaymentMethodRequest,
      go.micro.service.signup.SignupOuterClass.SetPaymentMethodResponse> getSetPaymentMethodMethod() {
    io.grpc.MethodDescriptor<go.micro.service.signup.SignupOuterClass.SetPaymentMethodRequest, go.micro.service.signup.SignupOuterClass.SetPaymentMethodResponse> getSetPaymentMethodMethod;
    if ((getSetPaymentMethodMethod = SignupGrpc.getSetPaymentMethodMethod) == null) {
      synchronized (SignupGrpc.class) {
        if ((getSetPaymentMethodMethod = SignupGrpc.getSetPaymentMethodMethod) == null) {
          SignupGrpc.getSetPaymentMethodMethod = getSetPaymentMethodMethod =
              io.grpc.MethodDescriptor.<go.micro.service.signup.SignupOuterClass.SetPaymentMethodRequest, go.micro.service.signup.SignupOuterClass.SetPaymentMethodResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "SetPaymentMethod"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  go.micro.service.signup.SignupOuterClass.SetPaymentMethodRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  go.micro.service.signup.SignupOuterClass.SetPaymentMethodResponse.getDefaultInstance()))
              .setSchemaDescriptor(new SignupMethodDescriptorSupplier("SetPaymentMethod"))
              .build();
        }
      }
    }
    return getSetPaymentMethodMethod;
  }

  private static volatile io.grpc.MethodDescriptor<go.micro.service.signup.SignupOuterClass.HasPaymentMethodRequest,
      go.micro.service.signup.SignupOuterClass.HasPaymentMethodResponse> getHasPaymentMethodMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "HasPaymentMethod",
      requestType = go.micro.service.signup.SignupOuterClass.HasPaymentMethodRequest.class,
      responseType = go.micro.service.signup.SignupOuterClass.HasPaymentMethodResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<go.micro.service.signup.SignupOuterClass.HasPaymentMethodRequest,
      go.micro.service.signup.SignupOuterClass.HasPaymentMethodResponse> getHasPaymentMethodMethod() {
    io.grpc.MethodDescriptor<go.micro.service.signup.SignupOuterClass.HasPaymentMethodRequest, go.micro.service.signup.SignupOuterClass.HasPaymentMethodResponse> getHasPaymentMethodMethod;
    if ((getHasPaymentMethodMethod = SignupGrpc.getHasPaymentMethodMethod) == null) {
      synchronized (SignupGrpc.class) {
        if ((getHasPaymentMethodMethod = SignupGrpc.getHasPaymentMethodMethod) == null) {
          SignupGrpc.getHasPaymentMethodMethod = getHasPaymentMethodMethod =
              io.grpc.MethodDescriptor.<go.micro.service.signup.SignupOuterClass.HasPaymentMethodRequest, go.micro.service.signup.SignupOuterClass.HasPaymentMethodResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "HasPaymentMethod"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  go.micro.service.signup.SignupOuterClass.HasPaymentMethodRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  go.micro.service.signup.SignupOuterClass.HasPaymentMethodResponse.getDefaultInstance()))
              .setSchemaDescriptor(new SignupMethodDescriptorSupplier("HasPaymentMethod"))
              .build();
        }
      }
    }
    return getHasPaymentMethodMethod;
  }

  private static volatile io.grpc.MethodDescriptor<go.micro.service.signup.SignupOuterClass.CompleteSignupRequest,
      go.micro.service.signup.SignupOuterClass.CompleteSignupResponse> getCompleteSignupMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "CompleteSignup",
      requestType = go.micro.service.signup.SignupOuterClass.CompleteSignupRequest.class,
      responseType = go.micro.service.signup.SignupOuterClass.CompleteSignupResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<go.micro.service.signup.SignupOuterClass.CompleteSignupRequest,
      go.micro.service.signup.SignupOuterClass.CompleteSignupResponse> getCompleteSignupMethod() {
    io.grpc.MethodDescriptor<go.micro.service.signup.SignupOuterClass.CompleteSignupRequest, go.micro.service.signup.SignupOuterClass.CompleteSignupResponse> getCompleteSignupMethod;
    if ((getCompleteSignupMethod = SignupGrpc.getCompleteSignupMethod) == null) {
      synchronized (SignupGrpc.class) {
        if ((getCompleteSignupMethod = SignupGrpc.getCompleteSignupMethod) == null) {
          SignupGrpc.getCompleteSignupMethod = getCompleteSignupMethod =
              io.grpc.MethodDescriptor.<go.micro.service.signup.SignupOuterClass.CompleteSignupRequest, go.micro.service.signup.SignupOuterClass.CompleteSignupResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "CompleteSignup"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  go.micro.service.signup.SignupOuterClass.CompleteSignupRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  go.micro.service.signup.SignupOuterClass.CompleteSignupResponse.getDefaultInstance()))
              .setSchemaDescriptor(new SignupMethodDescriptorSupplier("CompleteSignup"))
              .build();
        }
      }
    }
    return getCompleteSignupMethod;
  }

  private static volatile io.grpc.MethodDescriptor<go.micro.service.signup.SignupOuterClass.RecoverRequest,
      go.micro.service.signup.SignupOuterClass.RecoverResponse> getRecoverMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "Recover",
      requestType = go.micro.service.signup.SignupOuterClass.RecoverRequest.class,
      responseType = go.micro.service.signup.SignupOuterClass.RecoverResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<go.micro.service.signup.SignupOuterClass.RecoverRequest,
      go.micro.service.signup.SignupOuterClass.RecoverResponse> getRecoverMethod() {
    io.grpc.MethodDescriptor<go.micro.service.signup.SignupOuterClass.RecoverRequest, go.micro.service.signup.SignupOuterClass.RecoverResponse> getRecoverMethod;
    if ((getRecoverMethod = SignupGrpc.getRecoverMethod) == null) {
      synchronized (SignupGrpc.class) {
        if ((getRecoverMethod = SignupGrpc.getRecoverMethod) == null) {
          SignupGrpc.getRecoverMethod = getRecoverMethod =
              io.grpc.MethodDescriptor.<go.micro.service.signup.SignupOuterClass.RecoverRequest, go.micro.service.signup.SignupOuterClass.RecoverResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "Recover"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  go.micro.service.signup.SignupOuterClass.RecoverRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  go.micro.service.signup.SignupOuterClass.RecoverResponse.getDefaultInstance()))
              .setSchemaDescriptor(new SignupMethodDescriptorSupplier("Recover"))
              .build();
        }
      }
    }
    return getRecoverMethod;
  }

  /**
   * Creates a new async stub that supports all call types for the service
   */
  public static SignupStub newStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<SignupStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<SignupStub>() {
        @java.lang.Override
        public SignupStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new SignupStub(channel, callOptions);
        }
      };
    return SignupStub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static SignupBlockingStub newBlockingStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<SignupBlockingStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<SignupBlockingStub>() {
        @java.lang.Override
        public SignupBlockingStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new SignupBlockingStub(channel, callOptions);
        }
      };
    return SignupBlockingStub.newStub(factory, channel);
  }

  /**
   * Creates a new ListenableFuture-style stub that supports unary calls on the service
   */
  public static SignupFutureStub newFutureStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<SignupFutureStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<SignupFutureStub>() {
        @java.lang.Override
        public SignupFutureStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new SignupFutureStub(channel, callOptions);
        }
      };
    return SignupFutureStub.newStub(factory, channel);
  }

  /**
   */
  public static abstract class SignupImplBase implements io.grpc.BindableService {

    /**
     * <pre>
     * Sends the verification email to the user
     * </pre>
     */
    public void sendVerificationEmail(go.micro.service.signup.SignupOuterClass.SendVerificationEmailRequest request,
        io.grpc.stub.StreamObserver<go.micro.service.signup.SignupOuterClass.SendVerificationEmailResponse> responseObserver) {
      asyncUnimplementedUnaryCall(getSendVerificationEmailMethod(), responseObserver);
    }

    /**
     * <pre>
     * Verify kicks off the process of verification
     * </pre>
     */
    public void verify(go.micro.service.signup.SignupOuterClass.VerifyRequest request,
        io.grpc.stub.StreamObserver<go.micro.service.signup.SignupOuterClass.VerifyResponse> responseObserver) {
      asyncUnimplementedUnaryCall(getVerifyMethod(), responseObserver);
    }

    /**
     */
    public void setPaymentMethod(go.micro.service.signup.SignupOuterClass.SetPaymentMethodRequest request,
        io.grpc.stub.StreamObserver<go.micro.service.signup.SignupOuterClass.SetPaymentMethodResponse> responseObserver) {
      asyncUnimplementedUnaryCall(getSetPaymentMethodMethod(), responseObserver);
    }

    /**
     */
    public void hasPaymentMethod(go.micro.service.signup.SignupOuterClass.HasPaymentMethodRequest request,
        io.grpc.stub.StreamObserver<go.micro.service.signup.SignupOuterClass.HasPaymentMethodResponse> responseObserver) {
      asyncUnimplementedUnaryCall(getHasPaymentMethodMethod(), responseObserver);
    }

    /**
     * <pre>
     * Creates a subscription and an account
     * </pre>
     */
    public void completeSignup(go.micro.service.signup.SignupOuterClass.CompleteSignupRequest request,
        io.grpc.stub.StreamObserver<go.micro.service.signup.SignupOuterClass.CompleteSignupResponse> responseObserver) {
      asyncUnimplementedUnaryCall(getCompleteSignupMethod(), responseObserver);
    }

    /**
     */
    public void recover(go.micro.service.signup.SignupOuterClass.RecoverRequest request,
        io.grpc.stub.StreamObserver<go.micro.service.signup.SignupOuterClass.RecoverResponse> responseObserver) {
      asyncUnimplementedUnaryCall(getRecoverMethod(), responseObserver);
    }

    @java.lang.Override public final io.grpc.ServerServiceDefinition bindService() {
      return io.grpc.ServerServiceDefinition.builder(getServiceDescriptor())
          .addMethod(
            getSendVerificationEmailMethod(),
            asyncUnaryCall(
              new MethodHandlers<
                go.micro.service.signup.SignupOuterClass.SendVerificationEmailRequest,
                go.micro.service.signup.SignupOuterClass.SendVerificationEmailResponse>(
                  this, METHODID_SEND_VERIFICATION_EMAIL)))
          .addMethod(
            getVerifyMethod(),
            asyncUnaryCall(
              new MethodHandlers<
                go.micro.service.signup.SignupOuterClass.VerifyRequest,
                go.micro.service.signup.SignupOuterClass.VerifyResponse>(
                  this, METHODID_VERIFY)))
          .addMethod(
            getSetPaymentMethodMethod(),
            asyncUnaryCall(
              new MethodHandlers<
                go.micro.service.signup.SignupOuterClass.SetPaymentMethodRequest,
                go.micro.service.signup.SignupOuterClass.SetPaymentMethodResponse>(
                  this, METHODID_SET_PAYMENT_METHOD)))
          .addMethod(
            getHasPaymentMethodMethod(),
            asyncUnaryCall(
              new MethodHandlers<
                go.micro.service.signup.SignupOuterClass.HasPaymentMethodRequest,
                go.micro.service.signup.SignupOuterClass.HasPaymentMethodResponse>(
                  this, METHODID_HAS_PAYMENT_METHOD)))
          .addMethod(
            getCompleteSignupMethod(),
            asyncUnaryCall(
              new MethodHandlers<
                go.micro.service.signup.SignupOuterClass.CompleteSignupRequest,
                go.micro.service.signup.SignupOuterClass.CompleteSignupResponse>(
                  this, METHODID_COMPLETE_SIGNUP)))
          .addMethod(
            getRecoverMethod(),
            asyncUnaryCall(
              new MethodHandlers<
                go.micro.service.signup.SignupOuterClass.RecoverRequest,
                go.micro.service.signup.SignupOuterClass.RecoverResponse>(
                  this, METHODID_RECOVER)))
          .build();
    }
  }

  /**
   */
  public static final class SignupStub extends io.grpc.stub.AbstractAsyncStub<SignupStub> {
    private SignupStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected SignupStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new SignupStub(channel, callOptions);
    }

    /**
     * <pre>
     * Sends the verification email to the user
     * </pre>
     */
    public void sendVerificationEmail(go.micro.service.signup.SignupOuterClass.SendVerificationEmailRequest request,
        io.grpc.stub.StreamObserver<go.micro.service.signup.SignupOuterClass.SendVerificationEmailResponse> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(getSendVerificationEmailMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * Verify kicks off the process of verification
     * </pre>
     */
    public void verify(go.micro.service.signup.SignupOuterClass.VerifyRequest request,
        io.grpc.stub.StreamObserver<go.micro.service.signup.SignupOuterClass.VerifyResponse> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(getVerifyMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void setPaymentMethod(go.micro.service.signup.SignupOuterClass.SetPaymentMethodRequest request,
        io.grpc.stub.StreamObserver<go.micro.service.signup.SignupOuterClass.SetPaymentMethodResponse> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(getSetPaymentMethodMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void hasPaymentMethod(go.micro.service.signup.SignupOuterClass.HasPaymentMethodRequest request,
        io.grpc.stub.StreamObserver<go.micro.service.signup.SignupOuterClass.HasPaymentMethodResponse> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(getHasPaymentMethodMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * Creates a subscription and an account
     * </pre>
     */
    public void completeSignup(go.micro.service.signup.SignupOuterClass.CompleteSignupRequest request,
        io.grpc.stub.StreamObserver<go.micro.service.signup.SignupOuterClass.CompleteSignupResponse> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(getCompleteSignupMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void recover(go.micro.service.signup.SignupOuterClass.RecoverRequest request,
        io.grpc.stub.StreamObserver<go.micro.service.signup.SignupOuterClass.RecoverResponse> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(getRecoverMethod(), getCallOptions()), request, responseObserver);
    }
  }

  /**
   */
  public static final class SignupBlockingStub extends io.grpc.stub.AbstractBlockingStub<SignupBlockingStub> {
    private SignupBlockingStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected SignupBlockingStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new SignupBlockingStub(channel, callOptions);
    }

    /**
     * <pre>
     * Sends the verification email to the user
     * </pre>
     */
    public go.micro.service.signup.SignupOuterClass.SendVerificationEmailResponse sendVerificationEmail(go.micro.service.signup.SignupOuterClass.SendVerificationEmailRequest request) {
      return blockingUnaryCall(
          getChannel(), getSendVerificationEmailMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * Verify kicks off the process of verification
     * </pre>
     */
    public go.micro.service.signup.SignupOuterClass.VerifyResponse verify(go.micro.service.signup.SignupOuterClass.VerifyRequest request) {
      return blockingUnaryCall(
          getChannel(), getVerifyMethod(), getCallOptions(), request);
    }

    /**
     */
    public go.micro.service.signup.SignupOuterClass.SetPaymentMethodResponse setPaymentMethod(go.micro.service.signup.SignupOuterClass.SetPaymentMethodRequest request) {
      return blockingUnaryCall(
          getChannel(), getSetPaymentMethodMethod(), getCallOptions(), request);
    }

    /**
     */
    public go.micro.service.signup.SignupOuterClass.HasPaymentMethodResponse hasPaymentMethod(go.micro.service.signup.SignupOuterClass.HasPaymentMethodRequest request) {
      return blockingUnaryCall(
          getChannel(), getHasPaymentMethodMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * Creates a subscription and an account
     * </pre>
     */
    public go.micro.service.signup.SignupOuterClass.CompleteSignupResponse completeSignup(go.micro.service.signup.SignupOuterClass.CompleteSignupRequest request) {
      return blockingUnaryCall(
          getChannel(), getCompleteSignupMethod(), getCallOptions(), request);
    }

    /**
     */
    public go.micro.service.signup.SignupOuterClass.RecoverResponse recover(go.micro.service.signup.SignupOuterClass.RecoverRequest request) {
      return blockingUnaryCall(
          getChannel(), getRecoverMethod(), getCallOptions(), request);
    }
  }

  /**
   */
  public static final class SignupFutureStub extends io.grpc.stub.AbstractFutureStub<SignupFutureStub> {
    private SignupFutureStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected SignupFutureStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new SignupFutureStub(channel, callOptions);
    }

    /**
     * <pre>
     * Sends the verification email to the user
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<go.micro.service.signup.SignupOuterClass.SendVerificationEmailResponse> sendVerificationEmail(
        go.micro.service.signup.SignupOuterClass.SendVerificationEmailRequest request) {
      return futureUnaryCall(
          getChannel().newCall(getSendVerificationEmailMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * Verify kicks off the process of verification
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<go.micro.service.signup.SignupOuterClass.VerifyResponse> verify(
        go.micro.service.signup.SignupOuterClass.VerifyRequest request) {
      return futureUnaryCall(
          getChannel().newCall(getVerifyMethod(), getCallOptions()), request);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<go.micro.service.signup.SignupOuterClass.SetPaymentMethodResponse> setPaymentMethod(
        go.micro.service.signup.SignupOuterClass.SetPaymentMethodRequest request) {
      return futureUnaryCall(
          getChannel().newCall(getSetPaymentMethodMethod(), getCallOptions()), request);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<go.micro.service.signup.SignupOuterClass.HasPaymentMethodResponse> hasPaymentMethod(
        go.micro.service.signup.SignupOuterClass.HasPaymentMethodRequest request) {
      return futureUnaryCall(
          getChannel().newCall(getHasPaymentMethodMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * Creates a subscription and an account
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<go.micro.service.signup.SignupOuterClass.CompleteSignupResponse> completeSignup(
        go.micro.service.signup.SignupOuterClass.CompleteSignupRequest request) {
      return futureUnaryCall(
          getChannel().newCall(getCompleteSignupMethod(), getCallOptions()), request);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<go.micro.service.signup.SignupOuterClass.RecoverResponse> recover(
        go.micro.service.signup.SignupOuterClass.RecoverRequest request) {
      return futureUnaryCall(
          getChannel().newCall(getRecoverMethod(), getCallOptions()), request);
    }
  }

  private static final int METHODID_SEND_VERIFICATION_EMAIL = 0;
  private static final int METHODID_VERIFY = 1;
  private static final int METHODID_SET_PAYMENT_METHOD = 2;
  private static final int METHODID_HAS_PAYMENT_METHOD = 3;
  private static final int METHODID_COMPLETE_SIGNUP = 4;
  private static final int METHODID_RECOVER = 5;

  private static final class MethodHandlers<Req, Resp> implements
      io.grpc.stub.ServerCalls.UnaryMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ServerStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ClientStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.BidiStreamingMethod<Req, Resp> {
    private final SignupImplBase serviceImpl;
    private final int methodId;

    MethodHandlers(SignupImplBase serviceImpl, int methodId) {
      this.serviceImpl = serviceImpl;
      this.methodId = methodId;
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public void invoke(Req request, io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        case METHODID_SEND_VERIFICATION_EMAIL:
          serviceImpl.sendVerificationEmail((go.micro.service.signup.SignupOuterClass.SendVerificationEmailRequest) request,
              (io.grpc.stub.StreamObserver<go.micro.service.signup.SignupOuterClass.SendVerificationEmailResponse>) responseObserver);
          break;
        case METHODID_VERIFY:
          serviceImpl.verify((go.micro.service.signup.SignupOuterClass.VerifyRequest) request,
              (io.grpc.stub.StreamObserver<go.micro.service.signup.SignupOuterClass.VerifyResponse>) responseObserver);
          break;
        case METHODID_SET_PAYMENT_METHOD:
          serviceImpl.setPaymentMethod((go.micro.service.signup.SignupOuterClass.SetPaymentMethodRequest) request,
              (io.grpc.stub.StreamObserver<go.micro.service.signup.SignupOuterClass.SetPaymentMethodResponse>) responseObserver);
          break;
        case METHODID_HAS_PAYMENT_METHOD:
          serviceImpl.hasPaymentMethod((go.micro.service.signup.SignupOuterClass.HasPaymentMethodRequest) request,
              (io.grpc.stub.StreamObserver<go.micro.service.signup.SignupOuterClass.HasPaymentMethodResponse>) responseObserver);
          break;
        case METHODID_COMPLETE_SIGNUP:
          serviceImpl.completeSignup((go.micro.service.signup.SignupOuterClass.CompleteSignupRequest) request,
              (io.grpc.stub.StreamObserver<go.micro.service.signup.SignupOuterClass.CompleteSignupResponse>) responseObserver);
          break;
        case METHODID_RECOVER:
          serviceImpl.recover((go.micro.service.signup.SignupOuterClass.RecoverRequest) request,
              (io.grpc.stub.StreamObserver<go.micro.service.signup.SignupOuterClass.RecoverResponse>) responseObserver);
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

  private static abstract class SignupBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier, io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    SignupBaseDescriptorSupplier() {}

    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return go.micro.service.signup.SignupOuterClass.getDescriptor();
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.ServiceDescriptor getServiceDescriptor() {
      return getFileDescriptor().findServiceByName("Signup");
    }
  }

  private static final class SignupFileDescriptorSupplier
      extends SignupBaseDescriptorSupplier {
    SignupFileDescriptorSupplier() {}
  }

  private static final class SignupMethodDescriptorSupplier
      extends SignupBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoMethodDescriptorSupplier {
    private final String methodName;

    SignupMethodDescriptorSupplier(String methodName) {
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
      synchronized (SignupGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor = result = io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
              .setSchemaDescriptor(new SignupFileDescriptorSupplier())
              .addMethod(getSendVerificationEmailMethod())
              .addMethod(getVerifyMethod())
              .addMethod(getSetPaymentMethodMethod())
              .addMethod(getHasPaymentMethodMethod())
              .addMethod(getCompleteSignupMethod())
              .addMethod(getRecoverMethod())
              .build();
        }
      }
    }
    return result;
  }
}
