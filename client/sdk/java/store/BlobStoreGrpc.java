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
public final class BlobStoreGrpc {

  private BlobStoreGrpc() {}

  public static final String SERVICE_NAME = "store.BlobStore";

  // Static method descriptors that strictly reflect the proto.
  private static volatile io.grpc.MethodDescriptor<store.StoreOuterClass.BlobReadRequest,
      store.StoreOuterClass.BlobReadResponse> getReadMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "Read",
      requestType = store.StoreOuterClass.BlobReadRequest.class,
      responseType = store.StoreOuterClass.BlobReadResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.SERVER_STREAMING)
  public static io.grpc.MethodDescriptor<store.StoreOuterClass.BlobReadRequest,
      store.StoreOuterClass.BlobReadResponse> getReadMethod() {
    io.grpc.MethodDescriptor<store.StoreOuterClass.BlobReadRequest, store.StoreOuterClass.BlobReadResponse> getReadMethod;
    if ((getReadMethod = BlobStoreGrpc.getReadMethod) == null) {
      synchronized (BlobStoreGrpc.class) {
        if ((getReadMethod = BlobStoreGrpc.getReadMethod) == null) {
          BlobStoreGrpc.getReadMethod = getReadMethod =
              io.grpc.MethodDescriptor.<store.StoreOuterClass.BlobReadRequest, store.StoreOuterClass.BlobReadResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.SERVER_STREAMING)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "Read"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  store.StoreOuterClass.BlobReadRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  store.StoreOuterClass.BlobReadResponse.getDefaultInstance()))
              .setSchemaDescriptor(new BlobStoreMethodDescriptorSupplier("Read"))
              .build();
        }
      }
    }
    return getReadMethod;
  }

  private static volatile io.grpc.MethodDescriptor<store.StoreOuterClass.BlobWriteRequest,
      store.StoreOuterClass.BlobWriteResponse> getWriteMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "Write",
      requestType = store.StoreOuterClass.BlobWriteRequest.class,
      responseType = store.StoreOuterClass.BlobWriteResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.CLIENT_STREAMING)
  public static io.grpc.MethodDescriptor<store.StoreOuterClass.BlobWriteRequest,
      store.StoreOuterClass.BlobWriteResponse> getWriteMethod() {
    io.grpc.MethodDescriptor<store.StoreOuterClass.BlobWriteRequest, store.StoreOuterClass.BlobWriteResponse> getWriteMethod;
    if ((getWriteMethod = BlobStoreGrpc.getWriteMethod) == null) {
      synchronized (BlobStoreGrpc.class) {
        if ((getWriteMethod = BlobStoreGrpc.getWriteMethod) == null) {
          BlobStoreGrpc.getWriteMethod = getWriteMethod =
              io.grpc.MethodDescriptor.<store.StoreOuterClass.BlobWriteRequest, store.StoreOuterClass.BlobWriteResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.CLIENT_STREAMING)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "Write"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  store.StoreOuterClass.BlobWriteRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  store.StoreOuterClass.BlobWriteResponse.getDefaultInstance()))
              .setSchemaDescriptor(new BlobStoreMethodDescriptorSupplier("Write"))
              .build();
        }
      }
    }
    return getWriteMethod;
  }

  private static volatile io.grpc.MethodDescriptor<store.StoreOuterClass.BlobDeleteRequest,
      store.StoreOuterClass.BlobDeleteResponse> getDeleteMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "Delete",
      requestType = store.StoreOuterClass.BlobDeleteRequest.class,
      responseType = store.StoreOuterClass.BlobDeleteResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<store.StoreOuterClass.BlobDeleteRequest,
      store.StoreOuterClass.BlobDeleteResponse> getDeleteMethod() {
    io.grpc.MethodDescriptor<store.StoreOuterClass.BlobDeleteRequest, store.StoreOuterClass.BlobDeleteResponse> getDeleteMethod;
    if ((getDeleteMethod = BlobStoreGrpc.getDeleteMethod) == null) {
      synchronized (BlobStoreGrpc.class) {
        if ((getDeleteMethod = BlobStoreGrpc.getDeleteMethod) == null) {
          BlobStoreGrpc.getDeleteMethod = getDeleteMethod =
              io.grpc.MethodDescriptor.<store.StoreOuterClass.BlobDeleteRequest, store.StoreOuterClass.BlobDeleteResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "Delete"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  store.StoreOuterClass.BlobDeleteRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  store.StoreOuterClass.BlobDeleteResponse.getDefaultInstance()))
              .setSchemaDescriptor(new BlobStoreMethodDescriptorSupplier("Delete"))
              .build();
        }
      }
    }
    return getDeleteMethod;
  }

  /**
   * Creates a new async stub that supports all call types for the service
   */
  public static BlobStoreStub newStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<BlobStoreStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<BlobStoreStub>() {
        @java.lang.Override
        public BlobStoreStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new BlobStoreStub(channel, callOptions);
        }
      };
    return BlobStoreStub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static BlobStoreBlockingStub newBlockingStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<BlobStoreBlockingStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<BlobStoreBlockingStub>() {
        @java.lang.Override
        public BlobStoreBlockingStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new BlobStoreBlockingStub(channel, callOptions);
        }
      };
    return BlobStoreBlockingStub.newStub(factory, channel);
  }

  /**
   * Creates a new ListenableFuture-style stub that supports unary calls on the service
   */
  public static BlobStoreFutureStub newFutureStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<BlobStoreFutureStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<BlobStoreFutureStub>() {
        @java.lang.Override
        public BlobStoreFutureStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new BlobStoreFutureStub(channel, callOptions);
        }
      };
    return BlobStoreFutureStub.newStub(factory, channel);
  }

  /**
   */
  public static abstract class BlobStoreImplBase implements io.grpc.BindableService {

    /**
     */
    public void read(store.StoreOuterClass.BlobReadRequest request,
        io.grpc.stub.StreamObserver<store.StoreOuterClass.BlobReadResponse> responseObserver) {
      asyncUnimplementedUnaryCall(getReadMethod(), responseObserver);
    }

    /**
     */
    public io.grpc.stub.StreamObserver<store.StoreOuterClass.BlobWriteRequest> write(
        io.grpc.stub.StreamObserver<store.StoreOuterClass.BlobWriteResponse> responseObserver) {
      return asyncUnimplementedStreamingCall(getWriteMethod(), responseObserver);
    }

    /**
     */
    public void delete(store.StoreOuterClass.BlobDeleteRequest request,
        io.grpc.stub.StreamObserver<store.StoreOuterClass.BlobDeleteResponse> responseObserver) {
      asyncUnimplementedUnaryCall(getDeleteMethod(), responseObserver);
    }

    @java.lang.Override public final io.grpc.ServerServiceDefinition bindService() {
      return io.grpc.ServerServiceDefinition.builder(getServiceDescriptor())
          .addMethod(
            getReadMethod(),
            asyncServerStreamingCall(
              new MethodHandlers<
                store.StoreOuterClass.BlobReadRequest,
                store.StoreOuterClass.BlobReadResponse>(
                  this, METHODID_READ)))
          .addMethod(
            getWriteMethod(),
            asyncClientStreamingCall(
              new MethodHandlers<
                store.StoreOuterClass.BlobWriteRequest,
                store.StoreOuterClass.BlobWriteResponse>(
                  this, METHODID_WRITE)))
          .addMethod(
            getDeleteMethod(),
            asyncUnaryCall(
              new MethodHandlers<
                store.StoreOuterClass.BlobDeleteRequest,
                store.StoreOuterClass.BlobDeleteResponse>(
                  this, METHODID_DELETE)))
          .build();
    }
  }

  /**
   */
  public static final class BlobStoreStub extends io.grpc.stub.AbstractAsyncStub<BlobStoreStub> {
    private BlobStoreStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected BlobStoreStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new BlobStoreStub(channel, callOptions);
    }

    /**
     */
    public void read(store.StoreOuterClass.BlobReadRequest request,
        io.grpc.stub.StreamObserver<store.StoreOuterClass.BlobReadResponse> responseObserver) {
      asyncServerStreamingCall(
          getChannel().newCall(getReadMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public io.grpc.stub.StreamObserver<store.StoreOuterClass.BlobWriteRequest> write(
        io.grpc.stub.StreamObserver<store.StoreOuterClass.BlobWriteResponse> responseObserver) {
      return asyncClientStreamingCall(
          getChannel().newCall(getWriteMethod(), getCallOptions()), responseObserver);
    }

    /**
     */
    public void delete(store.StoreOuterClass.BlobDeleteRequest request,
        io.grpc.stub.StreamObserver<store.StoreOuterClass.BlobDeleteResponse> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(getDeleteMethod(), getCallOptions()), request, responseObserver);
    }
  }

  /**
   */
  public static final class BlobStoreBlockingStub extends io.grpc.stub.AbstractBlockingStub<BlobStoreBlockingStub> {
    private BlobStoreBlockingStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected BlobStoreBlockingStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new BlobStoreBlockingStub(channel, callOptions);
    }

    /**
     */
    public java.util.Iterator<store.StoreOuterClass.BlobReadResponse> read(
        store.StoreOuterClass.BlobReadRequest request) {
      return blockingServerStreamingCall(
          getChannel(), getReadMethod(), getCallOptions(), request);
    }

    /**
     */
    public store.StoreOuterClass.BlobDeleteResponse delete(store.StoreOuterClass.BlobDeleteRequest request) {
      return blockingUnaryCall(
          getChannel(), getDeleteMethod(), getCallOptions(), request);
    }
  }

  /**
   */
  public static final class BlobStoreFutureStub extends io.grpc.stub.AbstractFutureStub<BlobStoreFutureStub> {
    private BlobStoreFutureStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected BlobStoreFutureStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new BlobStoreFutureStub(channel, callOptions);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<store.StoreOuterClass.BlobDeleteResponse> delete(
        store.StoreOuterClass.BlobDeleteRequest request) {
      return futureUnaryCall(
          getChannel().newCall(getDeleteMethod(), getCallOptions()), request);
    }
  }

  private static final int METHODID_READ = 0;
  private static final int METHODID_DELETE = 1;
  private static final int METHODID_WRITE = 2;

  private static final class MethodHandlers<Req, Resp> implements
      io.grpc.stub.ServerCalls.UnaryMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ServerStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ClientStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.BidiStreamingMethod<Req, Resp> {
    private final BlobStoreImplBase serviceImpl;
    private final int methodId;

    MethodHandlers(BlobStoreImplBase serviceImpl, int methodId) {
      this.serviceImpl = serviceImpl;
      this.methodId = methodId;
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public void invoke(Req request, io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        case METHODID_READ:
          serviceImpl.read((store.StoreOuterClass.BlobReadRequest) request,
              (io.grpc.stub.StreamObserver<store.StoreOuterClass.BlobReadResponse>) responseObserver);
          break;
        case METHODID_DELETE:
          serviceImpl.delete((store.StoreOuterClass.BlobDeleteRequest) request,
              (io.grpc.stub.StreamObserver<store.StoreOuterClass.BlobDeleteResponse>) responseObserver);
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
        case METHODID_WRITE:
          return (io.grpc.stub.StreamObserver<Req>) serviceImpl.write(
              (io.grpc.stub.StreamObserver<store.StoreOuterClass.BlobWriteResponse>) responseObserver);
        default:
          throw new AssertionError();
      }
    }
  }

  private static abstract class BlobStoreBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier, io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    BlobStoreBaseDescriptorSupplier() {}

    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return store.StoreOuterClass.getDescriptor();
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.ServiceDescriptor getServiceDescriptor() {
      return getFileDescriptor().findServiceByName("BlobStore");
    }
  }

  private static final class BlobStoreFileDescriptorSupplier
      extends BlobStoreBaseDescriptorSupplier {
    BlobStoreFileDescriptorSupplier() {}
  }

  private static final class BlobStoreMethodDescriptorSupplier
      extends BlobStoreBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoMethodDescriptorSupplier {
    private final String methodName;

    BlobStoreMethodDescriptorSupplier(String methodName) {
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
      synchronized (BlobStoreGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor = result = io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
              .setSchemaDescriptor(new BlobStoreFileDescriptorSupplier())
              .addMethod(getReadMethod())
              .addMethod(getWriteMethod())
              .addMethod(getDeleteMethod())
              .build();
        }
      }
    }
    return result;
  }
}
