// GENERATED CODE -- DO NOT EDIT!

'use strict';
var grpc = require('grpc');
var runtime_runtime_pb = require('../runtime/runtime_pb.js');

function serialize_runtime_BuildReadResponse(arg) {
  if (!(arg instanceof runtime_runtime_pb.BuildReadResponse)) {
    throw new Error('Expected argument of type runtime.BuildReadResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_runtime_BuildReadResponse(buffer_arg) {
  return runtime_runtime_pb.BuildReadResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_runtime_CreateRequest(arg) {
  if (!(arg instanceof runtime_runtime_pb.CreateRequest)) {
    throw new Error('Expected argument of type runtime.CreateRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_runtime_CreateRequest(buffer_arg) {
  return runtime_runtime_pb.CreateRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_runtime_CreateResponse(arg) {
  if (!(arg instanceof runtime_runtime_pb.CreateResponse)) {
    throw new Error('Expected argument of type runtime.CreateResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_runtime_CreateResponse(buffer_arg) {
  return runtime_runtime_pb.CreateResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_runtime_DeleteRequest(arg) {
  if (!(arg instanceof runtime_runtime_pb.DeleteRequest)) {
    throw new Error('Expected argument of type runtime.DeleteRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_runtime_DeleteRequest(buffer_arg) {
  return runtime_runtime_pb.DeleteRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_runtime_DeleteResponse(arg) {
  if (!(arg instanceof runtime_runtime_pb.DeleteResponse)) {
    throw new Error('Expected argument of type runtime.DeleteResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_runtime_DeleteResponse(buffer_arg) {
  return runtime_runtime_pb.DeleteResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_runtime_LogRecord(arg) {
  if (!(arg instanceof runtime_runtime_pb.LogRecord)) {
    throw new Error('Expected argument of type runtime.LogRecord');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_runtime_LogRecord(buffer_arg) {
  return runtime_runtime_pb.LogRecord.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_runtime_LogsRequest(arg) {
  if (!(arg instanceof runtime_runtime_pb.LogsRequest)) {
    throw new Error('Expected argument of type runtime.LogsRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_runtime_LogsRequest(buffer_arg) {
  return runtime_runtime_pb.LogsRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_runtime_ReadRequest(arg) {
  if (!(arg instanceof runtime_runtime_pb.ReadRequest)) {
    throw new Error('Expected argument of type runtime.ReadRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_runtime_ReadRequest(buffer_arg) {
  return runtime_runtime_pb.ReadRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_runtime_ReadResponse(arg) {
  if (!(arg instanceof runtime_runtime_pb.ReadResponse)) {
    throw new Error('Expected argument of type runtime.ReadResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_runtime_ReadResponse(buffer_arg) {
  return runtime_runtime_pb.ReadResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_runtime_Service(arg) {
  if (!(arg instanceof runtime_runtime_pb.Service)) {
    throw new Error('Expected argument of type runtime.Service');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_runtime_Service(buffer_arg) {
  return runtime_runtime_pb.Service.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_runtime_UpdateRequest(arg) {
  if (!(arg instanceof runtime_runtime_pb.UpdateRequest)) {
    throw new Error('Expected argument of type runtime.UpdateRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_runtime_UpdateRequest(buffer_arg) {
  return runtime_runtime_pb.UpdateRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_runtime_UpdateResponse(arg) {
  if (!(arg instanceof runtime_runtime_pb.UpdateResponse)) {
    throw new Error('Expected argument of type runtime.UpdateResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_runtime_UpdateResponse(buffer_arg) {
  return runtime_runtime_pb.UpdateResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_runtime_UploadRequest(arg) {
  if (!(arg instanceof runtime_runtime_pb.UploadRequest)) {
    throw new Error('Expected argument of type runtime.UploadRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_runtime_UploadRequest(buffer_arg) {
  return runtime_runtime_pb.UploadRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_runtime_UploadResponse(arg) {
  if (!(arg instanceof runtime_runtime_pb.UploadResponse)) {
    throw new Error('Expected argument of type runtime.UploadResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_runtime_UploadResponse(buffer_arg) {
  return runtime_runtime_pb.UploadResponse.deserializeBinary(new Uint8Array(buffer_arg));
}


var RuntimeService = exports.RuntimeService = {
  create: {
    path: '/runtime.Runtime/Create',
    requestStream: false,
    responseStream: false,
    requestType: runtime_runtime_pb.CreateRequest,
    responseType: runtime_runtime_pb.CreateResponse,
    requestSerialize: serialize_runtime_CreateRequest,
    requestDeserialize: deserialize_runtime_CreateRequest,
    responseSerialize: serialize_runtime_CreateResponse,
    responseDeserialize: deserialize_runtime_CreateResponse,
  },
  read: {
    path: '/runtime.Runtime/Read',
    requestStream: false,
    responseStream: false,
    requestType: runtime_runtime_pb.ReadRequest,
    responseType: runtime_runtime_pb.ReadResponse,
    requestSerialize: serialize_runtime_ReadRequest,
    requestDeserialize: deserialize_runtime_ReadRequest,
    responseSerialize: serialize_runtime_ReadResponse,
    responseDeserialize: deserialize_runtime_ReadResponse,
  },
  delete: {
    path: '/runtime.Runtime/Delete',
    requestStream: false,
    responseStream: false,
    requestType: runtime_runtime_pb.DeleteRequest,
    responseType: runtime_runtime_pb.DeleteResponse,
    requestSerialize: serialize_runtime_DeleteRequest,
    requestDeserialize: deserialize_runtime_DeleteRequest,
    responseSerialize: serialize_runtime_DeleteResponse,
    responseDeserialize: deserialize_runtime_DeleteResponse,
  },
  update: {
    path: '/runtime.Runtime/Update',
    requestStream: false,
    responseStream: false,
    requestType: runtime_runtime_pb.UpdateRequest,
    responseType: runtime_runtime_pb.UpdateResponse,
    requestSerialize: serialize_runtime_UpdateRequest,
    requestDeserialize: deserialize_runtime_UpdateRequest,
    responseSerialize: serialize_runtime_UpdateResponse,
    responseDeserialize: deserialize_runtime_UpdateResponse,
  },
  logs: {
    path: '/runtime.Runtime/Logs',
    requestStream: false,
    responseStream: true,
    requestType: runtime_runtime_pb.LogsRequest,
    responseType: runtime_runtime_pb.LogRecord,
    requestSerialize: serialize_runtime_LogsRequest,
    requestDeserialize: deserialize_runtime_LogsRequest,
    responseSerialize: serialize_runtime_LogRecord,
    responseDeserialize: deserialize_runtime_LogRecord,
  },
};

exports.RuntimeClient = grpc.makeGenericClientConstructor(RuntimeService);
// Source service is used by the CLI to upload source to the service. The service will return
// a unique ID representing the location of that source. This ID can then be used as a source
// for the service when doing Runtime.Create. The server will handle cleanup of uploaded source.
var SourceService = exports.SourceService = {
  upload: {
    path: '/runtime.Source/Upload',
    requestStream: true,
    responseStream: false,
    requestType: runtime_runtime_pb.UploadRequest,
    responseType: runtime_runtime_pb.UploadResponse,
    requestSerialize: serialize_runtime_UploadRequest,
    requestDeserialize: deserialize_runtime_UploadRequest,
    responseSerialize: serialize_runtime_UploadResponse,
    responseDeserialize: deserialize_runtime_UploadResponse,
  },
};

exports.SourceClient = grpc.makeGenericClientConstructor(SourceService);
// Build service is used by containers to download prebuilt binaries. The client will pass the 
// service (name and version are required attributed) and the server will then stream the latest
// binary to the client.
var BuildService = exports.BuildService = {
  read: {
    path: '/runtime.Build/Read',
    requestStream: false,
    responseStream: true,
    requestType: runtime_runtime_pb.Service,
    responseType: runtime_runtime_pb.BuildReadResponse,
    requestSerialize: serialize_runtime_Service,
    requestDeserialize: deserialize_runtime_Service,
    responseSerialize: serialize_runtime_BuildReadResponse,
    responseDeserialize: deserialize_runtime_BuildReadResponse,
  },
};

exports.BuildClient = grpc.makeGenericClientConstructor(BuildService);
