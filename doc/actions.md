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
