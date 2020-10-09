package config;

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
    comments = "Source: config/config.proto")
public final class ConfigGrpc {

  private ConfigGrpc() {}

  public static final String SERVICE_NAME = "config.Config";

  // Static method descriptors that strictly reflect the proto.
  private static volatile io.grpc.MethodDescriptor<config.ConfigOuterClass.GetRequest,
      config.ConfigOuterClass.GetResponse> getGetMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "Get",
      requestType = config.ConfigOuterClass.GetRequest.class,
      responseType = config.ConfigOuterClass.GetResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<config.ConfigOuterClass.GetRequest,
      config.ConfigOuterClass.GetResponse> getGetMethod() {
    io.grpc.MethodDescriptor<config.ConfigOuterClass.GetRequest, config.ConfigOuterClass.GetResponse> getGetMethod;
    if ((getGetMethod = ConfigGrpc.getGetMethod) == null) {
      synchronized (ConfigGrpc.class) {
        if ((getGetMethod = ConfigGrpc.getGetMethod) == null) {
          ConfigGrpc.getGetMethod = getGetMethod =
              io.grpc.MethodDescriptor.<config.ConfigOuterClass.GetRequest, config.ConfigOuterClass.GetResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "Get"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  config.ConfigOuterClass.GetRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  config.ConfigOuterClass.GetResponse.getDefaultInstance()))
              .setSchemaDescriptor(new ConfigMethodDescriptorSupplier("Get"))
              .build();
        }
      }
    }
    return getGetMethod;
  }

  private static volatile io.grpc.MethodDescriptor<config.ConfigOuterClass.SetRequest,
      config.ConfigOuterClass.SetResponse> getSetMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "Set",
      requestType = config.ConfigOuterClass.SetRequest.class,
      responseType = config.ConfigOuterClass.SetResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<config.ConfigOuterClass.SetRequest,
      config.ConfigOuterClass.SetResponse> getSetMethod() {
    io.grpc.MethodDescriptor<config.ConfigOuterClass.SetRequest, config.ConfigOuterClass.SetResponse> getSetMethod;
    if ((getSetMethod = ConfigGrpc.getSetMethod) == null) {
      synchronized (ConfigGrpc.class) {
        if ((getSetMethod = ConfigGrpc.getSetMethod) == null) {
          ConfigGrpc.getSetMethod = getSetMethod =
              io.grpc.MethodDescriptor.<config.ConfigOuterClass.SetRequest, config.ConfigOuterClass.SetResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "Set"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  config.ConfigOuterClass.SetRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  config.ConfigOuterClass.SetResponse.getDefaultInstance()))
              .setSchemaDescriptor(new ConfigMethodDescriptorSupplier("Set"))
              .build();
        }
      }
    }
    return getSetMethod;
  }

  private static volatile io.grpc.MethodDescriptor<config.ConfigOuterClass.DeleteRequest,
      config.ConfigOuterClass.DeleteResponse> getDeleteMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "Delete",
      requestType = config.ConfigOuterClass.DeleteRequest.class,
      responseType = config.ConfigOuterClass.DeleteResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<config.ConfigOuterClass.DeleteRequest,
      config.ConfigOuterClass.DeleteResponse> getDeleteMethod() {
    io.grpc.MethodDescriptor<config.ConfigOuterClass.DeleteRequest, config.ConfigOuterClass.DeleteResponse> getDeleteMethod;
    if ((getDeleteMethod = ConfigGrpc.getDeleteMethod) == null) {
      synchronized (ConfigGrpc.class) {
        if ((getDeleteMethod = ConfigGrpc.getDeleteMethod) == null) {
          ConfigGrpc.getDeleteMethod = getDeleteMethod =
              io.grpc.MethodDescriptor.<config.ConfigOuterClass.DeleteRequest, config.ConfigOuterClass.DeleteResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "Delete"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  config.ConfigOuterClass.DeleteRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  config.ConfigOuterClass.DeleteResponse.getDefaultInstance()))
              .setSchemaDescriptor(new ConfigMethodDescriptorSupplier("Delete"))
              .build();
        }
      }
    }
    return getDeleteMethod;
  }

  private static volatile io.grpc.MethodDescriptor<config.ConfigOuterClass.ReadRequest,
      config.ConfigOuterClass.ReadResponse> getReadMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "Read",
      requestType = config.ConfigOuterClass.ReadRequest.class,
      responseType = config.ConfigOuterClass.ReadResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<config.ConfigOuterClass.ReadRequest,
      config.ConfigOuterClass.ReadResponse> getReadMethod() {
    io.grpc.MethodDescriptor<config.ConfigOuterClass.ReadRequest, config.ConfigOuterClass.ReadResponse> getReadMethod;
    if ((getReadMethod = ConfigGrpc.getReadMethod) == null) {
      synchronized (ConfigGrpc.class) {
        if ((getReadMethod = ConfigGrpc.getReadMethod) == null) {
          ConfigGrpc.getReadMethod = getReadMethod =
              io.grpc.MethodDescriptor.<config.ConfigOuterClass.ReadRequest, config.ConfigOuterClass.ReadResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "Read"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  config.ConfigOuterClass.ReadRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  config.ConfigOuterClass.ReadResponse.getDefaultInstance()))
              .setSchemaDescriptor(new ConfigMethodDescriptorSupplier("Read"))
              .build();
        }
      }
    }
    return getReadMethod;
  }

  /**
   * Creates a new async stub that supports all call types for the service
   */
  public static ConfigStub newStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<ConfigStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<ConfigStub>() {
        @java.lang.Override
        public ConfigStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new ConfigStub(channel, callOptions);
        }
      };
    return ConfigStub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static ConfigBlockingStub newBlockingStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<ConfigBlockingStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<ConfigBlockingStub>() {
        @java.lang.Override
        public ConfigBlockingStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new ConfigBlockingStub(channel, callOptions);
        }
      };
    return ConfigBlockingStub.newStub(factory, channel);
  }

  /**
   * Creates a new ListenableFuture-style stub that supports unary calls on the service
   */
  public static ConfigFutureStub newFutureStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<ConfigFutureStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<ConfigFutureStub>() {
        @java.lang.Override
        public ConfigFutureStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new ConfigFutureStub(channel, callOptions);
        }
      };
    return ConfigFutureStub.newStub(factory, channel);
  }

  /**
   */
  public static abstract class ConfigImplBase implements io.grpc.BindableService {

    /**
     */
    public void get(config.ConfigOuterClass.GetRequest request,
        io.grpc.stub.StreamObserver<config.ConfigOuterClass.GetResponse> responseObserver) {
      asyncUnimplementedUnaryCall(getGetMethod(), responseObserver);
    }

    /**
     */
    public void set(config.ConfigOuterClass.SetRequest request,
        io.grpc.stub.StreamObserver<config.ConfigOuterClass.SetResponse> responseObserver) {
      asyncUnimplementedUnaryCall(getSetMethod(), responseObserver);
    }

    /**
     */
    public void delete(config.ConfigOuterClass.DeleteRequest request,
        io.grpc.stub.StreamObserver<config.ConfigOuterClass.DeleteResponse> responseObserver) {
      asyncUnimplementedUnaryCall(getDeleteMethod(), responseObserver);
    }

    /**
     * <pre>
     * These methods are here for backwards compatibility reasons
     * </pre>
     */
    public void read(config.ConfigOuterClass.ReadRequest request,
        io.grpc.stub.StreamObserver<config.ConfigOuterClass.ReadResponse> responseObserver) {
      asyncUnimplementedUnaryCall(getReadMethod(), responseObserver);
    }

    @java.lang.Override public final io.grpc.ServerServiceDefinition bindService() {
      return io.grpc.ServerServiceDefinition.builder(getServiceDescriptor())
          .addMethod(
            getGetMethod(),
            asyncUnaryCall(
              new MethodHandlers<
                config.ConfigOuterClass.GetRequest,
                config.ConfigOuterClass.GetResponse>(
                  this, METHODID_GET)))
          .addMethod(
            getSetMethod(),
            asyncUnaryCall(
              new MethodHandlers<
                config.ConfigOuterClass.SetRequest,
                config.ConfigOuterClass.SetResponse>(
                  this, METHODID_SET)))
          .addMethod(
            getDeleteMethod(),
            asyncUnaryCall(
              new MethodHandlers<
                config.ConfigOuterClass.DeleteRequest,
                config.ConfigOuterClass.DeleteResponse>(
                  this, METHODID_DELETE)))
          .addMethod(
            getReadMethod(),
            asyncUnaryCall(
              new MethodHandlers<
                config.ConfigOuterClass.ReadRequest,
                config.ConfigOuterClass.ReadResponse>(
                  this, METHODID_READ)))
          .build();
    }
  }

  /**
   */
  public static final class ConfigStub extends io.grpc.stub.AbstractAsyncStub<ConfigStub> {
    private ConfigStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected ConfigStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new ConfigStub(channel, callOptions);
    }

    /**
     */
    public void get(config.ConfigOuterClass.GetRequest request,
        io.grpc.stub.StreamObserver<config.ConfigOuterClass.GetResponse> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(getGetMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void set(config.ConfigOuterClass.SetRequest request,
        io.grpc.stub.StreamObserver<config.ConfigOuterClass.SetResponse> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(getSetMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void delete(config.ConfigOuterClass.DeleteRequest request,
        io.grpc.stub.StreamObserver<config.ConfigOuterClass.DeleteResponse> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(getDeleteMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * These methods are here for backwards compatibility reasons
     * </pre>
     */
    public void read(config.ConfigOuterClass.ReadRequest request,
        io.grpc.stub.StreamObserver<config.ConfigOuterClass.ReadResponse> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(getReadMethod(), getCallOptions()), request, responseObserver);
    }
  }

  /**
   */
  public static final class ConfigBlockingStub extends io.grpc.stub.AbstractBlockingStub<ConfigBlockingStub> {
    private ConfigBlockingStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected ConfigBlockingStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new ConfigBlockingStub(channel, callOptions);
    }

    /**
     */
    public config.ConfigOuterClass.GetResponse get(config.ConfigOuterClass.GetRequest request) {
      return blockingUnaryCall(
          getChannel(), getGetMethod(), getCallOptions(), request);
    }

    /**
     */
    public config.ConfigOuterClass.SetResponse set(config.ConfigOuterClass.SetRequest request) {
      return blockingUnaryCall(
          getChannel(), getSetMethod(), getCallOptions(), request);
    }

    /**
     */
    public config.ConfigOuterClass.DeleteResponse delete(config.ConfigOuterClass.DeleteRequest request) {
      return blockingUnaryCall(
          getChannel(), getDeleteMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * These methods are here for backwards compatibility reasons
     * </pre>
     */
    public config.ConfigOuterClass.ReadResponse read(config.ConfigOuterClass.ReadRequest request) {
      return blockingUnaryCall(
          getChannel(), getReadMethod(), getCallOptions(), request);
    }
  }

  /**
   */
  public static final class ConfigFutureStub extends io.grpc.stub.AbstractFutureStub<ConfigFutureStub> {
    private ConfigFutureStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected ConfigFutureStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new ConfigFutureStub(channel, callOptions);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<config.ConfigOuterClass.GetResponse> get(
        config.ConfigOuterClass.GetRequest request) {
      return futureUnaryCall(
          getChannel().newCall(getGetMethod(), getCallOptions()), request);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<config.ConfigOuterClass.SetResponse> set(
        config.ConfigOuterClass.SetRequest request) {
      return futureUnaryCall(
          getChannel().newCall(getSetMethod(), getCallOptions()), request);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<config.ConfigOuterClass.DeleteResponse> delete(
        config.ConfigOuterClass.DeleteRequest request) {
      return futureUnaryCall(
          getChannel().newCall(getDeleteMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * These methods are here for backwards compatibility reasons
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<config.ConfigOuterClass.ReadResponse> read(
        config.ConfigOuterClass.ReadRequest request) {
      return futureUnaryCall(
          getChannel().newCall(getReadMethod(), getCallOptions()), request);
    }
  }

  private static final int METHODID_GET = 0;
  private static final int METHODID_SET = 1;
  private static final int METHODID_DELETE = 2;
  private static final int METHODID_READ = 3;

  private static final class MethodHandlers<Req, Resp> implements
      io.grpc.stub.ServerCalls.UnaryMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ServerStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ClientStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.BidiStreamingMethod<Req, Resp> {
    private final ConfigImplBase serviceImpl;
    private final int methodId;

    MethodHandlers(ConfigImplBase serviceImpl, int methodId) {
      this.serviceImpl = serviceImpl;
      this.methodId = methodId;
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public void invoke(Req request, io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        case METHODID_GET:
          serviceImpl.get((config.ConfigOuterClass.GetRequest) request,
              (io.grpc.stub.StreamObserver<config.ConfigOuterClass.GetResponse>) responseObserver);
          break;
        case METHODID_SET:
          serviceImpl.set((config.ConfigOuterClass.SetRequest) request,
              (io.grpc.stub.StreamObserver<config.ConfigOuterClass.SetResponse>) responseObserver);
          break;
        case METHODID_DELETE:
          serviceImpl.delete((config.ConfigOuterClass.DeleteRequest) request,
              (io.grpc.stub.StreamObserver<config.ConfigOuterClass.DeleteResponse>) responseObserver);
          break;
        case METHODID_READ:
          serviceImpl.read((config.ConfigOuterClass.ReadRequest) request,
              (io.grpc.stub.StreamObserver<config.ConfigOuterClass.ReadResponse>) responseObserver);
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

  private static abstract class ConfigBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier, io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    ConfigBaseDescriptorSupplier() {}

    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return config.ConfigOuterClass.getDescriptor();
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.ServiceDescriptor getServiceDescriptor() {
      return getFileDescriptor().findServiceByName("Config");
    }
  }

  private static final class ConfigFileDescriptorSupplier
      extends ConfigBaseDescriptorSupplier {
    ConfigFileDescriptorSupplier() {}
  }

  private static final class ConfigMethodDescriptorSupplier
      extends ConfigBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoMethodDescriptorSupplier {
    private final String methodName;

    ConfigMethodDescriptorSupplier(String methodName) {
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
      synchronized (ConfigGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor = result = io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
              .setSchemaDescriptor(new ConfigFileDescriptorSupplier())
              .addMethod(getGetMethod())
              .addMethod(getSetMethod())
              .addMethod(getDeleteMethod())
              .addMethod(getReadMethod())
              .build();
        }
      }
    }
    return result;
  }
}
