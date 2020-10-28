// GENERATED CODE -- DO NOT EDIT!

'use strict';
var grpc = require('grpc');
var config_config_pb = require('../config/config_pb.js');

function serialize_config_DeleteRequest(arg) {
  if (!(arg instanceof config_config_pb.DeleteRequest)) {
    throw new Error('Expected argument of type config.DeleteRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_config_DeleteRequest(buffer_arg) {
  return config_config_pb.DeleteRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_config_DeleteResponse(arg) {
  if (!(arg instanceof config_config_pb.DeleteResponse)) {
    throw new Error('Expected argument of type config.DeleteResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_config_DeleteResponse(buffer_arg) {
  return config_config_pb.DeleteResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_config_GetRequest(arg) {
  if (!(arg instanceof config_config_pb.GetRequest)) {
    throw new Error('Expected argument of type config.GetRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_config_GetRequest(buffer_arg) {
  return config_config_pb.GetRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_config_GetResponse(arg) {
  if (!(arg instanceof config_config_pb.GetResponse)) {
    throw new Error('Expected argument of type config.GetResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_config_GetResponse(buffer_arg) {
  return config_config_pb.GetResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_config_ReadRequest(arg) {
  if (!(arg instanceof config_config_pb.ReadRequest)) {
    throw new Error('Expected argument of type config.ReadRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_config_ReadRequest(buffer_arg) {
  return config_config_pb.ReadRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_config_ReadResponse(arg) {
  if (!(arg instanceof config_config_pb.ReadResponse)) {
    throw new Error('Expected argument of type config.ReadResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_config_ReadResponse(buffer_arg) {
  return config_config_pb.ReadResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_config_SetRequest(arg) {
  if (!(arg instanceof config_config_pb.SetRequest)) {
    throw new Error('Expected argument of type config.SetRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_config_SetRequest(buffer_arg) {
  return config_config_pb.SetRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_config_SetResponse(arg) {
  if (!(arg instanceof config_config_pb.SetResponse)) {
    throw new Error('Expected argument of type config.SetResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_config_SetResponse(buffer_arg) {
  return config_config_pb.SetResponse.deserializeBinary(new Uint8Array(buffer_arg));
}


var ConfigService = exports.ConfigService = {
  get: {
    path: '/config.Config/Get',
    requestStream: false,
    responseStream: false,
    requestType: config_config_pb.GetRequest,
    responseType: config_config_pb.GetResponse,
    requestSerialize: serialize_config_GetRequest,
    requestDeserialize: deserialize_config_GetRequest,
    responseSerialize: serialize_config_GetResponse,
    responseDeserialize: deserialize_config_GetResponse,
  },
  set: {
    path: '/config.Config/Set',
    requestStream: false,
    responseStream: false,
    requestType: config_config_pb.SetRequest,
    responseType: config_config_pb.SetResponse,
    requestSerialize: serialize_config_SetRequest,
    requestDeserialize: deserialize_config_SetRequest,
    responseSerialize: serialize_config_SetResponse,
    responseDeserialize: deserialize_config_SetResponse,
  },
  delete: {
    path: '/config.Config/Delete',
    requestStream: false,
    responseStream: false,
    requestType: config_config_pb.DeleteRequest,
    responseType: config_config_pb.DeleteResponse,
    requestSerialize: serialize_config_DeleteRequest,
    requestDeserialize: deserialize_config_DeleteRequest,
    responseSerialize: serialize_config_DeleteResponse,
    responseDeserialize: deserialize_config_DeleteResponse,
  },
  // These methods are here for backwards compatibility reasons
read: {
    path: '/config.Config/Read',
    requestStream: false,
    responseStream: false,
    requestType: config_config_pb.ReadRequest,
    responseType: config_config_pb.ReadResponse,
    requestSerialize: serialize_config_ReadRequest,
    requestDeserialize: deserialize_config_ReadRequest,
    responseSerialize: serialize_config_ReadResponse,
    responseDeserialize: deserialize_config_ReadResponse,
  },
};

exports.ConfigClient = grpc.makeGenericClientConstructor(ConfigService);
