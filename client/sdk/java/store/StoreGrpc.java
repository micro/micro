package store;

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
    comments = "Source: store/store.proto")
public final class StoreGrpc {

  private StoreGrpc() {}

  public static final String SERVICE_NAME = "store.Store";

  // Static method descriptors that strictly reflect the proto.
  private static volatile io.grpc.MethodDescriptor<store.StoreOuterClass.ReadRequest,
      store.StoreOuterClass.ReadResponse> getReadMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "Read",
      requestType = store.StoreOuterClass.ReadRequest.class,
      responseType = store.StoreOuterClass.ReadResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<store.StoreOuterClass.ReadRequest,
      store.StoreOuterClass.ReadResponse> getReadMethod() {
    io.grpc.MethodDescriptor<store.StoreOuterClass.ReadRequest, store.StoreOuterClass.ReadResponse> getReadMethod;
    if ((getReadMethod = StoreGrpc.getReadMethod) == null) {
      synchronized (StoreGrpc.class) {
        if ((getReadMethod = StoreGrpc.getReadMethod) == null) {
          StoreGrpc.getReadMethod = getReadMethod =
              io.grpc.MethodDescriptor.<store.StoreOuterClass.ReadRequest, store.StoreOuterClass.ReadResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "Read"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  store.StoreOuterClass.ReadRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  store.StoreOuterClass.ReadResponse.getDefaultInstance()))
              .setSchemaDescriptor(new StoreMethodDescriptorSupplier("Read"))
              .build();
        }
      }
    }
    return getReadMethod;
  }

  private static volatile io.grpc.MethodDescriptor<store.StoreOuterClass.WriteRequest,
      store.StoreOuterClass.WriteResponse> getWriteMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "Write",
      requestType = store.StoreOuterClass.WriteRequest.class,
      responseType = store.StoreOuterClass.WriteResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<store.StoreOuterClass.WriteRequest,
      store.StoreOuterClass.WriteResponse> getWriteMethod() {
    io.grpc.MethodDescriptor<store.StoreOuterClass.WriteRequest, store.StoreOuterClass.WriteResponse> getWriteMethod;
    if ((getWriteMethod = StoreGrpc.getWriteMethod) == null) {
      synchronized (StoreGrpc.class) {
        if ((getWriteMethod = StoreGrpc.getWriteMethod) == null) {
          StoreGrpc.getWriteMethod = getWriteMethod =
              io.grpc.MethodDescriptor.<store.StoreOuterClass.WriteRequest, store.StoreOuterClass.WriteResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "Write"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  store.StoreOuterClass.WriteRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  store.StoreOuterClass.WriteResponse.getDefaultInstance()))
              .setSchemaDescriptor(new StoreMethodDescriptorSupplier("Write"))
              .build();
        }
      }
    }
    return getWriteMethod;
  }

  private static volatile io.grpc.MethodDescriptor<store.StoreOuterClass.DeleteRequest,
      store.StoreOuterClass.DeleteResponse> getDeleteMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "Delete",
      requestType = store.StoreOuterClass.DeleteRequest.class,
      responseType = store.StoreOuterClass.DeleteResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<store.StoreOuterClass.DeleteRequest,
      store.StoreOuterClass.DeleteResponse> getDeleteMethod() {
    io.grpc.MethodDescriptor<store.StoreOuterClass.DeleteRequest, store.StoreOuterClass.DeleteResponse> getDeleteMethod;
    if ((getDeleteMethod = StoreGrpc.getDeleteMethod) == null) {
      synchronized (StoreGrpc.class) {
        if ((getDeleteMethod = StoreGrpc.getDeleteMethod) == null) {
          StoreGrpc.getDeleteMethod = getDeleteMethod =
              io.grpc.MethodDescriptor.<store.StoreOuterClass.DeleteRequest, store.StoreOuterClass.DeleteResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "Delete"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  store.StoreOuterClass.DeleteRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  store.StoreOuterClass.DeleteResponse.getDefaultInstance()))
              .setSchemaDescriptor(new StoreMethodDescriptorSupplier("Delete"))
              .build();
        }
      }
    }
    return getDeleteMethod;
  }

  private static volatile io.grpc.MethodDescriptor<store.StoreOuterClass.ListRequest,
      store.StoreOuterClass.ListResponse> getListMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "List",
      requestType = store.StoreOuterClass.ListRequest.class,
      responseType = store.StoreOuterClass.ListResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.SERVER_STREAMING)
  public static io.grpc.MethodDescriptor<store.StoreOuterClass.ListRequest,
      store.StoreOuterClass.ListResponse> getListMethod() {
    io.grpc.MethodDescriptor<store.StoreOuterClass.ListRequest, store.StoreOuterClass.ListResponse> getListMethod;
    if ((getListMethod = StoreGrpc.getListMethod) == null) {
      synchronized (StoreGrpc.class) {
        if ((getListMethod = StoreGrpc.getListMethod) == null) {
          StoreGrpc.getListMethod = getListMethod =
              io.grpc.MethodDescriptor.<store.StoreOuterClass.ListRequest, store.StoreOuterClass.ListResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.SERVER_STREAMING)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "List"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  store.StoreOuterClass.ListRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  store.StoreOuterClass.ListResponse.getDefaultInstance()))
              .setSchemaDescriptor(new StoreMethodDescriptorSupplier("List"))
              .build();
        }
      }
    }
    return getListMethod;
  }

  private static volatile io.grpc.MethodDescriptor<store.StoreOuterClass.DatabasesRequest,
      store.StoreOuterClass.DatabasesResponse> getDatabasesMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "Databases",
      requestType = store.StoreOuterClass.DatabasesRequest.class,
      responseType = store.StoreOuterClass.DatabasesResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<store.StoreOuterClass.DatabasesRequest,
      store.StoreOuterClass.DatabasesResponse> getDatabasesMethod() {
    io.grpc.MethodDescriptor<store.StoreOuterClass.DatabasesRequest, store.StoreOuterClass.DatabasesResponse> getDatabasesMethod;
    if ((getDatabasesMethod = StoreGrpc.getDatabasesMethod) == null) {
      synchronized (StoreGrpc.class) {
        if ((getDatabasesMethod = StoreGrpc.getDatabasesMethod) == null) {
          StoreGrpc.getDatabasesMethod = getDatabasesMethod =
              io.grpc.MethodDescriptor.<store.StoreOuterClass.DatabasesRequest, store.StoreOuterClass.DatabasesResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "Databases"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  store.StoreOuterClass.DatabasesRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  store.StoreOuterClass.DatabasesResponse.getDefaultInstance()))
              .setSchemaDescriptor(new StoreMethodDescriptorSupplier("Databases"))
              .build();
        }
      }
    }
    return getDatabasesMethod;
  }

  private static volatile io.grpc.MethodDescriptor<store.StoreOuterClass.TablesRequest,
      store.StoreOuterClass.TablesResponse> getTablesMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "Tables",
      requestType = store.StoreOuterClass.TablesRequest.class,
      responseType = store.StoreOuterClass.TablesResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<store.StoreOuterClass.TablesRequest,
      store.StoreOuterClass.TablesResponse> getTablesMethod() {
    io.grpc.MethodDescriptor<store.StoreOuterClass.TablesRequest, store.StoreOuterClass.TablesResponse> getTablesMethod;
    if ((getTablesMethod = StoreGrpc.getTablesMethod) == null) {
      synchronized (StoreGrpc.class) {
        if ((getTablesMethod = StoreGrpc.getTablesMethod) == null) {
          StoreGrpc.getTablesMethod = getTablesMethod =
              io.grpc.MethodDescriptor.<store.StoreOuterClass.TablesRequest, store.StoreOuterClass.TablesResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "Tables"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  store.StoreOuterClass.TablesRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  store.StoreOuterClass.TablesResponse.getDefaultInstance()))
              .setSchemaDescriptor(new StoreMethodDescriptorSupplier("Tables"))
              .build();
        }
      }
    }
    return getTablesMethod;
  }

  /**
   * Creates a new async stub that supports all call types for the service
   */
  public static StoreStub newStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<StoreStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<StoreStub>() {
        @java.lang.Override
        public StoreStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new StoreStub(channel, callOptions);
        }
      };
    return StoreStub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static StoreBlockingStub newBlockingStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<StoreBlockingStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<StoreBlockingStub>() {
        @java.lang.Override
        public StoreBlockingStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new StoreBlockingStub(channel, callOptions);
        }
      };
    return StoreBlockingStub.newStub(factory, channel);
  }

  /**
   * Creates a new ListenableFuture-style stub that supports unary calls on the service
   */
  public static StoreFutureStub newFutureStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<StoreFutureStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<StoreFutureStub>() {
        @java.lang.Override
        public StoreFutureStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new StoreFutureStub(channel, callOptions);
        }
      };
    return StoreFutureStub.newStub(factory, channel);
  }

  /**
   */
  public static abstract class StoreImplBase implements io.grpc.BindableService {

    /**
     */
    public void read(store.StoreOuterClass.ReadRequest request,
        io.grpc.stub.StreamObserver<store.StoreOuterClass.ReadResponse> responseObserver) {
      asyncUnimplementedUnaryCall(getReadMethod(), responseObserver);
    }

    /**
     */
    public void write(store.StoreOuterClass.WriteRequest request,
        io.grpc.stub.StreamObserver<store.StoreOuterClass.WriteResponse> responseObserver) {
      asyncUnimplementedUnaryCall(getWriteMethod(), responseObserver);
    }

    /**
     */
    public void delete(store.StoreOuterClass.DeleteRequest request,
        io.grpc.stub.StreamObserver<store.StoreOuterClass.DeleteResponse> responseObserver) {
      asyncUnimplementedUnaryCall(getDeleteMethod(), responseObserver);
    }

    /**
     */
    public void list(store.StoreOuterClass.ListRequest request,
        io.grpc.stub.StreamObserver<store.StoreOuterClass.ListResponse> responseObserver) {
      asyncUnimplementedUnaryCall(getListMethod(), responseObserver);
    }

    /**
     */
    public void databases(store.StoreOuterClass.DatabasesRequest request,
        io.grpc.stub.StreamObserver<store.StoreOuterClass.DatabasesResponse> responseObserver) {
      asyncUnimplementedUnaryCall(getDatabasesMethod(), responseObserver);
    }

    /**
     */
    public void tables(store.StoreOuterClass.TablesRequest request,
        io.grpc.stub.StreamObserver<store.StoreOuterClass.TablesResponse> responseObserver) {
      asyncUnimplementedUnaryCall(getTablesMethod(), responseObserver);
    }

    @java.lang.Override public final io.grpc.ServerServiceDefinition bindService() {
      return io.grpc.ServerServiceDefinition.builder(getServiceDescriptor())
          .addMethod(
            getReadMethod(),
            asyncUnaryCall(
              new MethodHandlers<
                store.StoreOuterClass.ReadRequest,
                store.StoreOuterClass.ReadResponse>(
                  this, METHODID_READ)))
          .addMethod(
            getWriteMethod(),
            asyncUnaryCall(
              new MethodHandlers<
                store.StoreOuterClass.WriteRequest,
                store.StoreOuterClass.WriteResponse>(
                  this, METHODID_WRITE)))
          .addMethod(
            getDeleteMethod(),
            asyncUnaryCall(
              new MethodHandlers<
                store.StoreOuterClass.DeleteRequest,
                store.StoreOuterClass.DeleteResponse>(
                  this, METHODID_DELETE)))
          .addMethod(
            getListMethod(),
            asyncServerStreamingCall(
              new MethodHandlers<
                store.StoreOuterClass.ListRequest,
                store.StoreOuterClass.ListResponse>(
                  this, METHODID_LIST)))
          .addMethod(
            getDatabasesMethod(),
            asyncUnaryCall(
              new MethodHandlers<
                store.StoreOuterClass.DatabasesRequest,
                store.StoreOuterClass.DatabasesResponse>(
                  this, METHODID_DATABASES)))
          .addMethod(
            getTablesMethod(),
            asyncUnaryCall(
              new MethodHandlers<
                store.StoreOuterClass.TablesRequest,
                store.StoreOuterClass.TablesResponse>(
                  this, METHODID_TABLES)))
          .build();
    }
  }

  /**
   */
  public static final class StoreStub extends io.grpc.stub.AbstractAsyncStub<StoreStub> {
    private StoreStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected StoreStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new StoreStub(channel, callOptions);
    }

    /**
     */
    public void read(store.StoreOuterClass.ReadRequest request,
        io.grpc.stub.StreamObserver<store.StoreOuterClass.ReadResponse> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(getReadMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void write(store.StoreOuterClass.WriteRequest request,
        io.grpc.stub.StreamObserver<store.StoreOuterClass.WriteResponse> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(getWriteMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void delete(store.StoreOuterClass.DeleteRequest request,
        io.grpc.stub.StreamObserver<store.StoreOuterClass.DeleteResponse> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(getDeleteMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void list(store.StoreOuterClass.ListRequest request,
        io.grpc.stub.StreamObserver<store.StoreOuterClass.ListResponse> responseObserver) {
      asyncServerStreamingCall(
          getChannel().newCall(getListMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void databases(store.StoreOuterClass.DatabasesRequest request,
        io.grpc.stub.StreamObserver<store.StoreOuterClass.DatabasesResponse> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(getDatabasesMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void tables(store.StoreOuterClass.TablesRequest request,
        io.grpc.stub.StreamObserver<store.StoreOuterClass.TablesResponse> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(getTablesMethod(), getCallOptions()), request, responseObserver);
    }
  }

  /**
   */
  public static final class StoreBlockingStub extends io.grpc.stub.AbstractBlockingStub<StoreBlockingStub> {
    private StoreBlockingStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected StoreBlockingStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new StoreBlockingStub(channel, callOptions);
    }

    /**
     */
    public store.StoreOuterClass.ReadResponse read(store.StoreOuterClass.ReadRequest request) {
      return blockingUnaryCall(
          getChannel(), getReadMethod(), getCallOptions(), request);
    }

    /**
     */
    public store.StoreOuterClass.WriteResponse write(store.StoreOuterClass.WriteRequest request) {
      return blockingUnaryCall(
          getChannel(), getWriteMethod(), getCallOptions(), request);
    }

    /**
     */
    public store.StoreOuterClass.DeleteResponse delete(store.StoreOuterClass.DeleteRequest request) {
      return blockingUnaryCall(
          getChannel(), getDeleteMethod(), getCallOptions(), request);
    }

    /**
     */
    public java.util.Iterator<store.StoreOuterClass.ListResponse> list(
        store.StoreOuterClass.ListRequest request) {
      return blockingServerStreamingCall(
          getChannel(), getListMethod(), getCallOptions(), request);
    }

    /**
     */
    public store.StoreOuterClass.DatabasesResponse databases(store.StoreOuterClass.DatabasesRequest request) {
      return blockingUnaryCall(
          getChannel(), getDatabasesMethod(), getCallOptions(), request);
    }

    /**
     */
    public store.StoreOuterClass.TablesResponse tables(store.StoreOuterClass.TablesRequest request) {
      return blockingUnaryCall(
          getChannel(), getTablesMethod(), getCallOptions(), request);
    }
  }

  /**
   */
  public static final class StoreFutureStub extends io.grpc.stub.AbstractFutureStub<StoreFutureStub> {
    private StoreFutureStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected StoreFutureStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new StoreFutureStub(channel, callOptions);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<store.StoreOuterClass.ReadResponse> read(
        store.StoreOuterClass.ReadRequest request) {
      return futureUnaryCall(
          getChannel().newCall(getReadMethod(), getCallOptions()), request);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<store.StoreOuterClass.WriteResponse> write(
        store.StoreOuterClass.WriteRequest request) {
      return futureUnaryCall(
          getChannel().newCall(getWriteMethod(), getCallOptions()), request);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<store.StoreOuterClass.DeleteResponse> delete(
        store.StoreOuterClass.DeleteRequest request) {
      return futureUnaryCall(
          getChannel().newCall(getDeleteMethod(), getCallOptions()), request);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<store.StoreOuterClass.DatabasesResponse> databases(
        store.StoreOuterClass.DatabasesRequest request) {
      return futureUnaryCall(
          getChannel().newCall(getDatabasesMethod(), getCallOptions()), request);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<store.StoreOuterClass.TablesResponse> tables(
        store.StoreOuterClass.TablesRequest request) {
      return futureUnaryCall(
          getChannel().newCall(getTablesMethod(), getCallOptions()), request);
    }
  }

  private static final int METHODID_READ = 0;
  private static final int METHODID_WRITE = 1;
  private static final int METHODID_DELETE = 2;
  private static final int METHODID_LIST = 3;
  private static final int METHODID_DATABASES = 4;
  private static final int METHODID_TABLES = 5;

  private static final class MethodHandlers<Req, Resp> implements
      io.grpc.stub.ServerCalls.UnaryMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ServerStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ClientStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.BidiStreamingMethod<Req, Resp> {
    private final StoreImplBase serviceImpl;
    private final int methodId;

    MethodHandlers(StoreImplBase serviceImpl, int methodId) {
      this.serviceImpl = serviceImpl;
      this.methodId = methodId;
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public void invoke(Req request, io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        case METHODID_READ:
          serviceImpl.read((store.StoreOuterClass.ReadRequest) request,
              (io.grpc.stub.StreamObserver<store.StoreOuterClass.ReadResponse>) responseObserver);
          break;
        case METHODID_WRITE:
          serviceImpl.write((store.StoreOuterClass.WriteRequest) request,
              (io.grpc.stub.StreamObserver<store.StoreOuterClass.WriteResponse>) responseObserver);
          break;
        case METHODID_DELETE:
          serviceImpl.delete((store.StoreOuterClass.DeleteRequest) request,
              (io.grpc.stub.StreamObserver<store.StoreOuterClass.DeleteResponse>) responseObserver);
          break;
        case METHODID_LIST:
          serviceImpl.list((store.StoreOuterClass.ListRequest) request,
              (io.grpc.stub.StreamObserver<store.StoreOuterClass.ListResponse>) responseObserver);
          break;
        case METHODID_DATABASES:
          serviceImpl.databases((store.StoreOuterClass.DatabasesRequest) request,
              (io.grpc.stub.StreamObserver<store.StoreOuterClass.DatabasesResponse>) responseObserver);
          break;
        case METHODID_TABLES:
          serviceImpl.tables((store.StoreOuterClass.TablesRequest) request,
              (io.grpc.stub.StreamObserver<store.StoreOuterClass.TablesResponse>) responseObserver);
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

  private static abstract class StoreBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier, io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    StoreBaseDescriptorSupplier() {}

    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return store.StoreOuterClass.getDescriptor();
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.ServiceDescriptor getServiceDescriptor() {
      return getFileDescriptor().findServiceByName("Store");
    }
  }

  private static final class StoreFileDescriptorSupplier
      extends StoreBaseDescriptorSupplier {
    StoreFileDescriptorSupplier() {}
  }

  private static final class StoreMethodDescriptorSupplier
      extends StoreBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoMethodDescriptorSupplier {
    private final String methodName;

    StoreMethodDescriptorSupplier(String methodName) {
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
      synchronized (StoreGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor = result = io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
              .setSchemaDescriptor(new StoreFileDescriptorSupplier())
              .addMethod(getReadMethod())
              .addMethod(getWriteMethod())
              .addMethod(getDeleteMethod())
              .addMethod(getListMethod())
              .addMethod(getDatabasesMethod())
              .addMethod(getTablesMethod())
              .build();
        }
      }
    }
    return result;
  }
}
