// GENERATED CODE -- DO NOT EDIT!

'use strict';
var grpc = require('grpc');
var alert_alert_pb = require('../alert/alert_pb.js');

function serialize_alert_ReportEventRequest(arg) {
  if (!(arg instanceof alert_alert_pb.ReportEventRequest)) {
    throw new Error('Expected argument of type alert.ReportEventRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_alert_ReportEventRequest(buffer_arg) {
  return alert_alert_pb.ReportEventRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_alert_ReportEventResponse(arg) {
  if (!(arg instanceof alert_alert_pb.ReportEventResponse)) {
    throw new Error('Expected argument of type alert.ReportEventResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_alert_ReportEventResponse(buffer_arg) {
  return alert_alert_pb.ReportEventResponse.deserializeBinary(new Uint8Array(buffer_arg));
}


var AlertService = exports.AlertService = {
  // ReportEvent does event ingestions.
reportEvent: {
    path: '/alert.Alert/ReportEvent',
    requestStream: false,
    responseStream: false,
    requestType: alert_alert_pb.ReportEventRequest,
    responseType: alert_alert_pb.ReportEventResponse,
    requestSerialize: serialize_alert_ReportEventRequest,
    requestDeserialize: deserialize_alert_ReportEventRequest,
    responseSerialize: serialize_alert_ReportEventResponse,
    responseDeserialize: deserialize_alert_ReportEventResponse,
  },
};

exports.AlertClient = grpc.makeGenericClientConstructor(AlertService);
