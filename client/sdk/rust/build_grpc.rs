// This file is generated. Do not edit
// @generated

// https://github.com/Manishearth/rust-clippy/issues/702
#![allow(unknown_lints)]
#![allow(clippy::all)]

#![cfg_attr(rustfmt, rustfmt_skip)]

#![allow(box_pointers)]
#![allow(dead_code)]
#![allow(missing_docs)]
#![allow(non_camel_case_types)]
#![allow(non_snake_case)]
#![allow(non_upper_case_globals)]
#![allow(trivial_casts)]
#![allow(unsafe_code)]
#![allow(unused_imports)]
#![allow(unused_results)]


// server interface

pub trait Build {
    fn build(&self, o: ::grpc::ServerHandlerContext, req: ::grpc::ServerRequest<super::build::BuildRequest>, resp: ::grpc::ServerResponseSink<super::build::Result>) -> ::grpc::Result<()>;
}

// client

pub struct BuildClient {
    grpc_client: ::std::sync::Arc<::grpc::Client>,
}

impl ::grpc::ClientStub for BuildClient {
    fn with_client(grpc_client: ::std::sync::Arc<::grpc::Client>) -> Self {
        BuildClient {
            grpc_client: grpc_client,
        }
    }
}

impl BuildClient {
    pub fn build(&self, o: ::grpc::RequestOptions) -> impl ::std::future::Future<Output=::grpc::Result<(::grpc::ClientRequestSink<super::build::BuildRequest>, ::grpc::StreamingResponse<super::build::Result>)>> {
        let descriptor = ::grpc::rt::ArcOrStatic::Static(&::grpc::rt::MethodDescriptor {
            name: ::grpc::rt::StringOrStatic::Static("/build.Build/Build"),
            streaming: ::grpc::rt::GrpcStreaming::Bidi,
            req_marshaller: ::grpc::rt::ArcOrStatic::Static(&::grpc_protobuf::MarshallerProtobuf),
            resp_marshaller: ::grpc::rt::ArcOrStatic::Static(&::grpc_protobuf::MarshallerProtobuf),
        });
        self.grpc_client.call_bidi(o, descriptor)
    }
}

// server

pub struct BuildServer;


impl BuildServer {
    pub fn new_service_def<H : Build + 'static + Sync + Send + 'static>(handler: H) -> ::grpc::rt::ServerServiceDefinition {
        let handler_arc = ::std::sync::Arc::new(handler);
        ::grpc::rt::ServerServiceDefinition::new("/build.Build",
            vec![
                ::grpc::rt::ServerMethod::new(
                    ::grpc::rt::ArcOrStatic::Static(&::grpc::rt::MethodDescriptor {
                        name: ::grpc::rt::StringOrStatic::Static("/build.Build/Build"),
                        streaming: ::grpc::rt::GrpcStreaming::Bidi,
                        req_marshaller: ::grpc::rt::ArcOrStatic::Static(&::grpc_protobuf::MarshallerProtobuf),
                        resp_marshaller: ::grpc::rt::ArcOrStatic::Static(&::grpc_protobuf::MarshallerProtobuf),
                    }),
                    {
                        let handler_copy = handler_arc.clone();
                        ::grpc::rt::MethodHandlerBidi::new(move |ctx, req, resp| (*handler_copy).build(ctx, req, resp))
                    },
                ),
            ],
        )
    }
}
