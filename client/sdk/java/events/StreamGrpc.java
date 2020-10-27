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
public final class StreamGrpc {

  private StreamGrpc() {}

  public static final String SERVICE_NAME = "events.Stream";

  // Static method descriptors that strictly reflect the proto.
  private static volatile io.grpc.MethodDescriptor<events.Events.PublishRequest,
      events.Events.PublishResponse> getPublishMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "Publish",
      requestType = events.Events.PublishRequest.class,
      responseType = events.Events.PublishResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<events.Events.PublishRequest,
      events.Events.PublishResponse> getPublishMethod() {
    io.grpc.MethodDescriptor<events.Events.PublishRequest, events.Events.PublishResponse> getPublishMethod;
    if ((getPublishMethod = StreamGrpc.getPublishMethod) == null) {
      synchronized (StreamGrpc.class) {
        if ((getPublishMethod = StreamGrpc.getPublishMethod) == null) {
          StreamGrpc.getPublishMethod = getPublishMethod =
              io.grpc.MethodDescriptor.<events.Events.PublishRequest, events.Events.PublishResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "Publish"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  events.Events.PublishRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  events.Events.PublishResponse.getDefaultInstance()))
              .setSchemaDescriptor(new StreamMethodDescriptorSupplier("Publish"))
              .build();
        }
      }
    }
    return getPublishMethod;
  }

  private static volatile io.grpc.MethodDescriptor<events.Events.ConsumeRequest,
      events.Events.Event> getConsumeMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "Consume",
      requestType = events.Events.ConsumeRequest.class,
      responseType = events.Events.Event.class,
      methodType = io.grpc.MethodDescriptor.MethodType.SERVER_STREAMING)
  public static io.grpc.MethodDescriptor<events.Events.ConsumeRequest,
      events.Events.Event> getConsumeMethod() {
    io.grpc.MethodDescriptor<events.Events.ConsumeRequest, events.Events.Event> getConsumeMethod;
    if ((getConsumeMethod = StreamGrpc.getConsumeMethod) == null) {
      synchronized (StreamGrpc.class) {
        if ((getConsumeMethod = StreamGrpc.getConsumeMethod) == null) {
          StreamGrpc.getConsumeMethod = getConsumeMethod =
              io.grpc.MethodDescriptor.<events.Events.ConsumeRequest, events.Events.Event>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.SERVER_STREAMING)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "Consume"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  events.Events.ConsumeRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  events.Events.Event.getDefaultInstance()))
              .setSchemaDescriptor(new StreamMethodDescriptorSupplier("Consume"))
              .build();
        }
      }
    }
    return getConsumeMethod;
  }

  /**
   * Creates a new async stub that supports all call types for the service
   */
  public static StreamStub newStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<StreamStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<StreamStub>() {
        @java.lang.Override
        public StreamStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new StreamStub(channel, callOptions);
        }
      };
    return StreamStub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static StreamBlockingStub newBlockingStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<StreamBlockingStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<StreamBlockingStub>() {
        @java.lang.Override
        public StreamBlockingStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new StreamBlockingStub(channel, callOptions);
        }
      };
    return StreamBlockingStub.newStub(factory, channel);
  }

  /**
   * Creates a new ListenableFuture-style stub that supports unary calls on the service
   */
  public static StreamFutureStub newFutureStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<StreamFutureStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<StreamFutureStub>() {
        @java.lang.Override
        public StreamFutureStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new StreamFutureStub(channel, callOptions);
        }
      };
    return StreamFutureStub.newStub(factory, channel);
  }

  /**
   */
  public static abstract class StreamImplBase implements io.grpc.BindableService {

    /**
     */
    public void publish(events.Events.PublishRequest request,
        io.grpc.stub.StreamObserver<events.Events.PublishResponse> responseObserver) {
      asyncUnimplementedUnaryCall(getPublishMethod(), responseObserver);
    }

    /**
     */
    public void consume(events.Events.ConsumeRequest request,
        io.grpc.stub.StreamObserver<events.Events.Event> responseObserver) {
      asyncUnimplementedUnaryCall(getConsumeMethod(), responseObserver);
    }

    @java.lang.Override public final io.grpc.ServerServiceDefinition bindService() {
      return io.grpc.ServerServiceDefinition.builder(getServiceDescriptor())
          .addMethod(
            getPublishMethod(),
            asyncUnaryCall(
              new MethodHandlers<
                events.Events.PublishRequest,
                events.Events.PublishResponse>(
                  this, METHODID_PUBLISH)))
          .addMethod(
            getConsumeMethod(),
            asyncServerStreamingCall(
              new MethodHandlers<
                events.Events.ConsumeRequest,
                events.Events.Event>(
                  this, METHODID_CONSUME)))
          .build();
    }
  }

  /**
   */
  public static final class StreamStub extends io.grpc.stub.AbstractAsyncStub<StreamStub> {
    private StreamStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected StreamStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new StreamStub(channel, callOptions);
    }

    /**
     */
    public void publish(events.Events.PublishRequest request,
        io.grpc.stub.StreamObserver<events.Events.PublishResponse> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(getPublishMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void consume(events.Events.ConsumeRequest request,
        io.grpc.stub.StreamObserver<events.Events.Event> responseObserver) {
      asyncServerStreamingCall(
          getChannel().newCall(getConsumeMethod(), getCallOptions()), request, responseObserver);
    }
  }

  /**
   */
  public static final class StreamBlockingStub extends io.grpc.stub.AbstractBlockingStub<StreamBlockingStub> {
    private StreamBlockingStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected StreamBlockingStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new StreamBlockingStub(channel, callOptions);
    }

    /**
     */
    public events.Events.PublishResponse publish(events.Events.PublishRequest request) {
      return blockingUnaryCall(
          getChannel(), getPublishMethod(), getCallOptions(), request);
    }

    /**
     */
    public java.util.Iterator<events.Events.Event> consume(
        events.Events.ConsumeRequest request) {
      return blockingServerStreamingCall(
          getChannel(), getConsumeMethod(), getCallOptions(), request);
    }
  }

  /**
   */
  public static final class StreamFutureStub extends io.grpc.stub.AbstractFutureStub<StreamFutureStub> {
    private StreamFutureStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected StreamFutureStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new StreamFutureStub(channel, callOptions);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<events.Events.PublishResponse> publish(
        events.Events.PublishRequest request) {
      return futureUnaryCall(
          getChannel().newCall(getPublishMethod(), getCallOptions()), request);
    }
  }

  private static final int METHODID_PUBLISH = 0;
  private static final int METHODID_CONSUME = 1;

  private static final class MethodHandlers<Req, Resp> implements
      io.grpc.stub.ServerCalls.UnaryMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ServerStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ClientStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.BidiStreamingMethod<Req, Resp> {
    private final StreamImplBase serviceImpl;
    private final int methodId;

    MethodHandlers(StreamImplBase serviceImpl, int methodId) {
      this.serviceImpl = serviceImpl;
      this.methodId = methodId;
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public void invoke(Req request, io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        case METHODID_PUBLISH:
          serviceImpl.publish((events.Events.PublishRequest) request,
              (io.grpc.stub.StreamObserver<events.Events.PublishResponse>) responseObserver);
          break;
        case METHODID_CONSUME:
          serviceImpl.consume((events.Events.ConsumeRequest) request,
              (io.grpc.stub.StreamObserver<events.Events.Event>) responseObserver);
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

  private static abstract class StreamBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier, io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    StreamBaseDescriptorSupplier() {}

    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return events.Events.getDescriptor();
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.ServiceDescriptor getServiceDescriptor() {
      return getFileDescriptor().findServiceByName("Stream");
    }
  }

  private static final class StreamFileDescriptorSupplier
      extends StreamBaseDescriptorSupplier {
    StreamFileDescriptorSupplier() {}
  }

  private static final class StreamMethodDescriptorSupplier
      extends StreamBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoMethodDescriptorSupplier {
    private final String methodName;

    StreamMethodDescriptorSupplier(String methodName) {
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
      synchronized (StreamGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor = result = io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
              .setSchemaDescriptor(new StreamFileDescriptorSupplier())
              .addMethod(getPublishMethod())
              .addMethod(getConsumeMethod())
              .build();
        }
      }
    }
    return result;
  }
}
