// GENERATED CODE -- DO NOT EDIT!

'use strict';
var grpc = require('grpc');
var events_events_pb = require('../events/events_pb.js');

function serialize_events_ConsumeRequest(arg) {
  if (!(arg instanceof events_events_pb.ConsumeRequest)) {
    throw new Error('Expected argument of type events.ConsumeRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_events_ConsumeRequest(buffer_arg) {
  return events_events_pb.ConsumeRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_events_Event(arg) {
  if (!(arg instanceof events_events_pb.Event)) {
    throw new Error('Expected argument of type events.Event');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_events_Event(buffer_arg) {
  return events_events_pb.Event.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_events_PublishRequest(arg) {
  if (!(arg instanceof events_events_pb.PublishRequest)) {
    throw new Error('Expected argument of type events.PublishRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_events_PublishRequest(buffer_arg) {
  return events_events_pb.PublishRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_events_PublishResponse(arg) {
  if (!(arg instanceof events_events_pb.PublishResponse)) {
    throw new Error('Expected argument of type events.PublishResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_events_PublishResponse(buffer_arg) {
  return events_events_pb.PublishResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_events_ReadRequest(arg) {
  if (!(arg instanceof events_events_pb.ReadRequest)) {
    throw new Error('Expected argument of type events.ReadRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_events_ReadRequest(buffer_arg) {
  return events_events_pb.ReadRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_events_ReadResponse(arg) {
  if (!(arg instanceof events_events_pb.ReadResponse)) {
    throw new Error('Expected argument of type events.ReadResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_events_ReadResponse(buffer_arg) {
  return events_events_pb.ReadResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_events_WriteRequest(arg) {
  if (!(arg instanceof events_events_pb.WriteRequest)) {
    throw new Error('Expected argument of type events.WriteRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_events_WriteRequest(buffer_arg) {
  return events_events_pb.WriteRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_events_WriteResponse(arg) {
  if (!(arg instanceof events_events_pb.WriteResponse)) {
    throw new Error('Expected argument of type events.WriteResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_events_WriteResponse(buffer_arg) {
  return events_events_pb.WriteResponse.deserializeBinary(new Uint8Array(buffer_arg));
}


var StreamService = exports.StreamService = {
  publish: {
    path: '/events.Stream/Publish',
    requestStream: false,
    responseStream: false,
    requestType: events_events_pb.PublishRequest,
    responseType: events_events_pb.PublishResponse,
    requestSerialize: serialize_events_PublishRequest,
    requestDeserialize: deserialize_events_PublishRequest,
    responseSerialize: serialize_events_PublishResponse,
    responseDeserialize: deserialize_events_PublishResponse,
  },
  consume: {
    path: '/events.Stream/Consume',
    requestStream: false,
    responseStream: true,
    requestType: events_events_pb.ConsumeRequest,
    responseType: events_events_pb.Event,
    requestSerialize: serialize_events_ConsumeRequest,
    requestDeserialize: deserialize_events_ConsumeRequest,
    responseSerialize: serialize_events_Event,
    responseDeserialize: deserialize_events_Event,
  },
};

exports.StreamClient = grpc.makeGenericClientConstructor(StreamService);
var StoreService = exports.StoreService = {
  read: {
    path: '/events.Store/Read',
    requestStream: false,
    responseStream: false,
    requestType: events_events_pb.ReadRequest,
    responseType: events_events_pb.ReadResponse,
    requestSerialize: serialize_events_ReadRequest,
    requestDeserialize: deserialize_events_ReadRequest,
    responseSerialize: serialize_events_ReadResponse,
    responseDeserialize: deserialize_events_ReadResponse,
  },
  write: {
    path: '/events.Store/Write',
    requestStream: false,
    responseStream: false,
    requestType: events_events_pb.WriteRequest,
    responseType: events_events_pb.WriteResponse,
    requestSerialize: serialize_events_WriteRequest,
    requestDeserialize: deserialize_events_WriteRequest,
    responseSerialize: serialize_events_WriteResponse,
    responseDeserialize: deserialize_events_WriteResponse,
  },
};

exports.StoreClient = grpc.makeGenericClientConstructor(StoreService);
