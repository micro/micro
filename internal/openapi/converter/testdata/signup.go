package testdata

const Signup = `{
  "components": {
    "requestBodies": {
      "SignupCompleteSignupRequest": {
        "content": {
          "application/json": {
            "schema": {
              "$ref": "#/components/schemas/CompleteSignupRequest"
            }
          }
        },
        "description": "SignupCompleteSignupRequest"
      },
      "SignupHasPaymentMethodRequest": {
        "content": {
          "application/json": {
            "schema": {
              "$ref": "#/components/schemas/HasPaymentMethodRequest"
            }
          }
        },
        "description": "SignupHasPaymentMethodRequest"
      },
      "SignupRecoverRequest": {
        "content": {
          "application/json": {
            "schema": {
              "$ref": "#/components/schemas/RecoverRequest"
            }
          }
        },
        "description": "SignupRecoverRequest"
      },
      "SignupSendVerificationEmailRequest": {
        "content": {
          "application/json": {
            "schema": {
              "$ref": "#/components/schemas/SendVerificationEmailRequest"
            }
          }
        },
        "description": "SignupSendVerificationEmailRequest"
      },
      "SignupSetPaymentMethodRequest": {
        "content": {
          "application/json": {
            "schema": {
              "$ref": "#/components/schemas/SetPaymentMethodRequest"
            }
          }
        },
        "description": "SignupSetPaymentMethodRequest"
      },
      "SignupVerifyRequest": {
        "content": {
          "application/json": {
            "schema": {
              "$ref": "#/components/schemas/VerifyRequest"
            }
          }
        },
        "description": "SignupVerifyRequest"
      }
    },
    "responses": {
      "MicroAPIError": {
        "content": {
          "application/json": {
            "schema": {
              "properties": {
                "Code": {
                  "description": "Error code",
                  "example": 500,
                  "type": "number"
                },
                "Detail": {
                  "description": "Error detail",
                  "example": "service not found",
                  "type": "string"
                },
                "Id": {
                  "description": "Error ID",
                  "type": "string"
                },
                "Status": {
                  "description": "Error status message",
                  "example": "Internal Server Error",
                  "type": "string"
                }
              },
              "title": "MicroAPIError",
              "type": "object"
            }
          }
        },
        "description": "Error from the Micro API"
      },
      "SignupCompleteSignupResponse": {
        "content": {
          "application/json": {
            "schema": {
              "$ref": "#/components/schemas/CompleteSignupResponse"
            }
          }
        },
        "description": "SignupCompleteSignupResponse"
      },
      "SignupHasPaymentMethodResponse": {
        "content": {
          "application/json": {
            "schema": {
              "$ref": "#/components/schemas/HasPaymentMethodResponse"
            }
          }
        },
        "description": "SignupHasPaymentMethodResponse"
      },
      "SignupRecoverResponse": {
        "content": {
          "application/json": {
            "schema": {
              "$ref": "#/components/schemas/RecoverResponse"
            }
          }
        },
        "description": "SignupRecoverResponse"
      },
      "SignupSendVerificationEmailResponse": {
        "content": {
          "application/json": {
            "schema": {
              "$ref": "#/components/schemas/SendVerificationEmailResponse"
            }
          }
        },
        "description": "SignupSendVerificationEmailResponse"
      },
      "SignupSetPaymentMethodResponse": {
        "content": {
          "application/json": {
            "schema": {
              "$ref": "#/components/schemas/SetPaymentMethodResponse"
            }
          }
        },
        "description": "SignupSetPaymentMethodResponse"
      },
      "SignupVerifyResponse": {
        "content": {
          "application/json": {
            "schema": {
              "$ref": "#/components/schemas/VerifyResponse"
            }
          }
        },
        "description": "SignupVerifyResponse"
      }
    },
    "schemas": {
      "AuthToken": {
        "description": "lifted from https://github.com/micro/go-micro/blob/master/auth/service/proto/auth.proto",
        "properties": {
          "accessToken": {
            "type": "string"
          },
          "created": {
            "format": "int64",
            "type": "number"
          },
          "expiry": {
            "format": "int64",
            "type": "number"
          },
          "refreshToken": {
            "type": "string"
          }
        },
        "title": "AuthToken",
        "type": "object"
      },
      "CompleteSignupRequest": {
        "properties": {
          "email": {
            "type": "string"
          },
          "namespace": {
            "description": "Which namespace to sign up to based on previous invite",
            "type": "string"
          },
          "paymentMethodID": {
            "description": "This payment method ID is the one we got back from Stripe on the frontend (ie. 'm3o.com/subscribe.html')\n deprecated: signup service now knows the payment method due to the\n SetPaymentMethod call issued by the frontend.",
            "type": "string"
          },
          "secret": {
            "description": "The secret/password to use for the account",
            "type": "string"
          },
          "token": {
            "description": "The token has to be passed here too for identification purposes.",
            "type": "string"
          }
        },
        "title": "CompleteSignupRequest",
        "type": "object"
      },
      "CompleteSignupResponse": {
        "properties": {
          "authToken": {
            "properties": {
              "accessToken": {
                "type": "string"
              },
              "created": {
                "format": "int64",
                "type": "number"
              },
              "expiry": {
                "format": "int64",
                "type": "number"
              },
              "refreshToken": {
                "type": "string"
              }
            },
            "type": "object"
          },
          "namespace": {
            "type": "string"
          }
        },
        "title": "CompleteSignupResponse",
        "type": "object"
      },
      "HasPaymentMethodRequest": {
        "properties": {
          "token": {
            "description": "We can't read by email because that would be too easy to guess.\n The token is already used for identification purposes during the signup\n so we will use that too to pull for the payment method.",
            "type": "string"
          }
        },
        "title": "HasPaymentMethodRequest",
        "type": "object"
      },
      "HasPaymentMethodResponse": {
        "properties": {
          "has": {
            "type": "boolean"
          }
        },
        "title": "HasPaymentMethodResponse",
        "type": "object"
      },
      "RecoverRequest": {
        "properties": {
          "email": {
            "type": "string"
          }
        },
        "title": "RecoverRequest",
        "type": "object"
      },
      "RecoverResponse": {
        "title": "RecoverResponse",
        "type": "object"
      },
      "SendVerificationEmailRequest": {
        "properties": {
          "email": {
            "type": "string"
          }
        },
        "title": "SendVerificationEmailRequest",
        "type": "object"
      },
      "SendVerificationEmailResponse": {
        "title": "SendVerificationEmailResponse",
        "type": "object"
      },
      "SetPaymentMethodRequest": {
        "properties": {
          "email": {
            "type": "string"
          },
          "paymentMethod": {
            "type": "string"
          }
        },
        "title": "SetPaymentMethodRequest",
        "type": "object"
      },
      "SetPaymentMethodResponse": {
        "title": "SetPaymentMethodResponse",
        "type": "object"
      },
      "VerifyRequest": {
        "properties": {
          "email": {
            "type": "string"
          },
          "token": {
            "description": "Email token that was received in an email.",
            "type": "string"
          }
        },
        "title": "VerifyRequest",
        "type": "object"
      },
      "VerifyResponse": {
        "properties": {
          "authToken": {
            "description": "Auth token to be saved into '~/.micro'\n For users who are already registered and paid,\n the flow stops here.\n For users who are yet to be registered\n the token will be acquired in the 'FinishSignup' step.",
            "properties": {
              "accessToken": {
                "type": "string"
              },
              "created": {
                "format": "int64",
                "type": "number"
              },
              "expiry": {
                "format": "int64",
                "type": "number"
              },
              "refreshToken": {
                "type": "string"
              }
            },
            "type": "object"
          },
          "customerID": {
            "description": "Payment provider custommer id that can be used to\n acquire a payment method, see 'micro login' flow for more.\n @todo this is likely not needed",
            "type": "string"
          },
          "message": {
            "description": "Message to display to the user",
            "type": "string"
          },
          "namespace": {
            "description": "Namespace to use\n @todod deprecated since we no longer support OTP logins",
            "type": "string"
          },
          "namespaces": {
            "description": "Namespaces one has access to based on previous invites\n Currently only 1 is supported",
            "items": {
              "type": "string"
            },
            "type": "array"
          },
          "paymentRequired": {
            "description": "Whether payment is required or not",
            "type": "boolean"
          }
        },
        "title": "VerifyResponse",
        "type": "object"
      }
    },
    "securitySchemes": {
      "MicroAPIToken": {
        "bearerFormat": "JWT",
        "description": "Micro API token",
        "scheme": "bearer",
        "type": "http"
      }
    }
  },
  "info": {
    "description": "Generated by protoc-gen-openapi",
    "title": "Go.Micro.Generator.Test",
    "version": "1",
    "x-logo": {
      "altText": "Micro logo",
      "backgroundColor": "#FFFFFF",
      "url": "https://micro.mu/images/brand.png"
    }
  },
  "openapi": "3.0.0",
  "paths": {
    "/test/Signup/CompleteSignup": {
      "parameters": [
        {
          "in": "header",
          "name": "Micro-Namespace",
          "required": true,
          "schema": {
            "type": "string"
          }
        }
      ],
      "post": {
        "requestBody": {
          "$ref": "#/components/requestBodies/SignupCompleteSignupRequest"
        },
        "responses": {
          "200": {
            "$ref": "#/components/responses/SignupCompleteSignupResponse"
          },
          "default": {
            "$ref": "#/components/responses/MicroAPIError"
          }
        },
        "security": [
          {
            "MicroAPIToken": []
          }
        ],
        "summary": "Signup.CompleteSignup(CompleteSignupRequest)"
      }
    },
    "/test/Signup/HasPaymentMethod": {
      "parameters": [
        {
          "in": "header",
          "name": "Micro-Namespace",
          "required": true,
          "schema": {
            "type": "string"
          }
        }
      ],
      "post": {
        "requestBody": {
          "$ref": "#/components/requestBodies/SignupHasPaymentMethodRequest"
        },
        "responses": {
          "200": {
            "$ref": "#/components/responses/SignupHasPaymentMethodResponse"
          },
          "default": {
            "$ref": "#/components/responses/MicroAPIError"
          }
        },
        "security": [
          {
            "MicroAPIToken": []
          }
        ],
        "summary": "Signup.HasPaymentMethod(HasPaymentMethodRequest)"
      }
    },
    "/test/Signup/Recover": {
      "parameters": [
        {
          "in": "header",
          "name": "Micro-Namespace",
          "required": true,
          "schema": {
            "type": "string"
          }
        }
      ],
      "post": {
        "requestBody": {
          "$ref": "#/components/requestBodies/SignupRecoverRequest"
        },
        "responses": {
          "200": {
            "$ref": "#/components/responses/SignupRecoverResponse"
          },
          "default": {
            "$ref": "#/components/responses/MicroAPIError"
          }
        },
        "security": [
          {
            "MicroAPIToken": []
          }
        ],
        "summary": "Signup.Recover(RecoverRequest)"
      }
    },
    "/test/Signup/SendVerificationEmail": {
      "parameters": [
        {
          "in": "header",
          "name": "Micro-Namespace",
          "required": true,
          "schema": {
            "type": "string"
          }
        }
      ],
      "post": {
        "requestBody": {
          "$ref": "#/components/requestBodies/SignupSendVerificationEmailRequest"
        },
        "responses": {
          "200": {
            "$ref": "#/components/responses/SignupSendVerificationEmailResponse"
          },
          "default": {
            "$ref": "#/components/responses/MicroAPIError"
          }
        },
        "security": [
          {
            "MicroAPIToken": []
          }
        ],
        "summary": "Signup.SendVerificationEmail(SendVerificationEmailRequest)"
      }
    },
    "/test/Signup/SetPaymentMethod": {
      "parameters": [
        {
          "in": "header",
          "name": "Micro-Namespace",
          "required": true,
          "schema": {
            "type": "string"
          }
        }
      ],
      "post": {
        "requestBody": {
          "$ref": "#/components/requestBodies/SignupSetPaymentMethodRequest"
        },
        "responses": {
          "200": {
            "$ref": "#/components/responses/SignupSetPaymentMethodResponse"
          },
          "default": {
            "$ref": "#/components/responses/MicroAPIError"
          }
        },
        "security": [
          {
            "MicroAPIToken": []
          }
        ],
        "summary": "Signup.SetPaymentMethod(SetPaymentMethodRequest)"
      }
    },
    "/test/Signup/Verify": {
      "parameters": [
        {
          "in": "header",
          "name": "Micro-Namespace",
          "required": true,
          "schema": {
            "type": "string"
          }
        }
      ],
      "post": {
        "requestBody": {
          "$ref": "#/components/requestBodies/SignupVerifyRequest"
        },
        "responses": {
          "200": {
            "$ref": "#/components/responses/SignupVerifyResponse"
          },
          "default": {
            "$ref": "#/components/responses/MicroAPIError"
          }
        },
        "security": [
          {
            "MicroAPIToken": []
          }
        ],
        "summary": "Signup.Verify(VerifyRequest)"
      }
    }
  },
  "servers": [
    {
      "url": "https://api.m3o.dev",
      "description": "Micro DEV environment"
    },
    {
      "url": "https://api.m3o.com",
      "description": "Micro LIVE environment"
    }
  ]
}`
