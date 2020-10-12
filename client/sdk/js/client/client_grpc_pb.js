// GENERATED CODE -- DO NOT EDIT!

'use strict';
var grpc = require('grpc');
var client_client_pb = require('../client/client_pb.js');

function serialize_client_Message(arg) {
  if (!(arg instanceof client_client_pb.Message)) {
    throw new Error('Expected argument of type client.Message');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_client_Message(buffer_arg) {
  return client_client_pb.Message.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_client_Request(arg) {
  if (!(arg instanceof client_client_pb.Request)) {
    throw new Error('Expected argument of type client.Request');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_client_Request(buffer_arg) {
  return client_client_pb.Request.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_client_Response(arg) {
  if (!(arg instanceof client_client_pb.Response)) {
    throw new Error('Expected argument of type client.Response');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_client_Response(buffer_arg) {
  return client_client_pb.Response.deserializeBinary(new Uint8Array(buffer_arg));
}


// Client is the micro client interface
var ClientService = exports.ClientService = {
  // Call allows a single request to be made
call: {
    path: '/client.Client/Call',
    requestStream: false,
    responseStream: false,
    requestType: client_client_pb.Request,
    responseType: client_client_pb.Response,
    requestSerialize: serialize_client_Request,
    requestDeserialize: deserialize_client_Request,
    responseSerialize: serialize_client_Response,
    responseDeserialize: deserialize_client_Response,
  },
  // Stream is a bidirectional stream
stream: {
    path: '/client.Client/Stream',
    requestStream: true,
    responseStream: true,
    requestType: client_client_pb.Request,
    responseType: client_client_pb.Response,
    requestSerialize: serialize_client_Request,
    requestDeserialize: deserialize_client_Request,
    responseSerialize: serialize_client_Response,
    responseDeserialize: deserialize_client_Response,
  },
  // Publish publishes a message and returns an empty Message
publish: {
    path: '/client.Client/Publish',
    requestStream: false,
    responseStream: false,
    requestType: client_client_pb.Message,
    responseType: client_client_pb.Message,
    requestSerialize: serialize_client_Message,
    requestDeserialize: deserialize_client_Message,
    responseSerialize: serialize_client_Message,
    responseDeserialize: deserialize_client_Message,
  },
};

exports.ClientClient = grpc.makeGenericClientConstructor(ClientService);
