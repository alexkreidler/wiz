# wiz tasks

Wiz tasks is bulit in on a unique internal architecture.

The user interacts with the **Manager**, which is run in one of two modes

1. bundled with the CLI, persisting state in a local file/DB
2. a long running daemon with an API, able to serve requests via multiple clients (authentication needed)

The manager interacts with various **type Executor interface**s, which simply set up the execution environment (local or Kubernetes) through its APIs and spawn the Wiz Executor componenet `wiz executor` which serves a gRPC API.

The manager needs network communication with the executor components, and pushes tasks to it.

The executor component is made of
1. bundled default processors
2. the ability to run any native program available on the host

Both of these options can run on either the raw Wiz Data sent via protobufs or on Wiz File data

The core DAG model is as follows:

## Spec:

A **Pipeline** is a DAG that contains fully configured **Processors** (sometimews referred to as **tasks**). Each pipeline can be **run** with different sets of input data. A **run** represents one attempt at running the pipeline, and can either be in the **Running** or **Failed** or **Succeeded** state. Each **Processor** can either be in the **Idle**, **Running**, **Failed**, or **Succeeded** state.

For v1.0, any *process-level fatal error* in any one of the processors in the Pipeline can cause a run to fail. However, some processors may simply emit warnings, depending on their configuration.

The WTF does not (at this time) prescribe how logs, metrics, or any instrumentation on each of the processors should be structured, aside from the basic four-state indication for each processor. In the future, we may record additional information such as
1. time taken for each processor
2. types of data used, size of data input/output
3. metrics from the execution environment such as CPU and memory usage to use in more advanced heuristics to schedule tasks.

## Internals

Key Choice: Moved from ProcessorNode to *ProcessorNode implementing the graph.Node interface.

Now we can access it with `n *ProcessorNode := p.Graph.Node(id)` allowing modifications from any APIs including traversal.

## IDs

We use KSUIDs of the format: `1Vh9QAjN9Hh1CoRYg0RHbW1fCZp` after careful evaluation of all UUID formats.

It has a simple + powerful go library, 128 bits of randomness (slightly more than RFC 4122 UUIDv4) and time-locality resolution of 1 second.