# Processor API

This document defines the Wiz Processor API.

## States
First, the state of a given processor is as follows:
- nonexistent --> create
- uninitialized --> configure
- configured [--> provide data, -> unintialized (error)]
- validating [-> configured (error), -> running]
- running [-> configured, --> data done]
- success
- failure

In the list above, the format is as follows
- API Action format: state_1 --> api_action (transition_to_state, default next item in list)
- Transition format: state_2 -> transition_to_state (transition_name)
- Either format: state_3 [format_1, format_2]

The following sections describe each state and the transitions to the next states

## Nonexistent

The processor has not been started and the Processor API is not yet available. This is not an explicit state in the procesor.

## Uninitialized

The processor has not yet been configured, and is a useless state. A processor can be configured with blank data (which is applied to the defaults) but it must be done explicitly.


## Configured

The processor has been configured successfully and is now ready to recieve data.

The following sections happen for each /data/:dataID that is provided, as each piece of data can be handled in parallel if the configuration allows.

## Data states
### (optional Determining)
The processor references external data or uses a nondeterministic/PRNG process that requires state to track deterministically. 

Thus in this stage the state is determined.

### Validating

The data handler has just recieved data and is validating that it will work with the current configuration.

It will automatically move to running state or return a non-fatal error that the data provided is invalid (resulting in a `CompletedWithErrors` terminal state)

### Running

The processor is running on the chunk of data. It will either

1. successfully complete on that chunk, and exit the data state (e.g.) return to the configured state
2. fail on the chunk, and depending on configuration either: log a warning and exit data state, or go into a Failed state, which once read to the manager, will terminate the processor

### Waiting/Ready
The waiting data state indicates that a data chunk has been processed successfully but has not been read by the manager or next processor. 

Thus, until the data is flushed, it must be kept around.

In the HTTP Data API version, this is until a get request at `/data/:dataID/output`, however this could be different for Kafka and other APIs by doing so automatically.

## Terminal States
The following states occur at the very end of the processors lifetime and the processor can terminate itself once these states are read by the manager.

### Success
Once the processor recieves a special API call that signifies all data has been sent, and if no errors have occured, it will resolve to the Success state.

### CompletedWithErrors
This signifies that the processor completed but had non-fatal errors. This happens for example

### Failed
This signifies a fatal error occured and the processor could not continue.

## Additional considerations

Some processors that are distributed across nodes will not follow the simple "data chunk parallelization" pattern, where a given set of processors with the same set of configuration can return the correct result loadbalanced by data chunks. For example, an ML algorithm may need data from multiple chunks.

Thus the processor must specify: `loadBalance: false` in its configuration so that the Wiz Executor does not automatically do this.

Some processors will need shared state across nodes: Etcd?

### Data APIs

By default, the data chunks specified above will come from simple HTTP requests that contain raw data like JSON or references to Wiz Filesystem data.

However, in the future we may allow changing the underlying data transport mechanisms to systems like Kafka, etc. However, the same state machine should work for these systems.

### Go library for building processors

We provide a simple Go library that defines all of these required handlers for various transitions as a simple interface to implement. The library then supports serving the appropraite HTTP interface to provide a standalone API for a custom Go processor. 

The library also supports multiple named processors for easy extensibility.

### Tools for building processors

All of this Processor API is also written in machine-readable form in an API design language that has support for generators in various languages, allowing you to quickly bootstrap your own processors in any language.

However, as the transport mechansim behind this API may change from HTTP to gRPC, CapNProto, or Apache Arrow RPC, we wanted to write it in a protocol-agnostic API IDL. Of course, no such thing exists.

In reality, the core Processor API for configuration and reading state will most likely always be in HTTP. However, the Data API for ingesting data will likely be in an Apache Arrow based format.

### extra

Data chunks are allocated with unique IDs by the manager, and they are global to the pipeline. They are also deterministic and are the basis of Provenance, which tracks how data is transformed at a granular basis.

### Driving of the state machine

This API is written using a state machine model, but all actions are driven by the manager.

This could cause a bottleneck at scale as the manager would need to 
1. Configure each node
2. Send it the raw data
3. Poll for updates on state until Waiting/Ready state
4. Recieve the new data, and send it to the next node/processor

However, the state machine concept can still be used if each Processor drives the majority of its own state transitions. For example:
1. Processor fetches its own config from manager
2. Fetches the raw data if needed from manager
3. When ready, it gets its downstream node/processor and sends it to it

This means the processors would drive themselves, and the manager would only be in charge of recording state changes, chunk IDs somehow, and persisting the record of all the transitions that the processors do.

In fact, at the terminal state of a processor, it could simply POST all of its state history and provenance data. This would cause a delay in fresh information about the state of the graph by the manager

It could also POST state to the manager whenever there is a transition

Either of these methods have drawbacks, including how each is authenticated and authorized.

In model 1, each processor simply has to verify that the manager is who it says it is and is trusted. Since the manager spawns the processor, this should be simple by just providing a nonce or token.

In model 2, the manager has to authenticate each processor and then authorize it to only modify state information about itself. 

Model 1 is more of a classic REST design whereas Model 2 is a streaming design that lends itself to gRPC and streams


 //"https://json-schema.org/draft/2019-09/schema",