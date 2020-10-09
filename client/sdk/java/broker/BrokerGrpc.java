package broker;

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
    comments = "Source: broker/broker.proto")
public final class BrokerGrpc {

  private BrokerGrpc() {}

  public static final String SERVICE_NAME = "broker.Broker";

  // Static method descriptors that strictly reflect the proto.
  private static volatile io.grpc.MethodDescriptor<broker.BrokerOuterClass.PublishRequest,
      broker.BrokerOuterClass.Empty> getPublishMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "Publish",
      requestType = broker.BrokerOuterClass.PublishRequest.class,
      responseType = broker.BrokerOuterClass.Empty.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<broker.BrokerOuterClass.PublishRequest,
      broker.BrokerOuterClass.Empty> getPublishMethod() {
    io.grpc.MethodDescriptor<broker.BrokerOuterClass.PublishRequest, broker.BrokerOuterClass.Empty> getPublishMethod;
    if ((getPublishMethod = BrokerGrpc.getPublishMethod) == null) {
      synchronized (BrokerGrpc.class) {
        if ((getPublishMethod = BrokerGrpc.getPublishMethod) == null) {
          BrokerGrpc.getPublishMethod = getPublishMethod =
              io.grpc.MethodDescriptor.<broker.BrokerOuterClass.PublishRequest, broker.BrokerOuterClass.Empty>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "Publish"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  broker.BrokerOuterClass.PublishRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  broker.BrokerOuterClass.Empty.getDefaultInstance()))
              .setSchemaDescriptor(new BrokerMethodDescriptorSupplier("Publish"))
              .build();
        }
      }
    }
    return getPublishMethod;
  }

  private static volatile io.grpc.MethodDescriptor<broker.BrokerOuterClass.SubscribeRequest,
      broker.BrokerOuterClass.Message> getSubscribeMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "Subscribe",
      requestType = broker.BrokerOuterClass.SubscribeRequest.class,
      responseType = broker.BrokerOuterClass.Message.class,
      methodType = io.grpc.MethodDescriptor.MethodType.SERVER_STREAMING)
  public static io.grpc.MethodDescriptor<broker.BrokerOuterClass.SubscribeRequest,
      broker.BrokerOuterClass.Message> getSubscribeMethod() {
    io.grpc.MethodDescriptor<broker.BrokerOuterClass.SubscribeRequest, broker.BrokerOuterClass.Message> getSubscribeMethod;
    if ((getSubscribeMethod = BrokerGrpc.getSubscribeMethod) == null) {
      synchronized (BrokerGrpc.class) {
        if ((getSubscribeMethod = BrokerGrpc.getSubscribeMethod) == null) {
          BrokerGrpc.getSubscribeMethod = getSubscribeMethod =
              io.grpc.MethodDescriptor.<broker.BrokerOuterClass.SubscribeRequest, broker.BrokerOuterClass.Message>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.SERVER_STREAMING)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "Subscribe"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  broker.BrokerOuterClass.SubscribeRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  broker.BrokerOuterClass.Message.getDefaultInstance()))
              .setSchemaDescriptor(new BrokerMethodDescriptorSupplier("Subscribe"))
              .build();
        }
      }
    }
    return getSubscribeMethod;
  }

  /**
   * Creates a new async stub that supports all call types for the service
   */
  public static BrokerStub newStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<BrokerStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<BrokerStub>() {
        @java.lang.Override
        public BrokerStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new BrokerStub(channel, callOptions);
        }
      };
    return BrokerStub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static BrokerBlockingStub newBlockingStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<BrokerBlockingStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<BrokerBlockingStub>() {
        @java.lang.Override
        public BrokerBlockingStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new BrokerBlockingStub(channel, callOptions);
        }
      };
    return BrokerBlockingStub.newStub(factory, channel);
  }

  /**
   * Creates a new ListenableFuture-style stub that supports unary calls on the service
   */
  public static BrokerFutureStub newFutureStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<BrokerFutureStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<BrokerFutureStub>() {
        @java.lang.Override
        public BrokerFutureStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new BrokerFutureStub(channel, callOptions);
        }
      };
    return BrokerFutureStub.newStub(factory, channel);
  }

  /**
   */
  public static abstract class BrokerImplBase implements io.grpc.BindableService {

    /**
     */
    public void publish(broker.BrokerOuterClass.PublishRequest request,
        io.grpc.stub.StreamObserver<broker.BrokerOuterClass.Empty> responseObserver) {
      asyncUnimplementedUnaryCall(getPublishMethod(), responseObserver);
    }

    /**
     */
    public void subscribe(broker.BrokerOuterClass.SubscribeRequest request,
        io.grpc.stub.StreamObserver<broker.BrokerOuterClass.Message> responseObserver) {
      asyncUnimplementedUnaryCall(getSubscribeMethod(), responseObserver);
    }

    @java.lang.Override public final io.grpc.ServerServiceDefinition bindService() {
      return io.grpc.ServerServiceDefinition.builder(getServiceDescriptor())
          .addMethod(
            getPublishMethod(),
            asyncUnaryCall(
              new MethodHandlers<
                broker.BrokerOuterClass.PublishRequest,
                broker.BrokerOuterClass.Empty>(
                  this, METHODID_PUBLISH)))
          .addMethod(
            getSubscribeMethod(),
            asyncServerStreamingCall(
              new MethodHandlers<
                broker.BrokerOuterClass.SubscribeRequest,
                broker.BrokerOuterClass.Message>(
                  this, METHODID_SUBSCRIBE)))
          .build();
    }
  }

  /**
   */
  public static final class BrokerStub extends io.grpc.stub.AbstractAsyncStub<BrokerStub> {
    private BrokerStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected BrokerStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new BrokerStub(channel, callOptions);
    }

    /**
     */
    public void publish(broker.BrokerOuterClass.PublishRequest request,
        io.grpc.stub.StreamObserver<broker.BrokerOuterClass.Empty> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(getPublishMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void subscribe(broker.BrokerOuterClass.SubscribeRequest request,
        io.grpc.stub.StreamObserver<broker.BrokerOuterClass.Message> responseObserver) {
      asyncServerStreamingCall(
          getChannel().newCall(getSubscribeMethod(), getCallOptions()), request, responseObserver);
    }
  }

  /**
   */
  public static final class BrokerBlockingStub extends io.grpc.stub.AbstractBlockingStub<BrokerBlockingStub> {
    private BrokerBlockingStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected BrokerBlockingStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new BrokerBlockingStub(channel, callOptions);
    }

    /**
     */
    public broker.BrokerOuterClass.Empty publish(broker.BrokerOuterClass.PublishRequest request) {
      return blockingUnaryCall(
          getChannel(), getPublishMethod(), getCallOptions(), request);
    }

    /**
     */
    public java.util.Iterator<broker.BrokerOuterClass.Message> subscribe(
        broker.BrokerOuterClass.SubscribeRequest request) {
      return blockingServerStreamingCall(
          getChannel(), getSubscribeMethod(), getCallOptions(), request);
    }
  }

  /**
   */
  public static final class BrokerFutureStub extends io.grpc.stub.AbstractFutureStub<BrokerFutureStub> {
    private BrokerFutureStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected BrokerFutureStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new BrokerFutureStub(channel, callOptions);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<broker.BrokerOuterClass.Empty> publish(
        broker.BrokerOuterClass.PublishRequest request) {
      return futureUnaryCall(
          getChannel().newCall(getPublishMethod(), getCallOptions()), request);
    }
  }

  private static final int METHODID_PUBLISH = 0;
  private static final int METHODID_SUBSCRIBE = 1;

  private static final class MethodHandlers<Req, Resp> implements
      io.grpc.stub.ServerCalls.UnaryMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ServerStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ClientStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.BidiStreamingMethod<Req, Resp> {
    private final BrokerImplBase serviceImpl;
    private final int methodId;

    MethodHandlers(BrokerImplBase serviceImpl, int methodId) {
      this.serviceImpl = serviceImpl;
      this.methodId = methodId;
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public void invoke(Req request, io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        case METHODID_PUBLISH:
          serviceImpl.publish((broker.BrokerOuterClass.PublishRequest) request,
              (io.grpc.stub.StreamObserver<broker.BrokerOuterClass.Empty>) responseObserver);
          break;
        case METHODID_SUBSCRIBE:
          serviceImpl.subscribe((broker.BrokerOuterClass.SubscribeRequest) request,
              (io.grpc.stub.StreamObserver<broker.BrokerOuterClass.Message>) responseObserver);
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

  private static abstract class BrokerBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier, io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    BrokerBaseDescriptorSupplier() {}

    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return broker.BrokerOuterClass.getDescriptor();
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.ServiceDescriptor getServiceDescriptor() {
      return getFileDescriptor().findServiceByName("Broker");
    }
  }

  private static final class BrokerFileDescriptorSupplier
      extends BrokerBaseDescriptorSupplier {
    BrokerFileDescriptorSupplier() {}
  }

  private static final class BrokerMethodDescriptorSupplier
      extends BrokerBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoMethodDescriptorSupplier {
    private final String methodName;

    BrokerMethodDescriptorSupplier(String methodName) {
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
      synchronized (BrokerGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor = result = io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
              .setSchemaDescriptor(new BrokerFileDescriptorSupplier())
              .addMethod(getPublishMethod())
              .addMethod(getSubscribeMethod())
              .build();
        }
      }
    }
    return result;
  }
}
