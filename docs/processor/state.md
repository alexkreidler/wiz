# Wiz Processor State + Status

**State** is an internal processor representation of any sources of non-determinism in the pipeline. This includes:
- external data sources (they can be overwritten, versions must be tracked)
- pseudorandomness (systems that use PRNGs must store the seed in state for deterministic outputs)

For example, an Apache HTTP server serving a directory of files can overwrite a file with a new version that may have different properties. 

The HTTP processor will then use properties of the HTTP response (e.g. `Last-Modified`, `ETag`, `Content-Length`, etc headers)

While none of these are 100% accurate, they are appropriate proxies to track the real **version** of the file.

If none of these are available, or the configuration requires a degree of certainty that these cannot provide, then we fall back to a simple hash of the file once it is downloaded. Wiz Executor has an optimized implementation of this and be configured to do it automatically given a certain response in the **state**. The hash is then stored as the state.

In Git for example, if an external data reference invokes `branch: master`, it could be updated to a new version. Thus, git simply stores the most recent commit SHA as the state. This ensures all the garuntees of Git apply to the data source. (however, branches can still be force pushed)

## Updating data

When a new version is available, there are two modes. The mode is specified by the user for a pipeline, or can be applied on a processor basis. They are:

- Update mode - updates the state to reflect the new resource
- Stale mode - tries to fetch the old resource. Actually this may not work as described in failure modes, so we just cache all resources with state in a hash table of their state and the local file reference.

## Failure modes

When trying to fetch an HTTP file that has been overwritten, or a branch that has been force pushed, or a repository deleted, and in many other cases, the version of the external resource requested may simply not be available

This can result either from
1. uninitialized state - e.g. requesting a file or repo that doesn't exist
2. Stale mode - trying to fetch a version that has been overwritten by the external system
3. Update mode - the entire repo or site has been taken offline, no longer available

Each option responds differently depending on which mode it came from, and this can be (you guessed right) configured!

## State Machine?
Wow this is now sounding very complex. State can be uninitialized, stale, or up to date, and there are many transistions between each state. These can each be configured to act differently.

Sounds like a state machine!

In fact, there is the **State Status** and the **Processor Status**. Both of these are states that have transitions and are modeled by state machines. We call states **Status**es just to prevent confusion with **State** which is the version of an external resource.   