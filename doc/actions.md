[//]: # (TODO: actions at the s/Service/Unit/g level)
[//]: # (TODO: make sure validation is clear)
[//]: # (TODO: what do charm authors need to know?  slice out technical validation)

# Actions

 - [What is an Action?](#actions-services-and-charms)
 - [Charm Creators](#charm-creators)
   - [Actions on the Charm](#actions-on-the-charm)
 - [Frontend Hackers](#frontend-hackers)
 - [Backend Hackers](#backend-hackers)
   - [API](#api)
     - [Client API for Actions](#client-api-for-actions)
     - [Internal API Methods (unit-facing)](#internal-api-methods-unit-facing)
     - [Details](#details)
   - [Lifecycle of an Action](#lifecycle-of-an-action)
   - [State Machinery](#state-machinery)
     - [Actions Collection](#actions-collection)
     - [Actions Results](#actions-results)
     - [Actions Log](#actions-log)
   - [Param Validation with gojsonschema](#param-validation-with-gojsonschema)
     - [Action Parameter Validation](#action-parameter-validation)
     - [Charm Actions Schema](#charm-actions-schema)
   - [Actions at the Service Level](#actions-at-the-service-level)
     - [Actions Watcher](#actions-watcher)
     - [Uniter Loop](#uniter-loop)
     - [Hook Environment](#hook-environment)
 - [Juju Users](#juju-users)

---

# Actions, Services, and Charms

An Action is an executable defined on the Charm.  Actions are controlled 
and observed by the Juju stateserver using a set of API endpoints,
which are called either via CLI or web frontend.  This document specifies the
[definition](#actions-services-and-charms), [client operation](#client-api-for-actions), [lifecycle](#lifecycle-of-an-action), [stateserver details](#state-machinery), and [service-level details](#actions-at-the-service-level) of Actions.


---

# Charm Creators

## Actions on the Charm

Actions are defined as a special case of [Hooks](charms-in-action.md#hooks) on
the Charm.  An Action is simply an executable script or file, which runs in a
[Hook environment](charms-in-action.md#execution-environment).  This is how an Action gets the params it is
called with.  It can also mark itself failed, or add results to a map to be
returned to the stateserver upon completion.  See [Hook Environment](#hook-environment) for
details.

The Charm author must define an actions.yml file in the Charm root
directory.  This document is used to define the schema which is used to
validate the parameters passed with Action commands.  See [Validation](#param-validation-with-gojsonschema).
The Charm author must also put the Action executables or scripts in the
`actions/` directory in the charm.  Their names must match the names given in
the map as directed below.

`actions.yaml` must contain a YAML map which conforms to the following pattern:

 - Top level key MUST be `actions`.
 - Next level keys MUST be the names of the Action scripts for the charm.
 - Each of these Action names MUST be the top level key of a map.
 - Each Action MUST have a string as its name, which MAY contain but MUST NOT
   begin or end with hyphens.  Each Action MUST have only lowercase names.
 - The map for each named action MUST have ONLY keys `description` and `params`
   and SHOULD have both, but both keys MAY be omitted.
 - `description` MUST have a string value.  This value SHOULD be the short
   description of the Action.
 - `params` MUST be the top level key of a map, but the map MAY be empty.
 - The map for `params` MUST specify [JSON-Schema](http://json-schema.org/)
   for the parameters for the action, as in the following example.
 - The map for `params` MUST NOT contain a `$schema` key.  `$schema` is not
   supported by Actions at this time.

Example:

```yaml
# actions.yaml
actions:
  snapshot:
    description: Take a snapshot of the database.
    params:
      title: Snapshot params
        description: Take a snapshot of the database.
        type: object
        properties:
          outfile:
            description: The file to write out to.
            type: string
          quality:
            description: Compression quality
            type: integer
            minimum: 0
            maximum: 9
        required: [outfile]
  kill:
    description: Kill the database.
```

> Note that the value for `params` is [a JSON-Schema conformant](http://json-schema.org/examples.html) YAML map.

> Note also that `kill` contains no `params` definition.  This Action takes no
  arguments.

This would support an Action call such as:

`$ juju do <unit name> snapshot --params snap.yml`

where `snap.yml` is:

```yaml
outfile:
  out-2014-11-03.bz2
  quality: 5
```

See [Validation](#param-validation-with-gojsonschema) for more details.

--- 

# Frontend Hackers

---

# Backend Hackers

## API

### Client API for Actions

> Documentation in progress

- Enqueue
- ListAll
- ListCompleted
- ListPending
- Cancel
- ServiceCharmActions
- WatchActions

### Internal API Methods (unit-facing)

- GetActions
- StartActions (maybe UpdateActions, or ActionsUpdate?)
- FinishActions
- WatchActions

### Details

#### Client Facing

##### Enqueue

`Enqueue` takes a list of `Actions` and queues them up to be executed by the
designated receiver, returning `ActionResults` for each queued Action, or an
error if there was a problem queueing up the Action. 

##### ListAll

`ListAll` takes a list of `Tags` representing receivers and returns all of the
`Actions` that have been queued or run by each of those entities as
`ActionsByReceivers`.

##### ListCompleted

`ListCompleted` takes a list of `Tags` representing receivers and returns all
of the `Actions` that have been run on each of those entities as
`ActionsByReceivers`.

##### ListPending

`ListPending` takes a list of `Tags` representing receivers and returns all of
the `Actions` that are queued for each of those entities as `ActionsByReceivers`.

##### Cancel

`Cancel` attempts to cancel queued up Actions from running. `ActionTags` are
passed as argument, `ActionResults` are returned.

##### ServiceCharmActions

`ServicesCharmActions` returns the `ServicesCharmActionsResults` for the
passed `ServiceTags`. 

##### WatchActions

:question:

#### Unit-facing

##### GetActions

:question:

##### StartActions

:question:

##### FinishActions

:question:

##### WatchActions

:question:

## Lifecycle of an Action

An Action is enqueued on State using the [`Enqueue`](#enqueue) API method. 
Its parameters are [validated](#action-parameter-validation) before it is enqueued.
If invalid, an error is returned; if valid, it will be successfully enqueued.
A listener for the unit targeted by the action will obtain a diff containing
the action's tag.  The unit uses this listener in Filter to receive queued
Actions, and passes a HookInfo to Uniter with the ActionTag.  Uniter calls 
runAction with the HookInfo, and runAction validates the params against its
version of the charm.  If the params fail to validate, an appropriate error is
returned.  Failed or errored Actions will not trigger a Uniter error state.

If successful, the named action executable in /actions will be run.  This
script may retreive its parameters by executing `action-get`.  It may set a
fail state and message using `action-fail`, and it may set a response map
using `action-set`.

When it is finished, its results will be returned to State.  A client can use
the [WatchActions](#watchactions) API method to get a watcher which will return the tag.
The client can then query State using the results tag to get the results of
the Action.

## State Machinery

### Actions Collection

### Actions Results

### Actions Log

## Param Validation with [gojsonschema](http://github.com/juju/gojsonschema)

### Action Parameter Validation

Action schemas are defined by the Charm in order to validate parameters
passed to Actions, and to define the GUI by which Action parameters are
defined.  The schema is loaded from YAML in the Charm and validated at
runtime, and is used to build a [JSON-Schema](http://json-schema.org) document using
[gojsonschema](http://github.com/juju/gojsonschema).

### Validation

Any JSON document defining an Action params map may be validated against its
respective Actions.ActionSpecs.Params[<action name>] map using gojsonschema
as follows:

```go
imports ( 
//	...
	"github.com/juju/gojsonschema"
	"github.com/juju/juju/charm"
//	...
)

// ...
	validationDoc := gojsonschema.NewJsonSchemaDocument(someCharm.Actions())
	validationResult := validationDoc.Validate(actionArgsJson)
	if !validationResult.Valid() {
		for i, vErr := validationResult.Errors() {
			// ...
		}
// ...
```

This technique is used to validate params when an Action comes into State
from the API, as well as to validate params upon delivery to the Unit as the
real states of those entities may change in flight.

### gojsonschema 

[gojsonschema](http://github.com/juju/gojsonschema) defines several useful types and methods which can be used
to validate or invalidate incoming JSON param maps against the Action's defined schema.

 - [JsonSchemaDocument](https://godoc.org/github.com/juju/gojsonschema#JsonSchemaDocument)
   * func (*JsonSchemaDocument) [Validate](https://godoc.org/github.com/juju/gojsonschema#JsonSchemaDocument.Validate) (document interface{}) *[ValidationResult](https://godoc.org/github.com/juju/gojsonschema#ValidationResult): Validates a JSON document against the defined JSON-Schema.
 - [ValidationResult](https://godoc.org/github.com/juju/gojsonschema#ValidationResult)
   * func (*ValidationResult) [Errors](https://godoc.org/github.com/juju/gojsonschema#ValidationResult.Errors) [][ValidationError](https://godoc.org/github.com/juju/gojsonschema#ValidationError): Returns a list of [ValidationError](https://godoc.org/github.com/juju/gojsonschema#ValidationError)s which occurred during the validation.  These can be introspected or simply have String() run on them.
   * func (*ValidationResult) [Valid](https://godoc.org/github.com/juju/gojsonschema#ValidationResult.Valid) bool: was the passed JSON document valid?
 - [ValidationError](https://godoc.org/github.com/juju/gojsonschema#ValidationError)
```go
type ValidationError struct {
    Context     *jsonContext // Tree like notation of the part that failed the validation. ex (root).a.b ...
    Description string       // A human readable error message
    Value       interface{}  // Value given by the JSON file that is the source of the error
}
```
   - func (ValidationError) [String](https://godoc.org/github.com/juju/gojsonschema#ValidationError.String) string: returns a human-readable string describing the validation error.

### Charm Actions Schema

The charm/actions.go Actions type is defined as follows:

```go
type Actions struct {
        ActionSpecs map[string]ActionSpec
}

type ActionSpec struct {
  Description string
  Params map[string]interface{}
}
```

See [Actions on the Charm](#actions-on-the-charm) for details on defining a
Charm's Actions schema.

The "Params" map for each Action in the list must conform to JSON-Schema,
since it is what the validation documents are built against.  When the Charm
is loaded via Bundle or Dir, [ReadActionsYaml](https://github.com/juju/juju/blob/master/charm/actions.go#L38) loads an Actions struct
from YAML and checks to ensure that it conforms to JSON-Schema Draft 4.  When
the Charm is added to the Charm representation in State, the Actions type is
serialized into the database for later use.  A Charm's Actions spec can be
retrieved from the Charm via [Charm.Actions()](https://github.com/juju/juju/blob/master/charm/charm.go#L16) on an implementing type.

Note that any types which implement Charm must also implement Actions().

## Actions at the Service Level 

### Actions Watcher

### Uniter Loop

### Hook Environment

When an Action script runs, it runs in a [Hook environment](#charms-in-action.md#execution-environment).  This means
that any Action script or executable has access to special commands from the
environment it runs in:

 - `action-get [<key>.<key>...]` -- retrieve the params passed with an action.

   Example:

   ```bash
   $ action-get outfile.name` 
   "foo.bz2"
   ```
 - `action-set <key>[.key.key]=<value>[ ...]` -- insert or replace values in a map to be returned after completion.

   Example:

   ```bash
   $ action-set outfile.size=10.2G success=true
   ```

 - `action-fail [<message>]` -- set the action to failed with a message.

   Example:

   ```bash
   $ action-fail "I'm afraid I can't let you do that, Dave."
   ```

---

# Juju Users

`$ juju actions mysql`

`$ juju do mysql/0 snapshot --params snap.yaml`
