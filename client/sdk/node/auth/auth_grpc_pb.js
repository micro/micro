// GENERATED CODE -- DO NOT EDIT!

'use strict';
var grpc = require('grpc');
var auth_auth_pb = require('../auth/auth_pb.js');

function serialize_auth_ChangeSecretRequest(arg) {
  if (!(arg instanceof auth_auth_pb.ChangeSecretRequest)) {
    throw new Error('Expected argument of type auth.ChangeSecretRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_auth_ChangeSecretRequest(buffer_arg) {
  return auth_auth_pb.ChangeSecretRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_auth_ChangeSecretResponse(arg) {
  if (!(arg instanceof auth_auth_pb.ChangeSecretResponse)) {
    throw new Error('Expected argument of type auth.ChangeSecretResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_auth_ChangeSecretResponse(buffer_arg) {
  return auth_auth_pb.ChangeSecretResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_auth_CreateRequest(arg) {
  if (!(arg instanceof auth_auth_pb.CreateRequest)) {
    throw new Error('Expected argument of type auth.CreateRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_auth_CreateRequest(buffer_arg) {
  return auth_auth_pb.CreateRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_auth_CreateResponse(arg) {
  if (!(arg instanceof auth_auth_pb.CreateResponse)) {
    throw new Error('Expected argument of type auth.CreateResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_auth_CreateResponse(buffer_arg) {
  return auth_auth_pb.CreateResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_auth_DeleteAccountRequest(arg) {
  if (!(arg instanceof auth_auth_pb.DeleteAccountRequest)) {
    throw new Error('Expected argument of type auth.DeleteAccountRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_auth_DeleteAccountRequest(buffer_arg) {
  return auth_auth_pb.DeleteAccountRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_auth_DeleteAccountResponse(arg) {
  if (!(arg instanceof auth_auth_pb.DeleteAccountResponse)) {
    throw new Error('Expected argument of type auth.DeleteAccountResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_auth_DeleteAccountResponse(buffer_arg) {
  return auth_auth_pb.DeleteAccountResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_auth_DeleteRequest(arg) {
  if (!(arg instanceof auth_auth_pb.DeleteRequest)) {
    throw new Error('Expected argument of type auth.DeleteRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_auth_DeleteRequest(buffer_arg) {
  return auth_auth_pb.DeleteRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_auth_DeleteResponse(arg) {
  if (!(arg instanceof auth_auth_pb.DeleteResponse)) {
    throw new Error('Expected argument of type auth.DeleteResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_auth_DeleteResponse(buffer_arg) {
  return auth_auth_pb.DeleteResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_auth_GenerateRequest(arg) {
  if (!(arg instanceof auth_auth_pb.GenerateRequest)) {
    throw new Error('Expected argument of type auth.GenerateRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_auth_GenerateRequest(buffer_arg) {
  return auth_auth_pb.GenerateRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_auth_GenerateResponse(arg) {
  if (!(arg instanceof auth_auth_pb.GenerateResponse)) {
    throw new Error('Expected argument of type auth.GenerateResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_auth_GenerateResponse(buffer_arg) {
  return auth_auth_pb.GenerateResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_auth_InspectRequest(arg) {
  if (!(arg instanceof auth_auth_pb.InspectRequest)) {
    throw new Error('Expected argument of type auth.InspectRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_auth_InspectRequest(buffer_arg) {
  return auth_auth_pb.InspectRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_auth_InspectResponse(arg) {
  if (!(arg instanceof auth_auth_pb.InspectResponse)) {
    throw new Error('Expected argument of type auth.InspectResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_auth_InspectResponse(buffer_arg) {
  return auth_auth_pb.InspectResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_auth_ListAccountsRequest(arg) {
  if (!(arg instanceof auth_auth_pb.ListAccountsRequest)) {
    throw new Error('Expected argument of type auth.ListAccountsRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_auth_ListAccountsRequest(buffer_arg) {
  return auth_auth_pb.ListAccountsRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_auth_ListAccountsResponse(arg) {
  if (!(arg instanceof auth_auth_pb.ListAccountsResponse)) {
    throw new Error('Expected argument of type auth.ListAccountsResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_auth_ListAccountsResponse(buffer_arg) {
  return auth_auth_pb.ListAccountsResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_auth_ListRequest(arg) {
  if (!(arg instanceof auth_auth_pb.ListRequest)) {
    throw new Error('Expected argument of type auth.ListRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_auth_ListRequest(buffer_arg) {
  return auth_auth_pb.ListRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_auth_ListResponse(arg) {
  if (!(arg instanceof auth_auth_pb.ListResponse)) {
    throw new Error('Expected argument of type auth.ListResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_auth_ListResponse(buffer_arg) {
  return auth_auth_pb.ListResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_auth_TokenRequest(arg) {
  if (!(arg instanceof auth_auth_pb.TokenRequest)) {
    throw new Error('Expected argument of type auth.TokenRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_auth_TokenRequest(buffer_arg) {
  return auth_auth_pb.TokenRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_auth_TokenResponse(arg) {
  if (!(arg instanceof auth_auth_pb.TokenResponse)) {
    throw new Error('Expected argument of type auth.TokenResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_auth_TokenResponse(buffer_arg) {
  return auth_auth_pb.TokenResponse.deserializeBinary(new Uint8Array(buffer_arg));
}


var AuthService = exports.AuthService = {
  generate: {
    path: '/auth.Auth/Generate',
    requestStream: false,
    responseStream: false,
    requestType: auth_auth_pb.GenerateRequest,
    responseType: auth_auth_pb.GenerateResponse,
    requestSerialize: serialize_auth_GenerateRequest,
    requestDeserialize: deserialize_auth_GenerateRequest,
    responseSerialize: serialize_auth_GenerateResponse,
    responseDeserialize: deserialize_auth_GenerateResponse,
  },
  inspect: {
    path: '/auth.Auth/Inspect',
    requestStream: false,
    responseStream: false,
    requestType: auth_auth_pb.InspectRequest,
    responseType: auth_auth_pb.InspectResponse,
    requestSerialize: serialize_auth_InspectRequest,
    requestDeserialize: deserialize_auth_InspectRequest,
    responseSerialize: serialize_auth_InspectResponse,
    responseDeserialize: deserialize_auth_InspectResponse,
  },
  token: {
    path: '/auth.Auth/Token',
    requestStream: false,
    responseStream: false,
    requestType: auth_auth_pb.TokenRequest,
    responseType: auth_auth_pb.TokenResponse,
    requestSerialize: serialize_auth_TokenRequest,
    requestDeserialize: deserialize_auth_TokenRequest,
    responseSerialize: serialize_auth_TokenResponse,
    responseDeserialize: deserialize_auth_TokenResponse,
  },
};

exports.AuthClient = grpc.makeGenericClientConstructor(AuthService);
var AccountsService = exports.AccountsService = {
  list: {
    path: '/auth.Accounts/List',
    requestStream: false,
    responseStream: false,
    requestType: auth_auth_pb.ListAccountsRequest,
    responseType: auth_auth_pb.ListAccountsResponse,
    requestSerialize: serialize_auth_ListAccountsRequest,
    requestDeserialize: deserialize_auth_ListAccountsRequest,
    responseSerialize: serialize_auth_ListAccountsResponse,
    responseDeserialize: deserialize_auth_ListAccountsResponse,
  },
  delete: {
    path: '/auth.Accounts/Delete',
    requestStream: false,
    responseStream: false,
    requestType: auth_auth_pb.DeleteAccountRequest,
    responseType: auth_auth_pb.DeleteAccountResponse,
    requestSerialize: serialize_auth_DeleteAccountRequest,
    requestDeserialize: deserialize_auth_DeleteAccountRequest,
    responseSerialize: serialize_auth_DeleteAccountResponse,
    responseDeserialize: deserialize_auth_DeleteAccountResponse,
  },
  changeSecret: {
    path: '/auth.Accounts/ChangeSecret',
    requestStream: false,
    responseStream: false,
    requestType: auth_auth_pb.ChangeSecretRequest,
    responseType: auth_auth_pb.ChangeSecretResponse,
    requestSerialize: serialize_auth_ChangeSecretRequest,
    requestDeserialize: deserialize_auth_ChangeSecretRequest,
    responseSerialize: serialize_auth_ChangeSecretResponse,
    responseDeserialize: deserialize_auth_ChangeSecretResponse,
  },
};

exports.AccountsClient = grpc.makeGenericClientConstructor(AccountsService);
var RulesService = exports.RulesService = {
  create: {
    path: '/auth.Rules/Create',
    requestStream: false,
    responseStream: false,
    requestType: auth_auth_pb.CreateRequest,
    responseType: auth_auth_pb.CreateResponse,
    requestSerialize: serialize_auth_CreateRequest,
    requestDeserialize: deserialize_auth_CreateRequest,
    responseSerialize: serialize_auth_CreateResponse,
    responseDeserialize: deserialize_auth_CreateResponse,
  },
  delete: {
    path: '/auth.Rules/Delete',
    requestStream: false,
    responseStream: false,
    requestType: auth_auth_pb.DeleteRequest,
    responseType: auth_auth_pb.DeleteResponse,
    requestSerialize: serialize_auth_DeleteRequest,
    requestDeserialize: deserialize_auth_DeleteRequest,
    responseSerialize: serialize_auth_DeleteResponse,
    responseDeserialize: deserialize_auth_DeleteResponse,
  },
  list: {
    path: '/auth.Rules/List',
    requestStream: false,
    responseStream: false,
    requestType: auth_auth_pb.ListRequest,
    responseType: auth_auth_pb.ListResponse,
    requestSerialize: serialize_auth_ListRequest,
    requestDeserialize: deserialize_auth_ListRequest,
    responseSerialize: serialize_auth_ListResponse,
    responseDeserialize: deserialize_auth_ListResponse,
  },
};

exports.RulesClient = grpc.makeGenericClientConstructor(RulesService);
