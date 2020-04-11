## Workflow

Workflow defines the developer UX and the expected step-by-step flow.

## Overview

The workflow is the developer experience when using the platform. This document defines all 
aspects of the workflow from login to new service creation to rolling updates and debugging. 
Anything related to the user experience is the workflow.

Our expectations of what the platform provides in relation to the workflow:

- Global deployments
- Auto configuration
- Simplified debugging
- Connect from anywhere

## Steps

Defining the steps in the workflow. This is the all encompassing experience when joining a 
company and gaining access to the source code and platform.

1. [Onboarding](#onboarding) - First invite, signup, login and access
2. [New Service](#new-service) - Creating your first service
3. [Pushing Updates](#pushing-updates) - Rolling updates when source code changes
4. [Debugging](#debugging) - Using the UI for debugging
5. [Querying](#querying) - Querying the service from anywhere
6. [Failure Cases](#failure-cases) - The failures that could occur and expected behaviours

Missing components in the flow

- Monitoring - User defined healthchecks
- Alerting - Paging/Email when things go wrong
- Access control - limiting who can do what
- Versioning/routing - Label based routing for feature development
- Usage tracking - Monitoring of platform and service usage
- Quotas/Billing - Per service/team quotas and tracking

## Onboarding

Think of onboarding as a developer joining the company or team. It's an identical experience here.

Onboarding steps

1. User is invited to GitHub team e.g Community
2. User navigates to Platform login page and signs in
3. User now has access to the platform and ability to create new service
4. User should be taken through the Platform 1:1

TODOs:

Login currently generates an Auth account but this is localised to a single region. Accounts should either 
be global or Auth should include a Login endpoint which is additionally called by the CLI for region specific login.

## New Service

Once the user has access they need a 1:1 experience for creating a new service.

1. Hit "new" service link and follow steps
2. Clone services monorepo (as would exist in a company)
3. Create new service locally (with own name)
4. Git push to the services repo
5. Follow flow in UI

TODOs:

We've seen user onboarding already fail either due to event issues, page refreshes or because the user flow breakdown. 
There's an expectation that all the instructions needed are on the one page. This may mean that we need to entirely 
remove the `micro run` step and auto deploy the service.

## Pushing Updates

When a user updates the source code we should automatically rollout the update. This may initially not exist but should become 
part of the flow. Because the runtime has the ability update services (by patching the deployment) we can trigger this 
via the built in Scheduler or calling runtime.Update.

1. User pushes to github and triggers the build
2. Build passed and curls /v1/github/events
3. Event is sent to platform service and it triggers runtime.Update

TODOs:

In future we create new service versions based on github branches for cross cutting feature development. Then 
delete them when the branches are merged/closed. Request flow can then push through via Micro-Version or something similar.

## Debugging

Once services are running the user should now have the ability to see that its running by viewing service endpoints, 
stats, logs and tracing info. This should be available via the UI and CLI.

1. UI redirects to service page after creating service
2. User can flip between tabs to view the service info
3. User can call `micro {stats, log, trace}` from the CLI to query the service directly

## Querying

Once the service is deployed the user should be able to query the service from the UI and the CLI. This should continue 
to mimic the existing user experience of `micro call` and the web dashboard form for calling services. The User should 
be pushed towards calling the service to ensure its actually working and behaving as expected (UI hint).

1. micro call from CLI (through proxy.micro.mu)
2. Platform UI call via services pages

## Failure Cases

We need to identify failure cases in the flow and rectify them. These are scenarios in which the default flow breaks down.

1. When events are not received the new service flow fails
2. When builds fail events are not received the flow fails
3. When the service fails to start events are not received
4. Any failure case needs to surface in the UI
