package build;

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
    comments = "Source: build/build.proto")
public final class BuildGrpc {

  private BuildGrpc() {}

  public static final String SERVICE_NAME = "build.Build";

  // Static method descriptors that strictly reflect the proto.
  private static volatile io.grpc.MethodDescriptor<build.BuildOuterClass.BuildRequest,
      build.BuildOuterClass.Result> getBuildMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "Build",
      requestType = build.BuildOuterClass.BuildRequest.class,
      responseType = build.BuildOuterClass.Result.class,
      methodType = io.grpc.MethodDescriptor.MethodType.BIDI_STREAMING)
  public static io.grpc.MethodDescriptor<build.BuildOuterClass.BuildRequest,
      build.BuildOuterClass.Result> getBuildMethod() {
    io.grpc.MethodDescriptor<build.BuildOuterClass.BuildRequest, build.BuildOuterClass.Result> getBuildMethod;
    if ((getBuildMethod = BuildGrpc.getBuildMethod) == null) {
      synchronized (BuildGrpc.class) {
        if ((getBuildMethod = BuildGrpc.getBuildMethod) == null) {
          BuildGrpc.getBuildMethod = getBuildMethod =
              io.grpc.MethodDescriptor.<build.BuildOuterClass.BuildRequest, build.BuildOuterClass.Result>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.BIDI_STREAMING)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "Build"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  build.BuildOuterClass.BuildRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  build.BuildOuterClass.Result.getDefaultInstance()))
              .setSchemaDescriptor(new BuildMethodDescriptorSupplier("Build"))
              .build();
        }
      }
    }
    return getBuildMethod;
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
   */
  public static abstract class BuildImplBase implements io.grpc.BindableService {

    /**
     */
    public io.grpc.stub.StreamObserver<build.BuildOuterClass.BuildRequest> build(
        io.grpc.stub.StreamObserver<build.BuildOuterClass.Result> responseObserver) {
      return asyncUnimplementedStreamingCall(getBuildMethod(), responseObserver);
    }

    @java.lang.Override public final io.grpc.ServerServiceDefinition bindService() {
      return io.grpc.ServerServiceDefinition.builder(getServiceDescriptor())
          .addMethod(
            getBuildMethod(),
            asyncBidiStreamingCall(
              new MethodHandlers<
                build.BuildOuterClass.BuildRequest,
                build.BuildOuterClass.Result>(
                  this, METHODID_BUILD)))
          .build();
    }
  }

  /**
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
    public io.grpc.stub.StreamObserver<build.BuildOuterClass.BuildRequest> build(
        io.grpc.stub.StreamObserver<build.BuildOuterClass.Result> responseObserver) {
      return asyncBidiStreamingCall(
          getChannel().newCall(getBuildMethod(), getCallOptions()), responseObserver);
    }
  }

  /**
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
  }

  /**
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

  private static final int METHODID_BUILD = 0;

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
        default:
          throw new AssertionError();
      }
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public io.grpc.stub.StreamObserver<Req> invoke(
        io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        case METHODID_BUILD:
          return (io.grpc.stub.StreamObserver<Req>) serviceImpl.build(
              (io.grpc.stub.StreamObserver<build.BuildOuterClass.Result>) responseObserver);
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
      return build.BuildOuterClass.getDescriptor();
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
              .addMethod(getBuildMethod())
              .build();
        }
      }
    }
    return result;
  }
}
