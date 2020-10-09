// GENERATED CODE -- DO NOT EDIT!

'use strict';
var grpc = require('grpc');
var api_api_pb = require('../api/api_pb.js');

function serialize_api_EmptyResponse(arg) {
  if (!(arg instanceof api_api_pb.EmptyResponse)) {
    throw new Error('Expected argument of type api.EmptyResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_api_EmptyResponse(buffer_arg) {
  return api_api_pb.EmptyResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_api_Endpoint(arg) {
  if (!(arg instanceof api_api_pb.Endpoint)) {
    throw new Error('Expected argument of type api.Endpoint');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_api_Endpoint(buffer_arg) {
  return api_api_pb.Endpoint.deserializeBinary(new Uint8Array(buffer_arg));
}


var ApiService = exports.ApiService = {
  register: {
    path: '/api.Api/Register',
    requestStream: false,
    responseStream: false,
    requestType: api_api_pb.Endpoint,
    responseType: api_api_pb.EmptyResponse,
    requestSerialize: serialize_api_Endpoint,
    requestDeserialize: deserialize_api_Endpoint,
    responseSerialize: serialize_api_EmptyResponse,
    responseDeserialize: deserialize_api_EmptyResponse,
  },
  deregister: {
    path: '/api.Api/Deregister',
    requestStream: false,
    responseStream: false,
    requestType: api_api_pb.Endpoint,
    responseType: api_api_pb.EmptyResponse,
    requestSerialize: serialize_api_Endpoint,
    requestDeserialize: deserialize_api_Endpoint,
    responseSerialize: serialize_api_EmptyResponse,
    responseDeserialize: deserialize_api_EmptyResponse,
  },
};

exports.ApiClient = grpc.makeGenericClientConstructor(ApiService);
