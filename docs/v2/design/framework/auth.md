# Authentication

Micro needs an authentication story. In the beginning go-micro had no auth, on the premise that the base requirement 
for distributed systems was solely discovery and communication. Our default experience will continue to operate 
in a zero auth noop model until we can identify how zero trust will work.

## Overview

Auth will include both authentication and authorization. Authentication is the basis for checking whether a user 
or service is "logged in" or has an access token to use across the system. Authorization is used to check 
whether a user or service actually has the privileges to access a resource.

Our story always begins with

- go-micro interface with implementations
  * zero dep default experience
  * industry standard highly available system
  * micro rpc service implementation

The go-micro interface should interop with the rest of the framework then have the capability of being swapped 
out in production for a centralised system. The micro service implemenation enables an anti corruption layer 
to abstract away the underlying infrastructure and further usage through the surrounding micro ecosystem.

## Implemenations

- Zero dep - likely noop because it does not need to be included by default
- Service - go.micro.auth is responsible for managing rules and accounts.
- Casbin, Hydra, OPA - these are becoming open source standards for oauth/rbac and make the most sense here

## Interface

The auth interface provides the following methods. Each one is explained in detail below.
```go
// Auth providers authentication and authorization
type Auth interface {
	// Initialise the auth implementation. This must be called before any other methods are called.
	Init(opts ...Option)
	// Options returns all the options set when initialising the auth implementation, such as credentials etc.
	Options() Options
	// Generate creates a new auth account. The only required argument is ID, however roles, metadata and a secret can all be set using the GenerateOptions. Secret is not always required since it wouldn't make sense for some resources such services to have passwords.
	Generate(id string, opts ...GenerateOption) (*Account, error)
	// Destroy allows an account to be deleted
	Destroy(id string) error
	// RBAC (role based access control) is used for auth. Roles can be provided to an account on Generate as an option. Roles can be granted access to a resource, e.g. grant the role "user.finance" access to any endpoint on the service named "go.micro.service.reporting".
	Grant(role string, res *Resource) error
	// The inverse of grant, revoke removes a roles access to a resource.
	Revoke(role string, res *Resource) error
	// Verify takes an account and verifies is has access to a resource based on the RBAC rules. The implementation will keep a record of the roles granted access to the resource, it will then compare those roles to the roles given to the user and return an error if a match is not found.
	Verify(acc *Account, res *Resource) error
	// Inspect takes an access token (normally a JWT), inspects the token (this can be done client-side if the token is a JWT and the client has access to the public key), and then returns the account which the token was generated for.
	Inspect(token string) (*Account, error)
	// Token generates a new token for an account. Tokens contain a short-lived access token, which can be used to perform calls in the system and a long lived refresh token which can later be exchanged for a new token. Token requires some form of authentication, which is provided as a TokenOption, this can either be the accounts credentials (id, secret) or a refresh token which was provided by a previous call to Token.
	Token(opts ...TokenOption) (*Token, error)
	// String returns the name of the implementation, e.g. service.
	String() string
}

// Resource is an entity such a service. RBAC ensures anyone calling this resource has the necessary roles.
type Resource struct {
	// Name of the resource, e.g. "go.micro.store"
	Name string `json:"name"`
	// Type of resource, e.g. "service"
	Type string `json:"type"`
	// Endpoint of the resource, e.g. "Store.Read". We specify endpoint as this allows us to use RBAC at an endpoint level. '*' can be used as a wilcard to specify any endpoint.
	Endpoint string `json:"endpoint"`
	// Namespace the resource belongs to, the default is "micro" (auth.DefaultNamespace). Namespace allows for multi-tenancy RBAC, since there could be multiple versions of "go.micro.store" running in different namespaces.
	Namespace string `json:"namespace"`
}

// Account is a resource such as a user or a service who needs to make requests and be authenticated by micro.
type Account struct {
	// ID of the account. For users this is normally their email (e.g. 'johndoe@micro.mu') and for services this is normally their name (e.g. 'go.micro.store').
	ID string `json:"id"`
	// Type of the account, e.g. service. Account types should always be lowercase. 
	Type string `json:"type"`
	// Provider who issued the account, e.g. "oauth/google". This is currentlys used as additional information when auditing the account. 
	Provider string `json:"provider"`
	// Roles the account was provided with as a GenerateOption. These rules are used when doing RBAC.
	Roles []string `json:"roles"`
	// Metadata is a key/value map which can be used to store additonal information about the account, such as their name and avatar.
	Metadata map[string]string `json:"metadata"`
	// Namespace the account belongs to. The default is "micro" (auth.DefaultNamespace). This allows for IDs to be scoped to namespace and not need to be globally unique.
	Namespace string `json:"namespace"`
	// Secret for the account, e.g. the password or a secret persisted by the accounts provider.
	Secret string `json:"secret"`
}

// Token contains the credentials needed for an account to perform requests and refresh its identity.
type Token struct {
	// AccessToken is a short lived token provided on each request the account makes. This is either a JWT or a standard token (UUID V4).
	AccessToken string `json:"access_token"`
	// RefreshToken is a long lived token which is only used when calling the Token method to generate a new AccessToken.
	RefreshToken string `json:"refresh_token"`
	// Created is the time the token was created.
	Created time.Time `json:"created"`
	// Expiry is the time the access token will expire. The client will need to call the Token method before this time and replace this token with a new one.
	Expiry time.Time `json:"expiry"`
}
```

## Wrapper

go-micro/util/wrapper.go contains an `AuthHandler`. This handler wraps all incoming requests and is responsible for authentication and authorization. This wrapper is always enabled, even if MICRO_AUTH is not specified, however the noop implementation will verify all requests.

The wrapper will return immediately for any requests made against Debug endpoints, which are added to all go-micro services. We may change this in the future and extend RBAC to include Debug.

### Token

The first thing the wrapper does is check for an auth token. This is the `access-token` which is provided in the Token object. This token is passed in context as `Authorization` header and is prefixed by `Bearer ` (auth.BearerScheme). If the cookie `micro-token` is provided, micro web will set this in the context.

### Load the account

The token will be inspected using the auth.Inspect method to determine the account. If no account is retrieved, we fall back on a blank account for the scope of the wrapper, since the noop auth implementation will allow requests through regardless.

### Namespace

When a request enters the platform, micro web / api will determine the namespace. If the request host is micro.mu or any subdomain, the default namespace is used: `micro` (auth.DefaultNamespace). 

If a non-micro domain is used, e.g. `dev.m3o.app`, the subdomain is determined to be the namespace, in this example: `dev`. If a subdomain has not been used, we fall back to the default namespace.

Micro API / Web will set the namespace in the context using the `auth.NamespaceKey` key. The wrapper will check for this key, or fall back to the default namespace if it is not found.

The namespace is then compared to the namespace in the account. If the namespaces do not match, a forbidden error is returned and the occurrance logged.


### RBAC

Next, the wrapper will call the `auth.Verify` method to determine if the account has access to the resource it's calling. The resource is determined using a combination of the namespace calculated in the previous step, along with the service name / endpoint provided to the wrapper as part of the request.

If Verify disallows the request, the user is unauthorised and the request is terminated. 

### Setting the account

Handlers often need to access the auth account to retrieve data such as account ID, so auth provides two helper methods to enable this: `AccountFromContext` to set the account in the context (used by the wrapper) and `ContextWithAccount` (optionally used by the handler).
