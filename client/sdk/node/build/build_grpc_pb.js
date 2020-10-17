// GENERATED CODE -- DO NOT EDIT!

'use strict';
var grpc = require('grpc');
var build_build_pb = require('../build/build_pb.js');

function serialize_build_BuildRequest(arg) {
  if (!(arg instanceof build_build_pb.BuildRequest)) {
    throw new Error('Expected argument of type build.BuildRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_build_BuildRequest(buffer_arg) {
  return build_build_pb.BuildRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_build_Result(arg) {
  if (!(arg instanceof build_build_pb.Result)) {
    throw new Error('Expected argument of type build.Result');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_build_Result(buffer_arg) {
  return build_build_pb.Result.deserializeBinary(new Uint8Array(buffer_arg));
}


var BuildService = exports.BuildService = {
  build: {
    path: '/build.Build/Build',
    requestStream: true,
    responseStream: true,
    requestType: build_build_pb.BuildRequest,
    responseType: build_build_pb.Result,
    requestSerialize: serialize_build_BuildRequest,
    requestDeserialize: deserialize_build_BuildRequest,
    responseSerialize: serialize_build_Result,
    responseDeserialize: deserialize_build_Result,
  },
};

exports.BuildClient = grpc.makeGenericClientConstructor(BuildService);
