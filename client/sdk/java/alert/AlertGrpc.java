package alert;

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
    comments = "Source: alert/alert.proto")
public final class AlertGrpc {

  private AlertGrpc() {}

  public static final String SERVICE_NAME = "alert.Alert";

  // Static method descriptors that strictly reflect the proto.
  private static volatile io.grpc.MethodDescriptor<alert.AlertOuterClass.ReportEventRequest,
      alert.AlertOuterClass.ReportEventResponse> getReportEventMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "ReportEvent",
      requestType = alert.AlertOuterClass.ReportEventRequest.class,
      responseType = alert.AlertOuterClass.ReportEventResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<alert.AlertOuterClass.ReportEventRequest,
      alert.AlertOuterClass.ReportEventResponse> getReportEventMethod() {
    io.grpc.MethodDescriptor<alert.AlertOuterClass.ReportEventRequest, alert.AlertOuterClass.ReportEventResponse> getReportEventMethod;
    if ((getReportEventMethod = AlertGrpc.getReportEventMethod) == null) {
      synchronized (AlertGrpc.class) {
        if ((getReportEventMethod = AlertGrpc.getReportEventMethod) == null) {
          AlertGrpc.getReportEventMethod = getReportEventMethod =
              io.grpc.MethodDescriptor.<alert.AlertOuterClass.ReportEventRequest, alert.AlertOuterClass.ReportEventResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "ReportEvent"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  alert.AlertOuterClass.ReportEventRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  alert.AlertOuterClass.ReportEventResponse.getDefaultInstance()))
              .setSchemaDescriptor(new AlertMethodDescriptorSupplier("ReportEvent"))
              .build();
        }
      }
    }
    return getReportEventMethod;
  }

  /**
   * Creates a new async stub that supports all call types for the service
   */
  public static AlertStub newStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<AlertStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<AlertStub>() {
        @java.lang.Override
        public AlertStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new AlertStub(channel, callOptions);
        }
      };
    return AlertStub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static AlertBlockingStub newBlockingStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<AlertBlockingStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<AlertBlockingStub>() {
        @java.lang.Override
        public AlertBlockingStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new AlertBlockingStub(channel, callOptions);
        }
      };
    return AlertBlockingStub.newStub(factory, channel);
  }

  /**
   * Creates a new ListenableFuture-style stub that supports unary calls on the service
   */
  public static AlertFutureStub newFutureStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<AlertFutureStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<AlertFutureStub>() {
        @java.lang.Override
        public AlertFutureStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new AlertFutureStub(channel, callOptions);
        }
      };
    return AlertFutureStub.newStub(factory, channel);
  }

  /**
   */
  public static abstract class AlertImplBase implements io.grpc.BindableService {

    /**
     * <pre>
     * ReportEvent does event ingestions.
     * </pre>
     */
    public void reportEvent(alert.AlertOuterClass.ReportEventRequest request,
        io.grpc.stub.StreamObserver<alert.AlertOuterClass.ReportEventResponse> responseObserver) {
      asyncUnimplementedUnaryCall(getReportEventMethod(), responseObserver);
    }

    @java.lang.Override public final io.grpc.ServerServiceDefinition bindService() {
      return io.grpc.ServerServiceDefinition.builder(getServiceDescriptor())
          .addMethod(
            getReportEventMethod(),
            asyncUnaryCall(
              new MethodHandlers<
                alert.AlertOuterClass.ReportEventRequest,
                alert.AlertOuterClass.ReportEventResponse>(
                  this, METHODID_REPORT_EVENT)))
          .build();
    }
  }

  /**
   */
  public static final class AlertStub extends io.grpc.stub.AbstractAsyncStub<AlertStub> {
    private AlertStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected AlertStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new AlertStub(channel, callOptions);
    }

    /**
     * <pre>
     * ReportEvent does event ingestions.
     * </pre>
     */
    public void reportEvent(alert.AlertOuterClass.ReportEventRequest request,
        io.grpc.stub.StreamObserver<alert.AlertOuterClass.ReportEventResponse> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(getReportEventMethod(), getCallOptions()), request, responseObserver);
    }
  }

  /**
   */
  public static final class AlertBlockingStub extends io.grpc.stub.AbstractBlockingStub<AlertBlockingStub> {
    private AlertBlockingStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected AlertBlockingStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new AlertBlockingStub(channel, callOptions);
    }

    /**
     * <pre>
     * ReportEvent does event ingestions.
     * </pre>
     */
    public alert.AlertOuterClass.ReportEventResponse reportEvent(alert.AlertOuterClass.ReportEventRequest request) {
      return blockingUnaryCall(
          getChannel(), getReportEventMethod(), getCallOptions(), request);
    }
  }

  /**
   */
  public static final class AlertFutureStub extends io.grpc.stub.AbstractFutureStub<AlertFutureStub> {
    private AlertFutureStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected AlertFutureStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new AlertFutureStub(channel, callOptions);
    }

    /**
     * <pre>
     * ReportEvent does event ingestions.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<alert.AlertOuterClass.ReportEventResponse> reportEvent(
        alert.AlertOuterClass.ReportEventRequest request) {
      return futureUnaryCall(
          getChannel().newCall(getReportEventMethod(), getCallOptions()), request);
    }
  }

  private static final int METHODID_REPORT_EVENT = 0;

  private static final class MethodHandlers<Req, Resp> implements
      io.grpc.stub.ServerCalls.UnaryMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ServerStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ClientStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.BidiStreamingMethod<Req, Resp> {
    private final AlertImplBase serviceImpl;
    private final int methodId;

    MethodHandlers(AlertImplBase serviceImpl, int methodId) {
      this.serviceImpl = serviceImpl;
      this.methodId = methodId;
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public void invoke(Req request, io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        case METHODID_REPORT_EVENT:
          serviceImpl.reportEvent((alert.AlertOuterClass.ReportEventRequest) request,
              (io.grpc.stub.StreamObserver<alert.AlertOuterClass.ReportEventResponse>) responseObserver);
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

  private static abstract class AlertBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier, io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    AlertBaseDescriptorSupplier() {}

    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return alert.AlertOuterClass.getDescriptor();
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.ServiceDescriptor getServiceDescriptor() {
      return getFileDescriptor().findServiceByName("Alert");
    }
  }

  private static final class AlertFileDescriptorSupplier
      extends AlertBaseDescriptorSupplier {
    AlertFileDescriptorSupplier() {}
  }

  private static final class AlertMethodDescriptorSupplier
      extends AlertBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoMethodDescriptorSupplier {
    private final String methodName;

    AlertMethodDescriptorSupplier(String methodName) {
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
      synchronized (AlertGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor = result = io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
              .setSchemaDescriptor(new AlertFileDescriptorSupplier())
              .addMethod(getReportEventMethod())
              .build();
        }
      }
    }
    return result;
  }
}
