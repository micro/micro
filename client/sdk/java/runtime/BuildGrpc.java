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
 * Build service is used by containers to download prebuilt binaries. The client will pass the 
 * service (name and version are required attributed) and the server will then stream the latest
 * binary to the client.
 * </pre>
 */
@javax.annotation.Generated(
    value = "by gRPC proto compiler (version 1.28.0)",
    comments = "Source: runtime/runtime.proto")
public final class BuildGrpc {

  private BuildGrpc() {}

  public static final String SERVICE_NAME = "runtime.Build";

  // Static method descriptors that strictly reflect the proto.
  private static volatile io.grpc.MethodDescriptor<runtime.RuntimeOuterClass.Service,
      runtime.RuntimeOuterClass.BuildReadResponse> getReadMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "Read",
      requestType = runtime.RuntimeOuterClass.Service.class,
      responseType = runtime.RuntimeOuterClass.BuildReadResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.SERVER_STREAMING)
  public static io.grpc.MethodDescriptor<runtime.RuntimeOuterClass.Service,
      runtime.RuntimeOuterClass.BuildReadResponse> getReadMethod() {
    io.grpc.MethodDescriptor<runtime.RuntimeOuterClass.Service, runtime.RuntimeOuterClass.BuildReadResponse> getReadMethod;
    if ((getReadMethod = BuildGrpc.getReadMethod) == null) {
      synchronized (BuildGrpc.class) {
        if ((getReadMethod = BuildGrpc.getReadMethod) == null) {
          BuildGrpc.getReadMethod = getReadMethod =
              io.grpc.MethodDescriptor.<runtime.RuntimeOuterClass.Service, runtime.RuntimeOuterClass.BuildReadResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.SERVER_STREAMING)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "Read"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  runtime.RuntimeOuterClass.Service.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  runtime.RuntimeOuterClass.BuildReadResponse.getDefaultInstance()))
              .setSchemaDescriptor(new BuildMethodDescriptorSupplier("Read"))
              .build();
        }
      }
    }
    return getReadMethod;
  }

  /**
   * Creates a new async stub that supports all call types for the service
   */
  public static BuildStub newStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<BuildStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<BuildStub>() {
        @java.lang.Override
        public BuildStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new BuildStub(channel, callOptions);
        }
      };
    return BuildStub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static BuildBlockingStub newBlockingStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<BuildBlockingStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<BuildBlockingStub>() {
        @java.lang.Override
        public BuildBlockingStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new BuildBlockingStub(channel, callOptions);
        }
      };
    return BuildBlockingStub.newStub(factory, channel);
  }

  /**
   * Creates a new ListenableFuture-style stub that supports unary calls on the service
   */
  public static BuildFutureStub newFutureStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<BuildFutureStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<BuildFutureStub>() {
        @java.lang.Override
        public BuildFutureStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new BuildFutureStub(channel, callOptions);
        }
      };
    return BuildFutureStub.newStub(factory, channel);
  }

  /**
   * <pre>
   * Build service is used by containers to download prebuilt binaries. The client will pass the 
   * service (name and version are required attributed) and the server will then stream the latest
   * binary to the client.
   * </pre>
   */
  public static abstract class BuildImplBase implements io.grpc.BindableService {

    /**
     */
    public void read(runtime.RuntimeOuterClass.Service request,
        io.grpc.stub.StreamObserver<runtime.RuntimeOuterClass.BuildReadResponse> responseObserver) {
      asyncUnimplementedUnaryCall(getReadMethod(), responseObserver);
    }

    @java.lang.Override public final io.grpc.ServerServiceDefinition bindService() {
      return io.grpc.ServerServiceDefinition.builder(getServiceDescriptor())
          .addMethod(
            getReadMethod(),
            asyncServerStreamingCall(
              new MethodHandlers<
                runtime.RuntimeOuterClass.Service,
                runtime.RuntimeOuterClass.BuildReadResponse>(
                  this, METHODID_READ)))
          .build();
    }
  }

  /**
   * <pre>
   * Build service is used by containers to download prebuilt binaries. The client will pass the 
   * service (name and version are required attributed) and the server will then stream the latest
   * binary to the client.
   * </pre>
   */
  public static final class BuildStub extends io.grpc.stub.AbstractAsyncStub<BuildStub> {
    private BuildStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected BuildStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new BuildStub(channel, callOptions);
    }

    /**
     */
    public void read(runtime.RuntimeOuterClass.Service request,
        io.grpc.stub.StreamObserver<runtime.RuntimeOuterClass.BuildReadResponse> responseObserver) {
      asyncServerStreamingCall(
          getChannel().newCall(getReadMethod(), getCallOptions()), request, responseObserver);
    }
  }

  /**
   * <pre>
   * Build service is used by containers to download prebuilt binaries. The client will pass the 
   * service (name and version are required attributed) and the server will then stream the latest
   * binary to the client.
   * </pre>
   */
  public static final class BuildBlockingStub extends io.grpc.stub.AbstractBlockingStub<BuildBlockingStub> {
    private BuildBlockingStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected BuildBlockingStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new BuildBlockingStub(channel, callOptions);
    }

    /**
     */
    public java.util.Iterator<runtime.RuntimeOuterClass.BuildReadResponse> read(
        runtime.RuntimeOuterClass.Service request) {
      return blockingServerStreamingCall(
          getChannel(), getReadMethod(), getCallOptions(), request);
    }
  }

  /**
   * <pre>
   * Build service is used by containers to download prebuilt binaries. The client will pass the 
   * service (name and version are required attributed) and the server will then stream the latest
   * binary to the client.
   * </pre>
   */
  public static final class BuildFutureStub extends io.grpc.stub.AbstractFutureStub<BuildFutureStub> {
    private BuildFutureStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected BuildFutureStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new BuildFutureStub(channel, callOptions);
    }
  }

  private static final int METHODID_READ = 0;

  private static final class MethodHandlers<Req, Resp> implements
      io.grpc.stub.ServerCalls.UnaryMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ServerStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ClientStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.BidiStreamingMethod<Req, Resp> {
    private final BuildImplBase serviceImpl;
    private final int methodId;

    MethodHandlers(BuildImplBase serviceImpl, int methodId) {
      this.serviceImpl = serviceImpl;
      this.methodId = methodId;
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public void invoke(Req request, io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        case METHODID_READ:
          serviceImpl.read((runtime.RuntimeOuterClass.Service) request,
              (io.grpc.stub.StreamObserver<runtime.RuntimeOuterClass.BuildReadResponse>) responseObserver);
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

  private static abstract class BuildBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier, io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    BuildBaseDescriptorSupplier() {}

    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return runtime.RuntimeOuterClass.getDescriptor();
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.ServiceDescriptor getServiceDescriptor() {
      return getFileDescriptor().findServiceByName("Build");
    }
  }

  private static final class BuildFileDescriptorSupplier
      extends BuildBaseDescriptorSupplier {
    BuildFileDescriptorSupplier() {}
  }

  private static final class BuildMethodDescriptorSupplier
      extends BuildBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoMethodDescriptorSupplier {
    private final String methodName;

    BuildMethodDescriptorSupplier(String methodName) {
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
      synchronized (BuildGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor = result = io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
              .setSchemaDescriptor(new BuildFileDescriptorSupplier())
              .addMethod(getReadMethod())
              .build();
        }
      }
    }
    return result;
  }
}
