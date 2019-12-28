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

Super good Go concurrency resource: https://notes.shichao.io/gopl/ch9/

## Work handling models

There are two possible methods, the Pool method or the Lifetime method

A few requirements for both methods

1. be able to limit the number of total workers running at a time - this is necessary when a worker is resource intensive, e.g. takes up bandwidth or GPU memory. Spawning more workers after this point would either lead to a drop in performance or even an error.

Ways to evaluate the methods: 

1. overuse/underuse of: memory, CPU, etc at any given time
2. additional startup costs/overhead
3. cost/latency on actual new request/chunk

### Pool method

We follow the Go worker pool pattern similar to what is described here: https://gobyexample.com/worker-pools

Basically we have a pool of ProcessorRunners with associated Processors. The runners listen on a channel for data, and thus data is distributed across the Runners. The number of runners can be configured in the ExecutorConfiguration.

This may
1. overuse memory and CPU at any time because workers could be idle and
2. cause startup overhead for starting the pool, especially if they do computationally intensive tasks on startup.
3. may reduce latency for new chunk if workers are expensive to start and a worker is available. If no worker is available it will have to wait for a worker. Also if workers are not expensive to start the latency benefit is negligible. 

### Lifetime method

We simply spawn a new Processor for each new Chunk we receive. 

We can then limit the total number of Processors via a configuration parameter to satisfy requirement 1.

For evaluation it will
1. always use the least amount of memory/CPU needed
2. not have any startup cost
3. have a constant lateny on a new chunk unless we reach the worker limit, at which it will have to wait for a worker to finish to be able to process.


### Decision
For now, we will just use the Lifetime method because it is conceptually the simplest