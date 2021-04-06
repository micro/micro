// GENERATED CODE -- DO NOT EDIT!

'use strict';
var grpc = require('grpc');
var router_router_pb = require('../router/router_pb.js');

function serialize_router_CreateResponse(arg) {
  if (!(arg instanceof router_router_pb.CreateResponse)) {
    throw new Error('Expected argument of type router.CreateResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_router_CreateResponse(buffer_arg) {
  return router_router_pb.CreateResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_router_DeleteResponse(arg) {
  if (!(arg instanceof router_router_pb.DeleteResponse)) {
    throw new Error('Expected argument of type router.DeleteResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_router_DeleteResponse(buffer_arg) {
  return router_router_pb.DeleteResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_router_Event(arg) {
  if (!(arg instanceof router_router_pb.Event)) {
    throw new Error('Expected argument of type router.Event');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_router_Event(buffer_arg) {
  return router_router_pb.Event.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_router_LookupRequest(arg) {
  if (!(arg instanceof router_router_pb.LookupRequest)) {
    throw new Error('Expected argument of type router.LookupRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_router_LookupRequest(buffer_arg) {
  return router_router_pb.LookupRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_router_LookupResponse(arg) {
  if (!(arg instanceof router_router_pb.LookupResponse)) {
    throw new Error('Expected argument of type router.LookupResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_router_LookupResponse(buffer_arg) {
  return router_router_pb.LookupResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_router_ReadRequest(arg) {
  if (!(arg instanceof router_router_pb.ReadRequest)) {
    throw new Error('Expected argument of type router.ReadRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_router_ReadRequest(buffer_arg) {
  return router_router_pb.ReadRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_router_ReadResponse(arg) {
  if (!(arg instanceof router_router_pb.ReadResponse)) {
    throw new Error('Expected argument of type router.ReadResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_router_ReadResponse(buffer_arg) {
  return router_router_pb.ReadResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_router_Route(arg) {
  if (!(arg instanceof router_router_pb.Route)) {
    throw new Error('Expected argument of type router.Route');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_router_Route(buffer_arg) {
  return router_router_pb.Route.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_router_UpdateResponse(arg) {
  if (!(arg instanceof router_router_pb.UpdateResponse)) {
    throw new Error('Expected argument of type router.UpdateResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_router_UpdateResponse(buffer_arg) {
  return router_router_pb.UpdateResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_router_WatchRequest(arg) {
  if (!(arg instanceof router_router_pb.WatchRequest)) {
    throw new Error('Expected argument of type router.WatchRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_router_WatchRequest(buffer_arg) {
  return router_router_pb.WatchRequest.deserializeBinary(new Uint8Array(buffer_arg));
}


// Router service is used by the proxy to lookup routes
var RouterService = exports.RouterService = {
  lookup: {
    path: '/router.Router/Lookup',
    requestStream: false,
    responseStream: false,
    requestType: router_router_pb.LookupRequest,
    responseType: router_router_pb.LookupResponse,
    requestSerialize: serialize_router_LookupRequest,
    requestDeserialize: deserialize_router_LookupRequest,
    responseSerialize: serialize_router_LookupResponse,
    responseDeserialize: deserialize_router_LookupResponse,
  },
  watch: {
    path: '/router.Router/Watch',
    requestStream: false,
    responseStream: true,
    requestType: router_router_pb.WatchRequest,
    responseType: router_router_pb.Event,
    requestSerialize: serialize_router_WatchRequest,
    requestDeserialize: deserialize_router_WatchRequest,
    responseSerialize: serialize_router_Event,
    responseDeserialize: deserialize_router_Event,
  },
};

exports.RouterClient = grpc.makeGenericClientConstructor(RouterService);
var TableService = exports.TableService = {
  create: {
    path: '/router.Table/Create',
    requestStream: false,
    responseStream: false,
    requestType: router_router_pb.Route,
    responseType: router_router_pb.CreateResponse,
    requestSerialize: serialize_router_Route,
    requestDeserialize: deserialize_router_Route,
    responseSerialize: serialize_router_CreateResponse,
    responseDeserialize: deserialize_router_CreateResponse,
  },
  delete: {
    path: '/router.Table/Delete',
    requestStream: false,
    responseStream: false,
    requestType: router_router_pb.Route,
    responseType: router_router_pb.DeleteResponse,
    requestSerialize: serialize_router_Route,
    requestDeserialize: deserialize_router_Route,
    responseSerialize: serialize_router_DeleteResponse,
    responseDeserialize: deserialize_router_DeleteResponse,
  },
  update: {
    path: '/router.Table/Update',
    requestStream: false,
    responseStream: false,
    requestType: router_router_pb.Route,
    responseType: router_router_pb.UpdateResponse,
    requestSerialize: serialize_router_Route,
    requestDeserialize: deserialize_router_Route,
    responseSerialize: serialize_router_UpdateResponse,
    responseDeserialize: deserialize_router_UpdateResponse,
  },
  read: {
    path: '/router.Table/Read',
    requestStream: false,
    responseStream: false,
    requestType: router_router_pb.ReadRequest,
    responseType: router_router_pb.ReadResponse,
    requestSerialize: serialize_router_ReadRequest,
    requestDeserialize: deserialize_router_ReadRequest,
    responseSerialize: serialize_router_ReadResponse,
    responseDeserialize: deserialize_router_ReadResponse,
  },
};

exports.TableClient = grpc.makeGenericClientConstructor(TableService);
