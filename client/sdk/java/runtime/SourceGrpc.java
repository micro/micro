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
 * <pre>
 * Source service is used by the CLI to upload source to the service. The service will return
 * a unique ID representing the location of that source. This ID can then be used as a source
 * for the service when doing Runtime.Create. The server will handle cleanup of uploaded source.
 * </pre>
 */
@javax.annotation.Generated(
    value = "by gRPC proto compiler (version 1.28.0)",
    comments = "Source: runtime/runtime.proto")
public final class SourceGrpc {

  private SourceGrpc() {}

  public static final String SERVICE_NAME = "runtime.Source";

  // Static method descriptors that strictly reflect the proto.
  private static volatile io.grpc.MethodDescriptor<runtime.RuntimeOuterClass.UploadRequest,
      runtime.RuntimeOuterClass.UploadResponse> getUploadMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "Upload",
      requestType = runtime.RuntimeOuterClass.UploadRequest.class,
      responseType = runtime.RuntimeOuterClass.UploadResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.CLIENT_STREAMING)
  public static io.grpc.MethodDescriptor<runtime.RuntimeOuterClass.UploadRequest,
      runtime.RuntimeOuterClass.UploadResponse> getUploadMethod() {
    io.grpc.MethodDescriptor<runtime.RuntimeOuterClass.UploadRequest, runtime.RuntimeOuterClass.UploadResponse> getUploadMethod;
    if ((getUploadMethod = SourceGrpc.getUploadMethod) == null) {
      synchronized (SourceGrpc.class) {
        if ((getUploadMethod = SourceGrpc.getUploadMethod) == null) {
          SourceGrpc.getUploadMethod = getUploadMethod =
              io.grpc.MethodDescriptor.<runtime.RuntimeOuterClass.UploadRequest, runtime.RuntimeOuterClass.UploadResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.CLIENT_STREAMING)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "Upload"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  runtime.RuntimeOuterClass.UploadRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  runtime.RuntimeOuterClass.UploadResponse.getDefaultInstance()))
              .setSchemaDescriptor(new SourceMethodDescriptorSupplier("Upload"))
              .build();
        }
      }
    }
    return getUploadMethod;
  }

  /**
   * Creates a new async stub that supports all call types for the service
   */
  public static SourceStub newStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<SourceStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<SourceStub>() {
        @java.lang.Override
        public SourceStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new SourceStub(channel, callOptions);
        }
      };
    return SourceStub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static SourceBlockingStub newBlockingStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<SourceBlockingStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<SourceBlockingStub>() {
        @java.lang.Override
        public SourceBlockingStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new SourceBlockingStub(channel, callOptions);
        }
      };
    return SourceBlockingStub.newStub(factory, channel);
  }

  /**
   * Creates a new ListenableFuture-style stub that supports unary calls on the service
   */
  public static SourceFutureStub newFutureStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<SourceFutureStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<SourceFutureStub>() {
        @java.lang.Override
        public SourceFutureStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new SourceFutureStub(channel, callOptions);
        }
      };
    return SourceFutureStub.newStub(factory, channel);
  }

  /**
   * <pre>
   * Source service is used by the CLI to upload source to the service. The service will return
   * a unique ID representing the location of that source. This ID can then be used as a source
   * for the service when doing Runtime.Create. The server will handle cleanup of uploaded source.
   * </pre>
   */
  public static abstract class SourceImplBase implements io.grpc.BindableService {

    /**
     */
    public io.grpc.stub.StreamObserver<runtime.RuntimeOuterClass.UploadRequest> upload(
        io.grpc.stub.StreamObserver<runtime.RuntimeOuterClass.UploadResponse> responseObserver) {
      return asyncUnimplementedStreamingCall(getUploadMethod(), responseObserver);
    }

    @java.lang.Override public final io.grpc.ServerServiceDefinition bindService() {
      return io.grpc.ServerServiceDefinition.builder(getServiceDescriptor())
          .addMethod(
            getUploadMethod(),
            asyncClientStreamingCall(
              new MethodHandlers<
                runtime.RuntimeOuterClass.UploadRequest,
                runtime.RuntimeOuterClass.UploadResponse>(
                  this, METHODID_UPLOAD)))
          .build();
    }
  }

  /**
   * <pre>
   * Source service is used by the CLI to upload source to the service. The service will return
   * a unique ID representing the location of that source. This ID can then be used as a source
   * for the service when doing Runtime.Create. The server will handle cleanup of uploaded source.
   * </pre>
   */
  public static final class SourceStub extends io.grpc.stub.AbstractAsyncStub<SourceStub> {
    private SourceStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected SourceStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new SourceStub(channel, callOptions);
    }

    /**
     */
    public io.grpc.stub.StreamObserver<runtime.RuntimeOuterClass.UploadRequest> upload(
        io.grpc.stub.StreamObserver<runtime.RuntimeOuterClass.UploadResponse> responseObserver) {
      return asyncClientStreamingCall(
          getChannel().newCall(getUploadMethod(), getCallOptions()), responseObserver);
    }
  }

  /**
   * <pre>
   * Source service is used by the CLI to upload source to the service. The service will return
   * a unique ID representing the location of that source. This ID can then be used as a source
   * for the service when doing Runtime.Create. The server will handle cleanup of uploaded source.
   * </pre>
   */
  public static final class SourceBlockingStub extends io.grpc.stub.AbstractBlockingStub<SourceBlockingStub> {
    private SourceBlockingStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected SourceBlockingStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new SourceBlockingStub(channel, callOptions);
    }
  }

  /**
   * <pre>
   * Source service is used by the CLI to upload source to the service. The service will return
   * a unique ID representing the location of that source. This ID can then be used as a source
   * for the service when doing Runtime.Create. The server will handle cleanup of uploaded source.
   * </pre>
   */
  public static final class SourceFutureStub extends io.grpc.stub.AbstractFutureStub<SourceFutureStub> {
    private SourceFutureStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected SourceFutureStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new SourceFutureStub(channel, callOptions);
    }
  }

  private static final int METHODID_UPLOAD = 0;

  private static final class MethodHandlers<Req, Resp> implements
      io.grpc.stub.ServerCalls.UnaryMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ServerStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ClientStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.BidiStreamingMethod<Req, Resp> {
    private final SourceImplBase serviceImpl;
    private final int methodId;

    MethodHandlers(SourceImplBase serviceImpl, int methodId) {
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
        case METHODID_UPLOAD:
          return (io.grpc.stub.StreamObserver<Req>) serviceImpl.upload(
              (io.grpc.stub.StreamObserver<runtime.RuntimeOuterClass.UploadResponse>) responseObserver);
        default:
          throw new AssertionError();
      }
    }
  }

  private static abstract class SourceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier, io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    SourceBaseDescriptorSupplier() {}

    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return runtime.RuntimeOuterClass.getDescriptor();
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.ServiceDescriptor getServiceDescriptor() {
      return getFileDescriptor().findServiceByName("Source");
    }
  }

  private static final class SourceFileDescriptorSupplier
      extends SourceBaseDescriptorSupplier {
    SourceFileDescriptorSupplier() {}
  }

  private static final class SourceMethodDescriptorSupplier
      extends SourceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoMethodDescriptorSupplier {
    private final String methodName;

    SourceMethodDescriptorSupplier(String methodName) {
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
      synchronized (SourceGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor = result = io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
              .setSchemaDescriptor(new SourceFileDescriptorSupplier())
              .addMethod(getUploadMethod())
              .build();
        }
      }
    }
    return result;
  }
}
