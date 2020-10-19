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
public final class AccountsGrpc {

  private AccountsGrpc() {}

  public static final String SERVICE_NAME = "auth.Accounts";

  // Static method descriptors that strictly reflect the proto.
  private static volatile io.grpc.MethodDescriptor<auth.AuthOuterClass.ListAccountsRequest,
      auth.AuthOuterClass.ListAccountsResponse> getListMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "List",
      requestType = auth.AuthOuterClass.ListAccountsRequest.class,
      responseType = auth.AuthOuterClass.ListAccountsResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<auth.AuthOuterClass.ListAccountsRequest,
      auth.AuthOuterClass.ListAccountsResponse> getListMethod() {
    io.grpc.MethodDescriptor<auth.AuthOuterClass.ListAccountsRequest, auth.AuthOuterClass.ListAccountsResponse> getListMethod;
    if ((getListMethod = AccountsGrpc.getListMethod) == null) {
      synchronized (AccountsGrpc.class) {
        if ((getListMethod = AccountsGrpc.getListMethod) == null) {
          AccountsGrpc.getListMethod = getListMethod =
              io.grpc.MethodDescriptor.<auth.AuthOuterClass.ListAccountsRequest, auth.AuthOuterClass.ListAccountsResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "List"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  auth.AuthOuterClass.ListAccountsRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  auth.AuthOuterClass.ListAccountsResponse.getDefaultInstance()))
              .setSchemaDescriptor(new AccountsMethodDescriptorSupplier("List"))
              .build();
        }
      }
    }
    return getListMethod;
  }

  private static volatile io.grpc.MethodDescriptor<auth.AuthOuterClass.DeleteAccountRequest,
      auth.AuthOuterClass.DeleteAccountResponse> getDeleteMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "Delete",
      requestType = auth.AuthOuterClass.DeleteAccountRequest.class,
      responseType = auth.AuthOuterClass.DeleteAccountResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<auth.AuthOuterClass.DeleteAccountRequest,
      auth.AuthOuterClass.DeleteAccountResponse> getDeleteMethod() {
    io.grpc.MethodDescriptor<auth.AuthOuterClass.DeleteAccountRequest, auth.AuthOuterClass.DeleteAccountResponse> getDeleteMethod;
    if ((getDeleteMethod = AccountsGrpc.getDeleteMethod) == null) {
      synchronized (AccountsGrpc.class) {
        if ((getDeleteMethod = AccountsGrpc.getDeleteMethod) == null) {
          AccountsGrpc.getDeleteMethod = getDeleteMethod =
              io.grpc.MethodDescriptor.<auth.AuthOuterClass.DeleteAccountRequest, auth.AuthOuterClass.DeleteAccountResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "Delete"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  auth.AuthOuterClass.DeleteAccountRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  auth.AuthOuterClass.DeleteAccountResponse.getDefaultInstance()))
              .setSchemaDescriptor(new AccountsMethodDescriptorSupplier("Delete"))
              .build();
        }
      }
    }
    return getDeleteMethod;
  }

  private static volatile io.grpc.MethodDescriptor<auth.AuthOuterClass.ChangeSecretRequest,
      auth.AuthOuterClass.ChangeSecretResponse> getChangeSecretMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "ChangeSecret",
      requestType = auth.AuthOuterClass.ChangeSecretRequest.class,
      responseType = auth.AuthOuterClass.ChangeSecretResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<auth.AuthOuterClass.ChangeSecretRequest,
      auth.AuthOuterClass.ChangeSecretResponse> getChangeSecretMethod() {
    io.grpc.MethodDescriptor<auth.AuthOuterClass.ChangeSecretRequest, auth.AuthOuterClass.ChangeSecretResponse> getChangeSecretMethod;
    if ((getChangeSecretMethod = AccountsGrpc.getChangeSecretMethod) == null) {
      synchronized (AccountsGrpc.class) {
        if ((getChangeSecretMethod = AccountsGrpc.getChangeSecretMethod) == null) {
          AccountsGrpc.getChangeSecretMethod = getChangeSecretMethod =
              io.grpc.MethodDescriptor.<auth.AuthOuterClass.ChangeSecretRequest, auth.AuthOuterClass.ChangeSecretResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "ChangeSecret"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  auth.AuthOuterClass.ChangeSecretRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  auth.AuthOuterClass.ChangeSecretResponse.getDefaultInstance()))
              .setSchemaDescriptor(new AccountsMethodDescriptorSupplier("ChangeSecret"))
              .build();
        }
      }
    }
    return getChangeSecretMethod;
  }

  /**
   * Creates a new async stub that supports all call types for the service
   */
  public static AccountsStub newStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<AccountsStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<AccountsStub>() {
        @java.lang.Override
        public AccountsStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new AccountsStub(channel, callOptions);
        }
      };
    return AccountsStub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static AccountsBlockingStub newBlockingStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<AccountsBlockingStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<AccountsBlockingStub>() {
        @java.lang.Override
        public AccountsBlockingStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new AccountsBlockingStub(channel, callOptions);
        }
      };
    return AccountsBlockingStub.newStub(factory, channel);
  }

  /**
   * Creates a new ListenableFuture-style stub that supports unary calls on the service
   */
  public static AccountsFutureStub newFutureStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<AccountsFutureStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<AccountsFutureStub>() {
        @java.lang.Override
        public AccountsFutureStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new AccountsFutureStub(channel, callOptions);
        }
      };
    return AccountsFutureStub.newStub(factory, channel);
  }

  /**
   */
  public static abstract class AccountsImplBase implements io.grpc.BindableService {

    /**
     */
    public void list(auth.AuthOuterClass.ListAccountsRequest request,
        io.grpc.stub.StreamObserver<auth.AuthOuterClass.ListAccountsResponse> responseObserver) {
      asyncUnimplementedUnaryCall(getListMethod(), responseObserver);
    }

    /**
     */
    public void delete(auth.AuthOuterClass.DeleteAccountRequest request,
        io.grpc.stub.StreamObserver<auth.AuthOuterClass.DeleteAccountResponse> responseObserver) {
      asyncUnimplementedUnaryCall(getDeleteMethod(), responseObserver);
    }

    /**
     */
    public void changeSecret(auth.AuthOuterClass.ChangeSecretRequest request,
        io.grpc.stub.StreamObserver<auth.AuthOuterClass.ChangeSecretResponse> responseObserver) {
      asyncUnimplementedUnaryCall(getChangeSecretMethod(), responseObserver);
    }

    @java.lang.Override public final io.grpc.ServerServiceDefinition bindService() {
      return io.grpc.ServerServiceDefinition.builder(getServiceDescriptor())
          .addMethod(
            getListMethod(),
            asyncUnaryCall(
              new MethodHandlers<
                auth.AuthOuterClass.ListAccountsRequest,
                auth.AuthOuterClass.ListAccountsResponse>(
                  this, METHODID_LIST)))
          .addMethod(
            getDeleteMethod(),
            asyncUnaryCall(
              new MethodHandlers<
                auth.AuthOuterClass.DeleteAccountRequest,
                auth.AuthOuterClass.DeleteAccountResponse>(
                  this, METHODID_DELETE)))
          .addMethod(
            getChangeSecretMethod(),
            asyncUnaryCall(
              new MethodHandlers<
                auth.AuthOuterClass.ChangeSecretRequest,
                auth.AuthOuterClass.ChangeSecretResponse>(
                  this, METHODID_CHANGE_SECRET)))
          .build();
    }
  }

  /**
   */
  public static final class AccountsStub extends io.grpc.stub.AbstractAsyncStub<AccountsStub> {
    private AccountsStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected AccountsStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new AccountsStub(channel, callOptions);
    }

    /**
     */
    public void list(auth.AuthOuterClass.ListAccountsRequest request,
        io.grpc.stub.StreamObserver<auth.AuthOuterClass.ListAccountsResponse> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(getListMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void delete(auth.AuthOuterClass.DeleteAccountRequest request,
        io.grpc.stub.StreamObserver<auth.AuthOuterClass.DeleteAccountResponse> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(getDeleteMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void changeSecret(auth.AuthOuterClass.ChangeSecretRequest request,
        io.grpc.stub.StreamObserver<auth.AuthOuterClass.ChangeSecretResponse> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(getChangeSecretMethod(), getCallOptions()), request, responseObserver);
    }
  }

  /**
   */
  public static final class AccountsBlockingStub extends io.grpc.stub.AbstractBlockingStub<AccountsBlockingStub> {
    private AccountsBlockingStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected AccountsBlockingStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new AccountsBlockingStub(channel, callOptions);
    }

    /**
     */
    public auth.AuthOuterClass.ListAccountsResponse list(auth.AuthOuterClass.ListAccountsRequest request) {
      return blockingUnaryCall(
          getChannel(), getListMethod(), getCallOptions(), request);
    }

    /**
     */
    public auth.AuthOuterClass.DeleteAccountResponse delete(auth.AuthOuterClass.DeleteAccountRequest request) {
      return blockingUnaryCall(
          getChannel(), getDeleteMethod(), getCallOptions(), request);
    }

    /**
     */
    public auth.AuthOuterClass.ChangeSecretResponse changeSecret(auth.AuthOuterClass.ChangeSecretRequest request) {
      return blockingUnaryCall(
          getChannel(), getChangeSecretMethod(), getCallOptions(), request);
    }
  }

  /**
   */
  public static final class AccountsFutureStub extends io.grpc.stub.AbstractFutureStub<AccountsFutureStub> {
    private AccountsFutureStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected AccountsFutureStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new AccountsFutureStub(channel, callOptions);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<auth.AuthOuterClass.ListAccountsResponse> list(
        auth.AuthOuterClass.ListAccountsRequest request) {
      return futureUnaryCall(
          getChannel().newCall(getListMethod(), getCallOptions()), request);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<auth.AuthOuterClass.DeleteAccountResponse> delete(
        auth.AuthOuterClass.DeleteAccountRequest request) {
      return futureUnaryCall(
          getChannel().newCall(getDeleteMethod(), getCallOptions()), request);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<auth.AuthOuterClass.ChangeSecretResponse> changeSecret(
        auth.AuthOuterClass.ChangeSecretRequest request) {
      return futureUnaryCall(
          getChannel().newCall(getChangeSecretMethod(), getCallOptions()), request);
    }
  }

  private static final int METHODID_LIST = 0;
  private static final int METHODID_DELETE = 1;
  private static final int METHODID_CHANGE_SECRET = 2;

  private static final class MethodHandlers<Req, Resp> implements
      io.grpc.stub.ServerCalls.UnaryMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ServerStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ClientStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.BidiStreamingMethod<Req, Resp> {
    private final AccountsImplBase serviceImpl;
    private final int methodId;

    MethodHandlers(AccountsImplBase serviceImpl, int methodId) {
      this.serviceImpl = serviceImpl;
      this.methodId = methodId;
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public void invoke(Req request, io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        case METHODID_LIST:
          serviceImpl.list((auth.AuthOuterClass.ListAccountsRequest) request,
              (io.grpc.stub.StreamObserver<auth.AuthOuterClass.ListAccountsResponse>) responseObserver);
          break;
        case METHODID_DELETE:
          serviceImpl.delete((auth.AuthOuterClass.DeleteAccountRequest) request,
              (io.grpc.stub.StreamObserver<auth.AuthOuterClass.DeleteAccountResponse>) responseObserver);
          break;
        case METHODID_CHANGE_SECRET:
          serviceImpl.changeSecret((auth.AuthOuterClass.ChangeSecretRequest) request,
              (io.grpc.stub.StreamObserver<auth.AuthOuterClass.ChangeSecretResponse>) responseObserver);
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

  private static abstract class AccountsBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier, io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    AccountsBaseDescriptorSupplier() {}

    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return auth.AuthOuterClass.getDescriptor();
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.ServiceDescriptor getServiceDescriptor() {
      return getFileDescriptor().findServiceByName("Accounts");
    }
  }

  private static final class AccountsFileDescriptorSupplier
      extends AccountsBaseDescriptorSupplier {
    AccountsFileDescriptorSupplier() {}
  }

  private static final class AccountsMethodDescriptorSupplier
      extends AccountsBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoMethodDescriptorSupplier {
    private final String methodName;

    AccountsMethodDescriptorSupplier(String methodName) {
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
      synchronized (AccountsGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor = result = io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
              .setSchemaDescriptor(new AccountsFileDescriptorSupplier())
              .addMethod(getListMethod())
              .addMethod(getDeleteMethod())
              .addMethod(getChangeSecretMethod())
              .build();
        }
      }
    }
    return result;
  }
}
