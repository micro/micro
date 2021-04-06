package runtime;

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
    comments = "Source: runtime/runtime.proto")
public final class RuntimeGrpc {

  private RuntimeGrpc() {}

  public static final String SERVICE_NAME = "runtime.Runtime";

  // Static method descriptors that strictly reflect the proto.
  private static volatile io.grpc.MethodDescriptor<runtime.RuntimeOuterClass.CreateRequest,
      runtime.RuntimeOuterClass.CreateResponse> getCreateMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "Create",
      requestType = runtime.RuntimeOuterClass.CreateRequest.class,
      responseType = runtime.RuntimeOuterClass.CreateResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<runtime.RuntimeOuterClass.CreateRequest,
      runtime.RuntimeOuterClass.CreateResponse> getCreateMethod() {
    io.grpc.MethodDescriptor<runtime.RuntimeOuterClass.CreateRequest, runtime.RuntimeOuterClass.CreateResponse> getCreateMethod;
    if ((getCreateMethod = RuntimeGrpc.getCreateMethod) == null) {
      synchronized (RuntimeGrpc.class) {
        if ((getCreateMethod = RuntimeGrpc.getCreateMethod) == null) {
          RuntimeGrpc.getCreateMethod = getCreateMethod =
              io.grpc.MethodDescriptor.<runtime.RuntimeOuterClass.CreateRequest, runtime.RuntimeOuterClass.CreateResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "Create"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  runtime.RuntimeOuterClass.CreateRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  runtime.RuntimeOuterClass.CreateResponse.getDefaultInstance()))
              .setSchemaDescriptor(new RuntimeMethodDescriptorSupplier("Create"))
              .build();
        }
      }
    }
    return getCreateMethod;
  }

  private static volatile io.grpc.MethodDescriptor<runtime.RuntimeOuterClass.ReadRequest,
      runtime.RuntimeOuterClass.ReadResponse> getReadMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "Read",
      requestType = runtime.RuntimeOuterClass.ReadRequest.class,
      responseType = runtime.RuntimeOuterClass.ReadResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<runtime.RuntimeOuterClass.ReadRequest,
      runtime.RuntimeOuterClass.ReadResponse> getReadMethod() {
    io.grpc.MethodDescriptor<runtime.RuntimeOuterClass.ReadRequest, runtime.RuntimeOuterClass.ReadResponse> getReadMethod;
    if ((getReadMethod = RuntimeGrpc.getReadMethod) == null) {
      synchronized (RuntimeGrpc.class) {
        if ((getReadMethod = RuntimeGrpc.getReadMethod) == null) {
          RuntimeGrpc.getReadMethod = getReadMethod =
              io.grpc.MethodDescriptor.<runtime.RuntimeOuterClass.ReadRequest, runtime.RuntimeOuterClass.ReadResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "Read"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  runtime.RuntimeOuterClass.ReadRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  runtime.RuntimeOuterClass.ReadResponse.getDefaultInstance()))
              .setSchemaDescriptor(new RuntimeMethodDescriptorSupplier("Read"))
              .build();
        }
      }
    }
    return getReadMethod;
  }

  private static volatile io.grpc.MethodDescriptor<runtime.RuntimeOuterClass.DeleteRequest,
      runtime.RuntimeOuterClass.DeleteResponse> getDeleteMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "Delete",
      requestType = runtime.RuntimeOuterClass.DeleteRequest.class,
      responseType = runtime.RuntimeOuterClass.DeleteResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<runtime.RuntimeOuterClass.DeleteRequest,
      runtime.RuntimeOuterClass.DeleteResponse> getDeleteMethod() {
    io.grpc.MethodDescriptor<runtime.RuntimeOuterClass.DeleteRequest, runtime.RuntimeOuterClass.DeleteResponse> getDeleteMethod;
    if ((getDeleteMethod = RuntimeGrpc.getDeleteMethod) == null) {
      synchronized (RuntimeGrpc.class) {
        if ((getDeleteMethod = RuntimeGrpc.getDeleteMethod) == null) {
          RuntimeGrpc.getDeleteMethod = getDeleteMethod =
              io.grpc.MethodDescriptor.<runtime.RuntimeOuterClass.DeleteRequest, runtime.RuntimeOuterClass.DeleteResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "Delete"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  runtime.RuntimeOuterClass.DeleteRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  runtime.RuntimeOuterClass.DeleteResponse.getDefaultInstance()))
              .setSchemaDescriptor(new RuntimeMethodDescriptorSupplier("Delete"))
              .build();
        }
      }
    }
    return getDeleteMethod;
  }

  private static volatile io.grpc.MethodDescriptor<runtime.RuntimeOuterClass.UpdateRequest,
      runtime.RuntimeOuterClass.UpdateResponse> getUpdateMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "Update",
      requestType = runtime.RuntimeOuterClass.UpdateRequest.class,
      responseType = runtime.RuntimeOuterClass.UpdateResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<runtime.RuntimeOuterClass.UpdateRequest,
      runtime.RuntimeOuterClass.UpdateResponse> getUpdateMethod() {
    io.grpc.MethodDescriptor<runtime.RuntimeOuterClass.UpdateRequest, runtime.RuntimeOuterClass.UpdateResponse> getUpdateMethod;
    if ((getUpdateMethod = RuntimeGrpc.getUpdateMethod) == null) {
      synchronized (RuntimeGrpc.class) {
        if ((getUpdateMethod = RuntimeGrpc.getUpdateMethod) == null) {
          RuntimeGrpc.getUpdateMethod = getUpdateMethod =
              io.grpc.MethodDescriptor.<runtime.RuntimeOuterClass.UpdateRequest, runtime.RuntimeOuterClass.UpdateResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "Update"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  runtime.RuntimeOuterClass.UpdateRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  runtime.RuntimeOuterClass.UpdateResponse.getDefaultInstance()))
              .setSchemaDescriptor(new RuntimeMethodDescriptorSupplier("Update"))
              .build();
        }
      }
    }
    return getUpdateMethod;
  }

  private static volatile io.grpc.MethodDescriptor<runtime.RuntimeOuterClass.LogsRequest,
      runtime.RuntimeOuterClass.LogRecord> getLogsMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "Logs",
      requestType = runtime.RuntimeOuterClass.LogsRequest.class,
      responseType = runtime.RuntimeOuterClass.LogRecord.class,
      methodType = io.grpc.MethodDescriptor.MethodType.SERVER_STREAMING)
  public static io.grpc.MethodDescriptor<runtime.RuntimeOuterClass.LogsRequest,
      runtime.RuntimeOuterClass.LogRecord> getLogsMethod() {
    io.grpc.MethodDescriptor<runtime.RuntimeOuterClass.LogsRequest, runtime.RuntimeOuterClass.LogRecord> getLogsMethod;
    if ((getLogsMethod = RuntimeGrpc.getLogsMethod) == null) {
      synchronized (RuntimeGrpc.class) {
        if ((getLogsMethod = RuntimeGrpc.getLogsMethod) == null) {
          RuntimeGrpc.getLogsMethod = getLogsMethod =
              io.grpc.MethodDescriptor.<runtime.RuntimeOuterClass.LogsRequest, runtime.RuntimeOuterClass.LogRecord>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.SERVER_STREAMING)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "Logs"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  runtime.RuntimeOuterClass.LogsRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  runtime.RuntimeOuterClass.LogRecord.getDefaultInstance()))
              .setSchemaDescriptor(new RuntimeMethodDescriptorSupplier("Logs"))
              .build();
        }
      }
    }
    return getLogsMethod;
  }

  /**
   * Creates a new async stub that supports all call types for the service
   */
  public static RuntimeStub newStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<RuntimeStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<RuntimeStub>() {
        @java.lang.Override
        public RuntimeStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new RuntimeStub(channel, callOptions);
        }
      };
    return RuntimeStub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static RuntimeBlockingStub newBlockingStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<RuntimeBlockingStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<RuntimeBlockingStub>() {
        @java.lang.Override
        public RuntimeBlockingStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new RuntimeBlockingStub(channel, callOptions);
        }
      };
    return RuntimeBlockingStub.newStub(factory, channel);
  }

  /**
   * Creates a new ListenableFuture-style stub that supports unary calls on the service
   */
  public static RuntimeFutureStub newFutureStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<RuntimeFutureStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<RuntimeFutureStub>() {
        @java.lang.Override
        public RuntimeFutureStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new RuntimeFutureStub(channel, callOptions);
        }
      };
    return RuntimeFutureStub.newStub(factory, channel);
  }

  /**
   */
  public static abstract class RuntimeImplBase implements io.grpc.BindableService {

    /**
     */
    public void create(runtime.RuntimeOuterClass.CreateRequest request,
        io.grpc.stub.StreamObserver<runtime.RuntimeOuterClass.CreateResponse> responseObserver) {
      asyncUnimplementedUnaryCall(getCreateMethod(), responseObserver);
    }

    /**
     */
    public void read(runtime.RuntimeOuterClass.ReadRequest request,
        io.grpc.stub.StreamObserver<runtime.RuntimeOuterClass.ReadResponse> responseObserver) {
      asyncUnimplementedUnaryCall(getReadMethod(), responseObserver);
    }

    /**
     */
    public void delete(runtime.RuntimeOuterClass.DeleteRequest request,
        io.grpc.stub.StreamObserver<runtime.RuntimeOuterClass.DeleteResponse> responseObserver) {
      asyncUnimplementedUnaryCall(getDeleteMethod(), responseObserver);
    }

    /**
     */
    public void update(runtime.RuntimeOuterClass.UpdateRequest request,
        io.grpc.stub.StreamObserver<runtime.RuntimeOuterClass.UpdateResponse> responseObserver) {
      asyncUnimplementedUnaryCall(getUpdateMethod(), responseObserver);
    }

    /**
     */
    public void logs(runtime.RuntimeOuterClass.LogsRequest request,
        io.grpc.stub.StreamObserver<runtime.RuntimeOuterClass.LogRecord> responseObserver) {
      asyncUnimplementedUnaryCall(getLogsMethod(), responseObserver);
    }

    @java.lang.Override public final io.grpc.ServerServiceDefinition bindService() {
      return io.grpc.ServerServiceDefinition.builder(getServiceDescriptor())
          .addMethod(
            getCreateMethod(),
            asyncUnaryCall(
              new MethodHandlers<
                runtime.RuntimeOuterClass.CreateRequest,
                runtime.RuntimeOuterClass.CreateResponse>(
                  this, METHODID_CREATE)))
          .addMethod(
            getReadMethod(),
            asyncUnaryCall(
              new MethodHandlers<
                runtime.RuntimeOuterClass.ReadRequest,
                runtime.RuntimeOuterClass.ReadResponse>(
                  this, METHODID_READ)))
          .addMethod(
            getDeleteMethod(),
            asyncUnaryCall(
              new MethodHandlers<
                runtime.RuntimeOuterClass.DeleteRequest,
                runtime.RuntimeOuterClass.DeleteResponse>(
                  this, METHODID_DELETE)))
          .addMethod(
            getUpdateMethod(),
            asyncUnaryCall(
              new MethodHandlers<
                runtime.RuntimeOuterClass.UpdateRequest,
                runtime.RuntimeOuterClass.UpdateResponse>(
                  this, METHODID_UPDATE)))
          .addMethod(
            getLogsMethod(),
            asyncServerStreamingCall(
              new MethodHandlers<
                runtime.RuntimeOuterClass.LogsRequest,
                runtime.RuntimeOuterClass.LogRecord>(
                  this, METHODID_LOGS)))
          .build();
    }
  }

  /**
   */
  public static final class RuntimeStub extends io.grpc.stub.AbstractAsyncStub<RuntimeStub> {
    private RuntimeStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected RuntimeStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new RuntimeStub(channel, callOptions);
    }

    /**
     */
    public void create(runtime.RuntimeOuterClass.CreateRequest request,
        io.grpc.stub.StreamObserver<runtime.RuntimeOuterClass.CreateResponse> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(getCreateMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void read(runtime.RuntimeOuterClass.ReadRequest request,
        io.grpc.stub.StreamObserver<runtime.RuntimeOuterClass.ReadResponse> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(getReadMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void delete(runtime.RuntimeOuterClass.DeleteRequest request,
        io.grpc.stub.StreamObserver<runtime.RuntimeOuterClass.DeleteResponse> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(getDeleteMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void update(runtime.RuntimeOuterClass.UpdateRequest request,
        io.grpc.stub.StreamObserver<runtime.RuntimeOuterClass.UpdateResponse> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(getUpdateMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void logs(runtime.RuntimeOuterClass.LogsRequest request,
        io.grpc.stub.StreamObserver<runtime.RuntimeOuterClass.LogRecord> responseObserver) {
      asyncServerStreamingCall(
          getChannel().newCall(getLogsMethod(), getCallOptions()), request, responseObserver);
    }
  }

  /**
   */
  public static final class RuntimeBlockingStub extends io.grpc.stub.AbstractBlockingStub<RuntimeBlockingStub> {
    private RuntimeBlockingStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected RuntimeBlockingStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new RuntimeBlockingStub(channel, callOptions);
    }

    /**
     */
    public runtime.RuntimeOuterClass.CreateResponse create(runtime.RuntimeOuterClass.CreateRequest request) {
      return blockingUnaryCall(
          getChannel(), getCreateMethod(), getCallOptions(), request);
    }

    /**
     */
    public runtime.RuntimeOuterClass.ReadResponse read(runtime.RuntimeOuterClass.ReadRequest request) {
      return blockingUnaryCall(
          getChannel(), getReadMethod(), getCallOptions(), request);
    }

    /**
     */
    public runtime.RuntimeOuterClass.DeleteResponse delete(runtime.RuntimeOuterClass.DeleteRequest request) {
      return blockingUnaryCall(
          getChannel(), getDeleteMethod(), getCallOptions(), request);
    }

    /**
     */
    public runtime.RuntimeOuterClass.UpdateResponse update(runtime.RuntimeOuterClass.UpdateRequest request) {
      return blockingUnaryCall(
          getChannel(), getUpdateMethod(), getCallOptions(), request);
    }

    /**
     */
    public java.util.Iterator<runtime.RuntimeOuterClass.LogRecord> logs(
        runtime.RuntimeOuterClass.LogsRequest request) {
      return blockingServerStreamingCall(
          getChannel(), getLogsMethod(), getCallOptions(), request);
    }
  }

  /**
   */
  public static final class RuntimeFutureStub extends io.grpc.stub.AbstractFutureStub<RuntimeFutureStub> {
    private RuntimeFutureStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected RuntimeFutureStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new RuntimeFutureStub(channel, callOptions);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<runtime.RuntimeOuterClass.CreateResponse> create(
        runtime.RuntimeOuterClass.CreateRequest request) {
      return futureUnaryCall(
          getChannel().newCall(getCreateMethod(), getCallOptions()), request);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<runtime.RuntimeOuterClass.ReadResponse> read(
        runtime.RuntimeOuterClass.ReadRequest request) {
      return futureUnaryCall(
          getChannel().newCall(getReadMethod(), getCallOptions()), request);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<runtime.RuntimeOuterClass.DeleteResponse> delete(
        runtime.RuntimeOuterClass.DeleteRequest request) {
      return futureUnaryCall(
          getChannel().newCall(getDeleteMethod(), getCallOptions()), request);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<runtime.RuntimeOuterClass.UpdateResponse> update(
        runtime.RuntimeOuterClass.UpdateRequest request) {
      return futureUnaryCall(
          getChannel().newCall(getUpdateMethod(), getCallOptions()), request);
    }
  }

  private static final int METHODID_CREATE = 0;
  private static final int METHODID_READ = 1;
  private static final int METHODID_DELETE = 2;
  private static final int METHODID_UPDATE = 3;
  private static final int METHODID_LOGS = 4;

  private static final class MethodHandlers<Req, Resp> implements
      io.grpc.stub.ServerCalls.UnaryMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ServerStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ClientStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.BidiStreamingMethod<Req, Resp> {
    private final RuntimeImplBase serviceImpl;
    private final int methodId;

    MethodHandlers(RuntimeImplBase serviceImpl, int methodId) {
      this.serviceImpl = serviceImpl;
      this.methodId = methodId;
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public void invoke(Req request, io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        case METHODID_CREATE:
          serviceImpl.create((runtime.RuntimeOuterClass.CreateRequest) request,
              (io.grpc.stub.StreamObserver<runtime.RuntimeOuterClass.CreateResponse>) responseObserver);
          break;
        case METHODID_READ:
          serviceImpl.read((runtime.RuntimeOuterClass.ReadRequest) request,
              (io.grpc.stub.StreamObserver<runtime.RuntimeOuterClass.ReadResponse>) responseObserver);
          break;
        case METHODID_DELETE:
          serviceImpl.delete((runtime.RuntimeOuterClass.DeleteRequest) request,
              (io.grpc.stub.StreamObserver<runtime.RuntimeOuterClass.DeleteResponse>) responseObserver);
          break;
        case METHODID_UPDATE:
          serviceImpl.update((runtime.RuntimeOuterClass.UpdateRequest) request,
              (io.grpc.stub.StreamObserver<runtime.RuntimeOuterClass.UpdateResponse>) responseObserver);
          break;
        case METHODID_LOGS:
          serviceImpl.logs((runtime.RuntimeOuterClass.LogsRequest) request,
              (io.grpc.stub.StreamObserver<runtime.RuntimeOuterClass.LogRecord>) responseObserver);
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

  private static abstract class RuntimeBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier, io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    RuntimeBaseDescriptorSupplier() {}

    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return runtime.RuntimeOuterClass.getDescriptor();
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.ServiceDescriptor getServiceDescriptor() {
      return getFileDescriptor().findServiceByName("Runtime");
    }
  }

  private static final class RuntimeFileDescriptorSupplier
      extends RuntimeBaseDescriptorSupplier {
    RuntimeFileDescriptorSupplier() {}
  }

  private static final class RuntimeMethodDescriptorSupplier
      extends RuntimeBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoMethodDescriptorSupplier {
    private final String methodName;

    RuntimeMethodDescriptorSupplier(String methodName) {
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
      synchronized (RuntimeGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor = result = io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
              .setSchemaDescriptor(new RuntimeFileDescriptorSupplier())
              .addMethod(getCreateMethod())
              .addMethod(getReadMethod())
              .addMethod(getDeleteMethod())
              .addMethod(getUpdateMethod())
              .addMethod(getLogsMethod())
              .build();
        }
      }
    }
    return result;
  }
}
