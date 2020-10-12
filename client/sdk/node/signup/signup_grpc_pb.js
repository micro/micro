// GENERATED CODE -- DO NOT EDIT!

'use strict';
var grpc = require('grpc');
var signup_signup_pb = require('../signup/signup_pb.js');

function serialize_go_micro_service_signup_CompleteSignupRequest(arg) {
  if (!(arg instanceof signup_signup_pb.CompleteSignupRequest)) {
    throw new Error('Expected argument of type go.micro.service.signup.CompleteSignupRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_go_micro_service_signup_CompleteSignupRequest(buffer_arg) {
  return signup_signup_pb.CompleteSignupRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_go_micro_service_signup_CompleteSignupResponse(arg) {
  if (!(arg instanceof signup_signup_pb.CompleteSignupResponse)) {
    throw new Error('Expected argument of type go.micro.service.signup.CompleteSignupResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_go_micro_service_signup_CompleteSignupResponse(buffer_arg) {
  return signup_signup_pb.CompleteSignupResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_go_micro_service_signup_HasPaymentMethodRequest(arg) {
  if (!(arg instanceof signup_signup_pb.HasPaymentMethodRequest)) {
    throw new Error('Expected argument of type go.micro.service.signup.HasPaymentMethodRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_go_micro_service_signup_HasPaymentMethodRequest(buffer_arg) {
  return signup_signup_pb.HasPaymentMethodRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_go_micro_service_signup_HasPaymentMethodResponse(arg) {
  if (!(arg instanceof signup_signup_pb.HasPaymentMethodResponse)) {
    throw new Error('Expected argument of type go.micro.service.signup.HasPaymentMethodResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_go_micro_service_signup_HasPaymentMethodResponse(buffer_arg) {
  return signup_signup_pb.HasPaymentMethodResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_go_micro_service_signup_RecoverRequest(arg) {
  if (!(arg instanceof signup_signup_pb.RecoverRequest)) {
    throw new Error('Expected argument of type go.micro.service.signup.RecoverRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_go_micro_service_signup_RecoverRequest(buffer_arg) {
  return signup_signup_pb.RecoverRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_go_micro_service_signup_RecoverResponse(arg) {
  if (!(arg instanceof signup_signup_pb.RecoverResponse)) {
    throw new Error('Expected argument of type go.micro.service.signup.RecoverResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_go_micro_service_signup_RecoverResponse(buffer_arg) {
  return signup_signup_pb.RecoverResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_go_micro_service_signup_SendVerificationEmailRequest(arg) {
  if (!(arg instanceof signup_signup_pb.SendVerificationEmailRequest)) {
    throw new Error('Expected argument of type go.micro.service.signup.SendVerificationEmailRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_go_micro_service_signup_SendVerificationEmailRequest(buffer_arg) {
  return signup_signup_pb.SendVerificationEmailRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_go_micro_service_signup_SendVerificationEmailResponse(arg) {
  if (!(arg instanceof signup_signup_pb.SendVerificationEmailResponse)) {
    throw new Error('Expected argument of type go.micro.service.signup.SendVerificationEmailResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_go_micro_service_signup_SendVerificationEmailResponse(buffer_arg) {
  return signup_signup_pb.SendVerificationEmailResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_go_micro_service_signup_SetPaymentMethodRequest(arg) {
  if (!(arg instanceof signup_signup_pb.SetPaymentMethodRequest)) {
    throw new Error('Expected argument of type go.micro.service.signup.SetPaymentMethodRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_go_micro_service_signup_SetPaymentMethodRequest(buffer_arg) {
  return signup_signup_pb.SetPaymentMethodRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_go_micro_service_signup_SetPaymentMethodResponse(arg) {
  if (!(arg instanceof signup_signup_pb.SetPaymentMethodResponse)) {
    throw new Error('Expected argument of type go.micro.service.signup.SetPaymentMethodResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_go_micro_service_signup_SetPaymentMethodResponse(buffer_arg) {
  return signup_signup_pb.SetPaymentMethodResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_go_micro_service_signup_VerifyRequest(arg) {
  if (!(arg instanceof signup_signup_pb.VerifyRequest)) {
    throw new Error('Expected argument of type go.micro.service.signup.VerifyRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_go_micro_service_signup_VerifyRequest(buffer_arg) {
  return signup_signup_pb.VerifyRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_go_micro_service_signup_VerifyResponse(arg) {
  if (!(arg instanceof signup_signup_pb.VerifyResponse)) {
    throw new Error('Expected argument of type go.micro.service.signup.VerifyResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_go_micro_service_signup_VerifyResponse(buffer_arg) {
  return signup_signup_pb.VerifyResponse.deserializeBinary(new Uint8Array(buffer_arg));
}


var SignupService = exports.SignupService = {
  // Sends the verification email to the user
sendVerificationEmail: {
    path: '/go.micro.service.signup.Signup/SendVerificationEmail',
    requestStream: false,
    responseStream: false,
    requestType: signup_signup_pb.SendVerificationEmailRequest,
    responseType: signup_signup_pb.SendVerificationEmailResponse,
    requestSerialize: serialize_go_micro_service_signup_SendVerificationEmailRequest,
    requestDeserialize: deserialize_go_micro_service_signup_SendVerificationEmailRequest,
    responseSerialize: serialize_go_micro_service_signup_SendVerificationEmailResponse,
    responseDeserialize: deserialize_go_micro_service_signup_SendVerificationEmailResponse,
  },
  // Verify kicks off the process of verification
verify: {
    path: '/go.micro.service.signup.Signup/Verify',
    requestStream: false,
    responseStream: false,
    requestType: signup_signup_pb.VerifyRequest,
    responseType: signup_signup_pb.VerifyResponse,
    requestSerialize: serialize_go_micro_service_signup_VerifyRequest,
    requestDeserialize: deserialize_go_micro_service_signup_VerifyRequest,
    responseSerialize: serialize_go_micro_service_signup_VerifyResponse,
    responseDeserialize: deserialize_go_micro_service_signup_VerifyResponse,
  },
  setPaymentMethod: {
    path: '/go.micro.service.signup.Signup/SetPaymentMethod',
    requestStream: false,
    responseStream: false,
    requestType: signup_signup_pb.SetPaymentMethodRequest,
    responseType: signup_signup_pb.SetPaymentMethodResponse,
    requestSerialize: serialize_go_micro_service_signup_SetPaymentMethodRequest,
    requestDeserialize: deserialize_go_micro_service_signup_SetPaymentMethodRequest,
    responseSerialize: serialize_go_micro_service_signup_SetPaymentMethodResponse,
    responseDeserialize: deserialize_go_micro_service_signup_SetPaymentMethodResponse,
  },
  hasPaymentMethod: {
    path: '/go.micro.service.signup.Signup/HasPaymentMethod',
    requestStream: false,
    responseStream: false,
    requestType: signup_signup_pb.HasPaymentMethodRequest,
    responseType: signup_signup_pb.HasPaymentMethodResponse,
    requestSerialize: serialize_go_micro_service_signup_HasPaymentMethodRequest,
    requestDeserialize: deserialize_go_micro_service_signup_HasPaymentMethodRequest,
    responseSerialize: serialize_go_micro_service_signup_HasPaymentMethodResponse,
    responseDeserialize: deserialize_go_micro_service_signup_HasPaymentMethodResponse,
  },
  // Creates a subscription and an account
completeSignup: {
    path: '/go.micro.service.signup.Signup/CompleteSignup',
    requestStream: false,
    responseStream: false,
    requestType: signup_signup_pb.CompleteSignupRequest,
    responseType: signup_signup_pb.CompleteSignupResponse,
    requestSerialize: serialize_go_micro_service_signup_CompleteSignupRequest,
    requestDeserialize: deserialize_go_micro_service_signup_CompleteSignupRequest,
    responseSerialize: serialize_go_micro_service_signup_CompleteSignupResponse,
    responseDeserialize: deserialize_go_micro_service_signup_CompleteSignupResponse,
  },
  recover: {
    path: '/go.micro.service.signup.Signup/Recover',
    requestStream: false,
    responseStream: false,
    requestType: signup_signup_pb.RecoverRequest,
    responseType: signup_signup_pb.RecoverResponse,
    requestSerialize: serialize_go_micro_service_signup_RecoverRequest,
    requestDeserialize: deserialize_go_micro_service_signup_RecoverRequest,
    responseSerialize: serialize_go_micro_service_signup_RecoverResponse,
    responseDeserialize: deserialize_go_micro_service_signup_RecoverResponse,
  },
};

exports.SignupClient = grpc.makeGenericClientConstructor(SignupService);
