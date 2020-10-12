// GENERATED CODE -- DO NOT EDIT!

'use strict';
var grpc = require('grpc');
var broker_broker_pb = require('../broker/broker_pb.js');

function serialize_broker_Empty(arg) {
  if (!(arg instanceof broker_broker_pb.Empty)) {
    throw new Error('Expected argument of type broker.Empty');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_broker_Empty(buffer_arg) {
  return broker_broker_pb.Empty.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_broker_Message(arg) {
  if (!(arg instanceof broker_broker_pb.Message)) {
    throw new Error('Expected argument of type broker.Message');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_broker_Message(buffer_arg) {
  return broker_broker_pb.Message.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_broker_PublishRequest(arg) {
  if (!(arg instanceof broker_broker_pb.PublishRequest)) {
    throw new Error('Expected argument of type broker.PublishRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_broker_PublishRequest(buffer_arg) {
  return broker_broker_pb.PublishRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_broker_SubscribeRequest(arg) {
  if (!(arg instanceof broker_broker_pb.SubscribeRequest)) {
    throw new Error('Expected argument of type broker.SubscribeRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_broker_SubscribeRequest(buffer_arg) {
  return broker_broker_pb.SubscribeRequest.deserializeBinary(new Uint8Array(buffer_arg));
}


var BrokerService = exports.BrokerService = {
  publish: {
    path: '/broker.Broker/Publish',
    requestStream: false,
    responseStream: false,
    requestType: broker_broker_pb.PublishRequest,
    responseType: broker_broker_pb.Empty,
    requestSerialize: serialize_broker_PublishRequest,
    requestDeserialize: deserialize_broker_PublishRequest,
    responseSerialize: serialize_broker_Empty,
    responseDeserialize: deserialize_broker_Empty,
  },
  subscribe: {
    path: '/broker.Broker/Subscribe',
    requestStream: false,
    responseStream: true,
    requestType: broker_broker_pb.SubscribeRequest,
    responseType: broker_broker_pb.Message,
    requestSerialize: serialize_broker_SubscribeRequest,
    requestDeserialize: deserialize_broker_SubscribeRequest,
    responseSerialize: serialize_broker_Message,
    responseDeserialize: deserialize_broker_Message,
  },
};

exports.BrokerClient = grpc.makeGenericClientConstructor(BrokerService);
