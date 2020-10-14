// GENERATED CODE -- DO NOT EDIT!

'use strict';
var grpc = require('grpc');
var transport_transport_pb = require('../transport/transport_pb.js');

function serialize_transport_Message(arg) {
  if (!(arg instanceof transport_transport_pb.Message)) {
    throw new Error('Expected argument of type transport.Message');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_transport_Message(buffer_arg) {
  return transport_transport_pb.Message.deserializeBinary(new Uint8Array(buffer_arg));
}


var TransportService = exports.TransportService = {
  stream: {
    path: '/transport.Transport/Stream',
    requestStream: true,
    responseStream: true,
    requestType: transport_transport_pb.Message,
    responseType: transport_transport_pb.Message,
    requestSerialize: serialize_transport_Message,
    requestDeserialize: deserialize_transport_Message,
    responseSerialize: serialize_transport_Message,
    responseDeserialize: deserialize_transport_Message,
  },
};

exports.TransportClient = grpc.makeGenericClientConstructor(TransportService);
