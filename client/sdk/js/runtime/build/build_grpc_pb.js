// GENERATED CODE -- DO NOT EDIT!

'use strict';
var grpc = require('grpc');
var runtime_build_build_pb = require('../../runtime/build/build_pb.js');

function serialize_runtime_build_BuildRequest(arg) {
  if (!(arg instanceof runtime_build_build_pb.BuildRequest)) {
    throw new Error('Expected argument of type runtime.build.BuildRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_runtime_build_BuildRequest(buffer_arg) {
  return runtime_build_build_pb.BuildRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_runtime_build_Result(arg) {
  if (!(arg instanceof runtime_build_build_pb.Result)) {
    throw new Error('Expected argument of type runtime.build.Result');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_runtime_build_Result(buffer_arg) {
  return runtime_build_build_pb.Result.deserializeBinary(new Uint8Array(buffer_arg));
}


var BuildService = exports.BuildService = {
  build: {
    path: '/runtime.build.Build/Build',
    requestStream: true,
    responseStream: true,
    requestType: runtime_build_build_pb.BuildRequest,
    responseType: runtime_build_build_pb.Result,
    requestSerialize: serialize_runtime_build_BuildRequest,
    requestDeserialize: deserialize_runtime_build_BuildRequest,
    responseSerialize: serialize_runtime_build_Result,
    responseDeserialize: deserialize_runtime_build_Result,
  },
};

exports.BuildClient = grpc.makeGenericClientConstructor(BuildService);
