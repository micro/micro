package events;

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
    comments = "Source: events/events.proto")
public final class StoreGrpc {

  private StoreGrpc() {}

  public static final String SERVICE_NAME = "events.Store";

  // Static method descriptors that strictly reflect the proto.
  private static volatile io.grpc.MethodDescriptor<events.Events.ReadRequest,
      events.Events.ReadResponse> getReadMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "Read",
      requestType = events.Events.ReadRequest.class,
      responseType = events.Events.ReadResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<events.Events.ReadRequest,
      events.Events.ReadResponse> getReadMethod() {
    io.grpc.MethodDescriptor<events.Events.ReadRequest, events.Events.ReadResponse> getReadMethod;
    if ((getReadMethod = StoreGrpc.getReadMethod) == null) {
      synchronized (StoreGrpc.class) {
        if ((getReadMethod = StoreGrpc.getReadMethod) == null) {
          StoreGrpc.getReadMethod = getReadMethod =
              io.grpc.MethodDescriptor.<events.Events.ReadRequest, events.Events.ReadResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "Read"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  events.Events.ReadRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  events.Events.ReadResponse.getDefaultInstance()))
              .setSchemaDescriptor(new StoreMethodDescriptorSupplier("Read"))
              .build();
        }
      }
    }
    return getReadMethod;
  }

  private static volatile io.grpc.MethodDescriptor<events.Events.WriteRequest,
      events.Events.WriteResponse> getWriteMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "Write",
      requestType = events.Events.WriteRequest.class,
      responseType = events.Events.WriteResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<events.Events.WriteRequest,
      events.Events.WriteResponse> getWriteMethod() {
    io.grpc.MethodDescriptor<events.Events.WriteRequest, events.Events.WriteResponse> getWriteMethod;
    if ((getWriteMethod = StoreGrpc.getWriteMethod) == null) {
      synchronized (StoreGrpc.class) {
        if ((getWriteMethod = StoreGrpc.getWriteMethod) == null) {
          StoreGrpc.getWriteMethod = getWriteMethod =
              io.grpc.MethodDescriptor.<events.Events.WriteRequest, events.Events.WriteResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "Write"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  events.Events.WriteRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  events.Events.WriteResponse.getDefaultInstance()))
              .setSchemaDescriptor(new StoreMethodDescriptorSupplier("Write"))
              .build();
        }
      }
    }
    return getWriteMethod;
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
    public void read(events.Events.ReadRequest request,
        io.grpc.stub.StreamObserver<events.Events.ReadResponse> responseObserver) {
      asyncUnimplementedUnaryCall(getReadMethod(), responseObserver);
    }

    /**
     */
    public void write(events.Events.WriteRequest request,
        io.grpc.stub.StreamObserver<events.Events.WriteResponse> responseObserver) {
      asyncUnimplementedUnaryCall(getWriteMethod(), responseObserver);
    }

    @java.lang.Override public final io.grpc.ServerServiceDefinition bindService() {
      return io.grpc.ServerServiceDefinition.builder(getServiceDescriptor())
          .addMethod(
            getReadMethod(),
            asyncUnaryCall(
              new MethodHandlers<
                events.Events.ReadRequest,
                events.Events.ReadResponse>(
                  this, METHODID_READ)))
          .addMethod(
            getWriteMethod(),
            asyncUnaryCall(
              new MethodHandlers<
                events.Events.WriteRequest,
                events.Events.WriteResponse>(
                  this, METHODID_WRITE)))
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
    public void read(events.Events.ReadRequest request,
        io.grpc.stub.StreamObserver<events.Events.ReadResponse> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(getReadMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void write(events.Events.WriteRequest request,
        io.grpc.stub.StreamObserver<events.Events.WriteResponse> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(getWriteMethod(), getCallOptions()), request, responseObserver);
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
    public events.Events.ReadResponse read(events.Events.ReadRequest request) {
      return blockingUnaryCall(
          getChannel(), getReadMethod(), getCallOptions(), request);
    }

    /**
     */
    public events.Events.WriteResponse write(events.Events.WriteRequest request) {
      return blockingUnaryCall(
          getChannel(), getWriteMethod(), getCallOptions(), request);
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
    public com.google.common.util.concurrent.ListenableFuture<events.Events.ReadResponse> read(
        events.Events.ReadRequest request) {
      return futureUnaryCall(
          getChannel().newCall(getReadMethod(), getCallOptions()), request);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<events.Events.WriteResponse> write(
        events.Events.WriteRequest request) {
      return futureUnaryCall(
          getChannel().newCall(getWriteMethod(), getCallOptions()), request);
    }
  }

  private static final int METHODID_READ = 0;
  private static final int METHODID_WRITE = 1;

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
          serviceImpl.read((events.Events.ReadRequest) request,
              (io.grpc.stub.StreamObserver<events.Events.ReadResponse>) responseObserver);
          break;
        case METHODID_WRITE:
          serviceImpl.write((events.Events.WriteRequest) request,
              (io.grpc.stub.StreamObserver<events.Events.WriteResponse>) responseObserver);
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
      return events.Events.getDescriptor();
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
              .build();
        }
      }
    }
    return result;
  }
}
