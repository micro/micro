package transport;

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
    comments = "Source: transport/transport.proto")
public final class TransportGrpc {

  private TransportGrpc() {}

  public static final String SERVICE_NAME = "transport.Transport";

  // Static method descriptors that strictly reflect the proto.
  private static volatile io.grpc.MethodDescriptor<transport.TransportOuterClass.Message,
      transport.TransportOuterClass.Message> getStreamMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "Stream",
      requestType = transport.TransportOuterClass.Message.class,
      responseType = transport.TransportOuterClass.Message.class,
      methodType = io.grpc.MethodDescriptor.MethodType.BIDI_STREAMING)
  public static io.grpc.MethodDescriptor<transport.TransportOuterClass.Message,
      transport.TransportOuterClass.Message> getStreamMethod() {
    io.grpc.MethodDescriptor<transport.TransportOuterClass.Message, transport.TransportOuterClass.Message> getStreamMethod;
    if ((getStreamMethod = TransportGrpc.getStreamMethod) == null) {
      synchronized (TransportGrpc.class) {
        if ((getStreamMethod = TransportGrpc.getStreamMethod) == null) {
          TransportGrpc.getStreamMethod = getStreamMethod =
              io.grpc.MethodDescriptor.<transport.TransportOuterClass.Message, transport.TransportOuterClass.Message>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.BIDI_STREAMING)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "Stream"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  transport.TransportOuterClass.Message.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  transport.TransportOuterClass.Message.getDefaultInstance()))
              .setSchemaDescriptor(new TransportMethodDescriptorSupplier("Stream"))
              .build();
        }
      }
    }
    return getStreamMethod;
  }

  /**
   * Creates a new async stub that supports all call types for the service
   */
  public static TransportStub newStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<TransportStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<TransportStub>() {
        @java.lang.Override
        public TransportStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new TransportStub(channel, callOptions);
        }
      };
    return TransportStub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static TransportBlockingStub newBlockingStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<TransportBlockingStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<TransportBlockingStub>() {
        @java.lang.Override
        public TransportBlockingStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new TransportBlockingStub(channel, callOptions);
        }
      };
    return TransportBlockingStub.newStub(factory, channel);
  }

  /**
   * Creates a new ListenableFuture-style stub that supports unary calls on the service
   */
  public static TransportFutureStub newFutureStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<TransportFutureStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<TransportFutureStub>() {
        @java.lang.Override
        public TransportFutureStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new TransportFutureStub(channel, callOptions);
        }
      };
    return TransportFutureStub.newStub(factory, channel);
  }

  /**
   */
  public static abstract class TransportImplBase implements io.grpc.BindableService {

    /**
     */
    public io.grpc.stub.StreamObserver<transport.TransportOuterClass.Message> stream(
        io.grpc.stub.StreamObserver<transport.TransportOuterClass.Message> responseObserver) {
      return asyncUnimplementedStreamingCall(getStreamMethod(), responseObserver);
    }

    @java.lang.Override public final io.grpc.ServerServiceDefinition bindService() {
      return io.grpc.ServerServiceDefinition.builder(getServiceDescriptor())
          .addMethod(
            getStreamMethod(),
            asyncBidiStreamingCall(
              new MethodHandlers<
                transport.TransportOuterClass.Message,
                transport.TransportOuterClass.Message>(
                  this, METHODID_STREAM)))
          .build();
    }
  }

  /**
   */
  public static final class TransportStub extends io.grpc.stub.AbstractAsyncStub<TransportStub> {
    private TransportStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected TransportStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new TransportStub(channel, callOptions);
    }

    /**
     */
    public io.grpc.stub.StreamObserver<transport.TransportOuterClass.Message> stream(
        io.grpc.stub.StreamObserver<transport.TransportOuterClass.Message> responseObserver) {
      return asyncBidiStreamingCall(
          getChannel().newCall(getStreamMethod(), getCallOptions()), responseObserver);
    }
  }

  /**
   */
  public static final class TransportBlockingStub extends io.grpc.stub.AbstractBlockingStub<TransportBlockingStub> {
    private TransportBlockingStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected TransportBlockingStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new TransportBlockingStub(channel, callOptions);
    }
  }

  /**
   */
  public static final class TransportFutureStub extends io.grpc.stub.AbstractFutureStub<TransportFutureStub> {
    private TransportFutureStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected TransportFutureStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new TransportFutureStub(channel, callOptions);
    }
  }

  private static final int METHODID_STREAM = 0;

  private static final class MethodHandlers<Req, Resp> implements
      io.grpc.stub.ServerCalls.UnaryMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ServerStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ClientStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.BidiStreamingMethod<Req, Resp> {
    private final TransportImplBase serviceImpl;
    private final int methodId;

    MethodHandlers(TransportImplBase serviceImpl, int methodId) {
      this.serviceImpl = serviceImpl;
      this.methodId = methodId;
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public void invoke(Req request, io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        default:
          throw new AssertionError();
      }
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public io.grpc.stub.StreamObserver<Req> invoke(
        io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        case METHODID_STREAM:
          return (io.grpc.stub.StreamObserver<Req>) serviceImpl.stream(
              (io.grpc.stub.StreamObserver<transport.TransportOuterClass.Message>) responseObserver);
        default:
          throw new AssertionError();
      }
    }
  }

  private static abstract class TransportBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier, io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    TransportBaseDescriptorSupplier() {}

    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return transport.TransportOuterClass.getDescriptor();
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.ServiceDescriptor getServiceDescriptor() {
      return getFileDescriptor().findServiceByName("Transport");
    }
  }

  private static final class TransportFileDescriptorSupplier
      extends TransportBaseDescriptorSupplier {
    TransportFileDescriptorSupplier() {}
  }

  private static final class TransportMethodDescriptorSupplier
      extends TransportBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoMethodDescriptorSupplier {
    private final String methodName;

    TransportMethodDescriptorSupplier(String methodName) {
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
      synchronized (TransportGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor = result = io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
              .setSchemaDescriptor(new TransportFileDescriptorSupplier())
              .addMethod(getStreamMethod())
              .build();
        }
      }
    }
    return result;
  }
}
