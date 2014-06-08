# Actions

 - [What is an Action?](#actions-services-and-charms)
 - [Charm Creators](#charm-creators)
   - [Actions on the Charm](#actions-on-the-charm)
 - [Frontend Hackers](#frontend-hackers)
 - [Backend Hackers](#backend-hackers)
   - [Client API for Actions](#client-api-for-actions)
   - [Lifecycle of an Action](#lifecycle-of-an-action)
   - [State Machinery](#state-machinery)
     * [Actions Collection](#actions-collection)
     * [Actions Results](#actions-results)
     * [Actions Log](#actions-log)
   - [Param Validation with gojsonschema](#param-validation-with-gojsonschema)
   - [Actions at the Service Level](#actions-at-the-service-level)
     * [Actions Watcher](#actions-watcher)
     * [Uniter Loop](#uniter-loop)
     * [Hook Environment](#hook-environment)
 - [Juju Users](#juju-users)

TODO: spruce up formatting?

---

# Actions, Services, and Charms

An Action is an executable defined on the Charm.  Actions are controlled 
and observed by the Juju client using a well-defined set of API endpoints,
which are called either via CLI or web frontend.  This document contains the
[definition](#actions-services-and-charms), [client operation](#client-api-for-actions), [lifecycle](#lifecycle-of-an-action), and [technical](#state-machinery) [details](#actions-at-the-service-level) of Actions.

---

# Charm Creators

## Actions on the Charm

Actions are defined as [Hooks](charms-in-action.md#hooks) on the Charm.  An Action is simply an executable
script or file, which runs in a [Hook environment](charms-in-action.md#execution-environment); therefore, certain
environment variables and calls can be used to interact with the environment.

The Charm author must also define an actions.yaml file in the Charm root
directory.  This file must begin with the key "actions:"; after that, each
Action must be listed with a "description:" key and an optional "params:"
key.  If "params:" is given, it must contain a schema compliant with
[JSON-Schema draft 4](http://json-schema.org/latest/json-schema-core.html).
See http://json-schema.org for examples and more details.

This document is used to define the schema which is used to validate the
parameters passed with inbound Action commands.  See [Validation](#param-validation-with-gojsonschema).

Names of actions and parameters must start and end with a lowercase alpha
character, and may only contain hyphens and lowercase alpha characters.

Example:

```yaml
#
# sample actions.yaml
#

actions: 
   snapshot:
      description: Take a snapshot of the database.
      params:
         outfile:
            description: The file to write out to.
            type: string
            default: foo.bz2
         compression-type:
            $schema: http://json-schema.org/draft-04/schema#
            title: Compression type
            description: The kind and quality of snapshot compression
            type: object
            properties:
               kind:
                  description: The compression tool to use.
                  type: string
               quality:
                  description: Compression quality from 0 to 9.
                  minimum: 0
                  maximum: 9
            required: [kind]
   kill:
      description: Kill the database.
```

TODO: Fill in details of what environment calls and variables exactly are
available via the Hook environment.

--- 

# Frontend Hackers

TODO: Fill me in!  Get talking with frontend crew!

---

# Backend Hackers

## Client API for Actions

## Lifecycle of an Action

## State Machinery

### Actions Collection

### Actions Results

### Actions Log

## Param Validation with [gojsonschema](http://github.com/binary132/gojsonschema)

## Actions at the Service Level 

### Actions Watcher

### Uniter Loop

### Hook Environment

---

# Juju Users

TODO: Fill me in!


***COMMENTS REQUESTED *** (especially from frontend hackers, see item 3)

I think there are three main consumers with different needs that this doc should address.  I think the doc should have a [TOC] at the top, with a #primary header for each consumer, addressing the following needs perhaps in ##secondary headers.

1. Juju Users
    - I think we have a good breakdown now.
    - Users need to answer a specific question, so I think we should ###sub-header each command.
2. Charm Creators
    - Charm Creators need to know how to define an Action on the Charm.  They need to know:
        * How to define a Charm Action script (i.e. it's a Hook)
        * How the Charm Action will interact with the Hook environment (e.g. ActionGet)
        * The requirement to create actions.yaml to specify the Charm's Actions' parameters.
3. Frontend Juju Hackers
    - They need a clear spec for the client API function by function.
    - We should ask Rick and Jeremy what API endpoints the frontend team needs.
        * Actions-Schema getter?  (serialized JSON-schema for GUI construction?)
        * What kind of data comes back from juju queue?  What about juju status?
    - They need a clear UX spec for various possible interactions.
        * What API endpoints are synchronous or asynchronous?
        * Other GUI questions should be discussed in mailing list with Frontend Hackers.
4. Backend Juju Hackers
    - A sketch from client API to State to Hook and back.
    - A sketch for Charm Actions.  Clarity between:
        * Actions schema (i.e. Charm Actions, and gojsonschema's role)
        * Lifecycle of "actual" Actions (e.g. a snapshot action)
        * Action hooks and their operation
        * Action responses and the Action log
