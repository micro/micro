// GENERATED CODE -- DO NOT EDIT!

'use strict';
var grpc = require('grpc');
var debug_debug_pb = require('../debug/debug_pb.js');

function serialize_debug_HealthRequest(arg) {
  if (!(arg instanceof debug_debug_pb.HealthRequest)) {
    throw new Error('Expected argument of type debug.HealthRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_debug_HealthRequest(buffer_arg) {
  return debug_debug_pb.HealthRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_debug_HealthResponse(arg) {
  if (!(arg instanceof debug_debug_pb.HealthResponse)) {
    throw new Error('Expected argument of type debug.HealthResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_debug_HealthResponse(buffer_arg) {
  return debug_debug_pb.HealthResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_debug_LogRequest(arg) {
  if (!(arg instanceof debug_debug_pb.LogRequest)) {
    throw new Error('Expected argument of type debug.LogRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_debug_LogRequest(buffer_arg) {
  return debug_debug_pb.LogRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_debug_LogResponse(arg) {
  if (!(arg instanceof debug_debug_pb.LogResponse)) {
    throw new Error('Expected argument of type debug.LogResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_debug_LogResponse(buffer_arg) {
  return debug_debug_pb.LogResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_debug_StatsRequest(arg) {
  if (!(arg instanceof debug_debug_pb.StatsRequest)) {
    throw new Error('Expected argument of type debug.StatsRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_debug_StatsRequest(buffer_arg) {
  return debug_debug_pb.StatsRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_debug_StatsResponse(arg) {
  if (!(arg instanceof debug_debug_pb.StatsResponse)) {
    throw new Error('Expected argument of type debug.StatsResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_debug_StatsResponse(buffer_arg) {
  return debug_debug_pb.StatsResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_debug_TraceRequest(arg) {
  if (!(arg instanceof debug_debug_pb.TraceRequest)) {
    throw new Error('Expected argument of type debug.TraceRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_debug_TraceRequest(buffer_arg) {
  return debug_debug_pb.TraceRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_debug_TraceResponse(arg) {
  if (!(arg instanceof debug_debug_pb.TraceResponse)) {
    throw new Error('Expected argument of type debug.TraceResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_debug_TraceResponse(buffer_arg) {
  return debug_debug_pb.TraceResponse.deserializeBinary(new Uint8Array(buffer_arg));
}


var DebugService = exports.DebugService = {
  log: {
    path: '/debug.Debug/Log',
    requestStream: false,
    responseStream: false,
    requestType: debug_debug_pb.LogRequest,
    responseType: debug_debug_pb.LogResponse,
    requestSerialize: serialize_debug_LogRequest,
    requestDeserialize: deserialize_debug_LogRequest,
    responseSerialize: serialize_debug_LogResponse,
    responseDeserialize: deserialize_debug_LogResponse,
  },
  health: {
    path: '/debug.Debug/Health',
    requestStream: false,
    responseStream: false,
    requestType: debug_debug_pb.HealthRequest,
    responseType: debug_debug_pb.HealthResponse,
    requestSerialize: serialize_debug_HealthRequest,
    requestDeserialize: deserialize_debug_HealthRequest,
    responseSerialize: serialize_debug_HealthResponse,
    responseDeserialize: deserialize_debug_HealthResponse,
  },
  stats: {
    path: '/debug.Debug/Stats',
    requestStream: false,
    responseStream: false,
    requestType: debug_debug_pb.StatsRequest,
    responseType: debug_debug_pb.StatsResponse,
    requestSerialize: serialize_debug_StatsRequest,
    requestDeserialize: deserialize_debug_StatsRequest,
    responseSerialize: serialize_debug_StatsResponse,
    responseDeserialize: deserialize_debug_StatsResponse,
  },
  trace: {
    path: '/debug.Debug/Trace',
    requestStream: false,
    responseStream: false,
    requestType: debug_debug_pb.TraceRequest,
    responseType: debug_debug_pb.TraceResponse,
    requestSerialize: serialize_debug_TraceRequest,
    requestDeserialize: deserialize_debug_TraceRequest,
    responseSerialize: serialize_debug_TraceResponse,
    responseDeserialize: deserialize_debug_TraceResponse,
  },
};

exports.DebugClient = grpc.makeGenericClientConstructor(DebugService);
