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
public final class RulesGrpc {

  private RulesGrpc() {}

  public static final String SERVICE_NAME = "auth.Rules";

  // Static method descriptors that strictly reflect the proto.
  private static volatile io.grpc.MethodDescriptor<auth.AuthOuterClass.CreateRequest,
      auth.AuthOuterClass.CreateResponse> getCreateMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "Create",
      requestType = auth.AuthOuterClass.CreateRequest.class,
      responseType = auth.AuthOuterClass.CreateResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<auth.AuthOuterClass.CreateRequest,
      auth.AuthOuterClass.CreateResponse> getCreateMethod() {
    io.grpc.MethodDescriptor<auth.AuthOuterClass.CreateRequest, auth.AuthOuterClass.CreateResponse> getCreateMethod;
    if ((getCreateMethod = RulesGrpc.getCreateMethod) == null) {
      synchronized (RulesGrpc.class) {
        if ((getCreateMethod = RulesGrpc.getCreateMethod) == null) {
          RulesGrpc.getCreateMethod = getCreateMethod =
              io.grpc.MethodDescriptor.<auth.AuthOuterClass.CreateRequest, auth.AuthOuterClass.CreateResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "Create"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  auth.AuthOuterClass.CreateRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  auth.AuthOuterClass.CreateResponse.getDefaultInstance()))
              .setSchemaDescriptor(new RulesMethodDescriptorSupplier("Create"))
              .build();
        }
      }
    }
    return getCreateMethod;
  }

  private static volatile io.grpc.MethodDescriptor<auth.AuthOuterClass.DeleteRequest,
      auth.AuthOuterClass.DeleteResponse> getDeleteMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "Delete",
      requestType = auth.AuthOuterClass.DeleteRequest.class,
      responseType = auth.AuthOuterClass.DeleteResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<auth.AuthOuterClass.DeleteRequest,
      auth.AuthOuterClass.DeleteResponse> getDeleteMethod() {
    io.grpc.MethodDescriptor<auth.AuthOuterClass.DeleteRequest, auth.AuthOuterClass.DeleteResponse> getDeleteMethod;
    if ((getDeleteMethod = RulesGrpc.getDeleteMethod) == null) {
      synchronized (RulesGrpc.class) {
        if ((getDeleteMethod = RulesGrpc.getDeleteMethod) == null) {
          RulesGrpc.getDeleteMethod = getDeleteMethod =
              io.grpc.MethodDescriptor.<auth.AuthOuterClass.DeleteRequest, auth.AuthOuterClass.DeleteResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "Delete"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  auth.AuthOuterClass.DeleteRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  auth.AuthOuterClass.DeleteResponse.getDefaultInstance()))
              .setSchemaDescriptor(new RulesMethodDescriptorSupplier("Delete"))
              .build();
        }
      }
    }
    return getDeleteMethod;
  }

  private static volatile io.grpc.MethodDescriptor<auth.AuthOuterClass.ListRequest,
      auth.AuthOuterClass.ListResponse> getListMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "List",
      requestType = auth.AuthOuterClass.ListRequest.class,
      responseType = auth.AuthOuterClass.ListResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<auth.AuthOuterClass.ListRequest,
      auth.AuthOuterClass.ListResponse> getListMethod() {
    io.grpc.MethodDescriptor<auth.AuthOuterClass.ListRequest, auth.AuthOuterClass.ListResponse> getListMethod;
    if ((getListMethod = RulesGrpc.getListMethod) == null) {
      synchronized (RulesGrpc.class) {
        if ((getListMethod = RulesGrpc.getListMethod) == null) {
          RulesGrpc.getListMethod = getListMethod =
              io.grpc.MethodDescriptor.<auth.AuthOuterClass.ListRequest, auth.AuthOuterClass.ListResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "List"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  auth.AuthOuterClass.ListRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  auth.AuthOuterClass.ListResponse.getDefaultInstance()))
              .setSchemaDescriptor(new RulesMethodDescriptorSupplier("List"))
              .build();
        }
      }
    }
    return getListMethod;
  }

  /**
   * Creates a new async stub that supports all call types for the service
   */
  public static RulesStub newStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<RulesStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<RulesStub>() {
        @java.lang.Override
        public RulesStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new RulesStub(channel, callOptions);
        }
      };
    return RulesStub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static RulesBlockingStub newBlockingStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<RulesBlockingStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<RulesBlockingStub>() {
        @java.lang.Override
        public RulesBlockingStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new RulesBlockingStub(channel, callOptions);
        }
      };
    return RulesBlockingStub.newStub(factory, channel);
  }

  /**
   * Creates a new ListenableFuture-style stub that supports unary calls on the service
   */
  public static RulesFutureStub newFutureStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<RulesFutureStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<RulesFutureStub>() {
        @java.lang.Override
        public RulesFutureStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new RulesFutureStub(channel, callOptions);
        }
      };
    return RulesFutureStub.newStub(factory, channel);
  }

  /**
   */
  public static abstract class RulesImplBase implements io.grpc.BindableService {

    /**
     */
    public void create(auth.AuthOuterClass.CreateRequest request,
        io.grpc.stub.StreamObserver<auth.AuthOuterClass.CreateResponse> responseObserver) {
      asyncUnimplementedUnaryCall(getCreateMethod(), responseObserver);
    }

    /**
     */
    public void delete(auth.AuthOuterClass.DeleteRequest request,
        io.grpc.stub.StreamObserver<auth.AuthOuterClass.DeleteResponse> responseObserver) {
      asyncUnimplementedUnaryCall(getDeleteMethod(), responseObserver);
    }

    /**
     */
    public void list(auth.AuthOuterClass.ListRequest request,
        io.grpc.stub.StreamObserver<auth.AuthOuterClass.ListResponse> responseObserver) {
      asyncUnimplementedUnaryCall(getListMethod(), responseObserver);
    }

    @java.lang.Override public final io.grpc.ServerServiceDefinition bindService() {
      return io.grpc.ServerServiceDefinition.builder(getServiceDescriptor())
          .addMethod(
            getCreateMethod(),
            asyncUnaryCall(
              new MethodHandlers<
                auth.AuthOuterClass.CreateRequest,
                auth.AuthOuterClass.CreateResponse>(
                  this, METHODID_CREATE)))
          .addMethod(
            getDeleteMethod(),
            asyncUnaryCall(
              new MethodHandlers<
                auth.AuthOuterClass.DeleteRequest,
                auth.AuthOuterClass.DeleteResponse>(
                  this, METHODID_DELETE)))
          .addMethod(
            getListMethod(),
            asyncUnaryCall(
              new MethodHandlers<
                auth.AuthOuterClass.ListRequest,
                auth.AuthOuterClass.ListResponse>(
                  this, METHODID_LIST)))
          .build();
    }
  }

  /**
   */
  public static final class RulesStub extends io.grpc.stub.AbstractAsyncStub<RulesStub> {
    private RulesStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected RulesStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new RulesStub(channel, callOptions);
    }

    /**
     */
    public void create(auth.AuthOuterClass.CreateRequest request,
        io.grpc.stub.StreamObserver<auth.AuthOuterClass.CreateResponse> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(getCreateMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void delete(auth.AuthOuterClass.DeleteRequest request,
        io.grpc.stub.StreamObserver<auth.AuthOuterClass.DeleteResponse> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(getDeleteMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void list(auth.AuthOuterClass.ListRequest request,
        io.grpc.stub.StreamObserver<auth.AuthOuterClass.ListResponse> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(getListMethod(), getCallOptions()), request, responseObserver);
    }
  }

  /**
   */
  public static final class RulesBlockingStub extends io.grpc.stub.AbstractBlockingStub<RulesBlockingStub> {
    private RulesBlockingStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected RulesBlockingStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new RulesBlockingStub(channel, callOptions);
    }

    /**
     */
    public auth.AuthOuterClass.CreateResponse create(auth.AuthOuterClass.CreateRequest request) {
      return blockingUnaryCall(
          getChannel(), getCreateMethod(), getCallOptions(), request);
    }

    /**
     */
    public auth.AuthOuterClass.DeleteResponse delete(auth.AuthOuterClass.DeleteRequest request) {
      return blockingUnaryCall(
          getChannel(), getDeleteMethod(), getCallOptions(), request);
    }

    /**
     */
    public auth.AuthOuterClass.ListResponse list(auth.AuthOuterClass.ListRequest request) {
      return blockingUnaryCall(
          getChannel(), getListMethod(), getCallOptions(), request);
    }
  }

  /**
   */
  public static final class RulesFutureStub extends io.grpc.stub.AbstractFutureStub<RulesFutureStub> {
    private RulesFutureStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected RulesFutureStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new RulesFutureStub(channel, callOptions);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<auth.AuthOuterClass.CreateResponse> create(
        auth.AuthOuterClass.CreateRequest request) {
      return futureUnaryCall(
          getChannel().newCall(getCreateMethod(), getCallOptions()), request);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<auth.AuthOuterClass.DeleteResponse> delete(
        auth.AuthOuterClass.DeleteRequest request) {
      return futureUnaryCall(
          getChannel().newCall(getDeleteMethod(), getCallOptions()), request);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<auth.AuthOuterClass.ListResponse> list(
        auth.AuthOuterClass.ListRequest request) {
      return futureUnaryCall(
          getChannel().newCall(getListMethod(), getCallOptions()), request);
    }
  }

  private static final int METHODID_CREATE = 0;
  private static final int METHODID_DELETE = 1;
  private static final int METHODID_LIST = 2;

  private static final class MethodHandlers<Req, Resp> implements
      io.grpc.stub.ServerCalls.UnaryMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ServerStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ClientStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.BidiStreamingMethod<Req, Resp> {
    private final RulesImplBase serviceImpl;
    private final int methodId;

    MethodHandlers(RulesImplBase serviceImpl, int methodId) {
      this.serviceImpl = serviceImpl;
      this.methodId = methodId;
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public void invoke(Req request, io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        case METHODID_CREATE:
          serviceImpl.create((auth.AuthOuterClass.CreateRequest) request,
              (io.grpc.stub.StreamObserver<auth.AuthOuterClass.CreateResponse>) responseObserver);
          break;
        case METHODID_DELETE:
          serviceImpl.delete((auth.AuthOuterClass.DeleteRequest) request,
              (io.grpc.stub.StreamObserver<auth.AuthOuterClass.DeleteResponse>) responseObserver);
          break;
        case METHODID_LIST:
          serviceImpl.list((auth.AuthOuterClass.ListRequest) request,
              (io.grpc.stub.StreamObserver<auth.AuthOuterClass.ListResponse>) responseObserver);
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

  private static abstract class RulesBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier, io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    RulesBaseDescriptorSupplier() {}

    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return auth.AuthOuterClass.getDescriptor();
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.ServiceDescriptor getServiceDescriptor() {
      return getFileDescriptor().findServiceByName("Rules");
    }
  }

  private static final class RulesFileDescriptorSupplier
      extends RulesBaseDescriptorSupplier {
    RulesFileDescriptorSupplier() {}
  }

  private static final class RulesMethodDescriptorSupplier
      extends RulesBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoMethodDescriptorSupplier {
    private final String methodName;

    RulesMethodDescriptorSupplier(String methodName) {
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
      synchronized (RulesGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor = result = io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
              .setSchemaDescriptor(new RulesFileDescriptorSupplier())
              .addMethod(getCreateMethod())
              .addMethod(getDeleteMethod())
              .addMethod(getListMethod())
              .build();
        }
      }
    }
    return result;
  }
}
