# Wiz API tooling

We use API Blueprint and Swagger v2.0 for its code generation capabilities.

We've built a pipeline that goes from API blueprint to swagger, and we hope in the future the tooling is good enough that can be our sole system.

However, it has a few quirks so as of the latest, we are just using the Swagger document as our source of truth.

## State of the union

Swagger is horrendously tedious, at least on the server side, so we're moving to either Protobufs or CapNProto for the Wiz Executor. Then, we'll also add support for external REST processors in the future.