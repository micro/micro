# Tenancy

Micro is moving towards becoming a multi-tenant system by default.

## Overview

[Multi-tenancy](https://en.wikipedia.org/wiki/Multitenancy) is the ability to support multiple tenants in a single system. 
Micro is now moving towards supporting multi-tenancy. We need the ability to define multi-tenancy in a clean way that 
does not require the developer or user to deal with the issue. We should be able to segregate data and services belonging 
to different users, customers and teams paying for isolated platforms.

## Design

Firstly, multi-tenancy is implemented by Micro and non-concern of the go-micro framework. Micro (the runtime) has an understanding of tenancy and is responsible for it, the services built on top of the runtime however should not need to have any notion of tenancy; whilst the runtime is mutli-tenant, services are single-tenant. Tenant is dependant on identity and therefore requires auth to work. If the auth implementation used is noop, the runtime should assume we're operating in a single tenant mode and default to using the default namespace as defined in `micro/internal/namespace` (at the time of writing this is currently *go.micro*, but we will likely move to just *micro* soon.).

We've found that go-micro has to at the very least support the options to be configured in a way that will enable multi-tenancy in micro itself e.g the store must support specifying database/table at the time of a query rather than just on initialisation. The runtime and registry must support passing through namespace or prefix and the router must be able to segregate networks.

## Access

Below we discuss access cross namespace and through entrypoints

### Cross-Namespace

As a rule, tenants cannot access any services outside their own namespace (each tenant is a "namespace" in micro), however there are some exceptions. All tenants can read/write to the default namespace, but tenants can only call (read) the runtime namespace, not deploy or amend services (write). Whilst the registry will restrict the services returned in a given namespace, it will be the responsibility of auth to enforce the rules noted above.

### API / Web

Micro API and Web both have a namespace flag, which now can be used to filter the services using a prefix (e.g. "go.micro.web" will only return the services in the format "go.micro.web.X"), or if the value "domain" is provided, they will determine the namespace on a request-by-request b asis using the domain. As a fixed (hardcoded) rule, all micro.mu and development hosts will use the default namespace, and other subdomains will use the subdomains to determine namespace. For example: "foo.m3o.app" will use the "foo" namespace, "staging.myapp.com" will use the "staging" namespace and a top level domain such as "myapp.com" will use the default namespace.

## Scoping

### Determining the Namespace

Namespaces are set in the request header at the start of the request by the auth wrapper (`micro/auth/wrapper`), which is injected into both the API and Web servers. All runtime services can use the `internal/namespace` package, `namespace.FromContext` function to determine the current namespace. If the request originated from an external call, the namespace key will be set on the context and this will be returned. If the request was a service=>service request, we will get the namespace from the auth account of the calling serviice. If the namespace is still unknown (an unauthenticated service made a request), we fallback to the default namespace.


### Namespace in Core Services

Whilst core services need to be responsible for managing multi-tenancy, they do not need to know how it's determined. They should simply scope the resources a context can access using the contexts namespace: `namespace.FromContext(ctx)`. Resources should be persisted so that the unique identifiers do not need to be globally unique, but only unique to the namespace. For example, if namespace A writes a config key "foo", it should not conflict with the same key previously written by namespace B, this goes back to services not needing to know about mutli-tenancy.

There is an exception to the above rule: the registry. Because micro service names use the format [namespace,type,alias], we already have the concept of namespace baked into go-micro. Hence, if a service in the "foo" namespace tries to register a service named "bar.web.x", they'll get a forbidden (403) error. In the future, we will likely move away from this as it prevents using mutli-tennancy for staging and test enviroments (in these scenarios the same services could exist in mutliple namespaces with the same name).

## Framework

One core concern we have is not conflating namespace/tenancy into Go Micro itself. Doing so adds a level of complexity we don't want to have to deal with here. Seeing Go Micro as both a foundational building block for distributed systems and the framework for writing services, neither really has any concern with tenancy of a platform.

Therefore in Go Micro we propose that tenancy is really geared around the support of configurability of interfaces and their domain e.g the Store interface has moved from generic key-value in a single flat namespace to supporting database/table as an option which enables us to scope usage to a specific database and sets of tables.

### Options

Below is a non-comprehensive list of how we map namespace/tenancy to configurable options in Go Micro

Interface | Options | Description
--- | --- | ---
Auth  | Scopes  | Logical scoping of an account to define is access policy
Broker | Exchange | The exchange through which the all routing occurs for a topic
Config | Environment  | The environment specific configuration for apps or services
Store | Database  | The database to read/write from
Registry | Domain | The domain to which the services belong
Router | Network  | The isolated network for a set of services
Runtime | Namespace | The boundary of isolation for processes
