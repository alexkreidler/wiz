# Wiz Processor Executor

This component of Wiz is a high performance golang binary that exposes a Wiz Processor API Server that handles all built-in and default processor requests.

Some architecture choices:

The general structure of the executor is the `runMap` type, which is a 2D map of ProcessorIDs to RunIDs to actual `runProcessors`, which are the instances of the processors for a specific run.

All of the operations use value receivers so they are concurrency safe and can be used by many clients. However the `runMap` is just a regular Go map so it is not concurrency safe at this time.

TODO: replace the runMap with a concurrency safe data type so many clients can manipulate at once.

It should be safe to use in multiple HTTP reqs.

Usually, builtin maps are fine for concurrent reads but a read and write at same time on 1 key can result in race conditions.

E.g. someone writes a config creating a new run, and someone tries to get all runs. These happen at the same time and we could get conflicts.

Concurrency is not a huge deal RN as for testing everything will be linear.

TODO: think about concurrency for data, e.g. how much to push through.