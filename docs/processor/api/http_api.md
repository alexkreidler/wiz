# Wiz Processor HTTP API

This document defines the Wiz Processor API as implemented using JSON over HTTP.

## Basics
The api is served at `/processors/v1/processor_name`, where `processor_name` is the unique global ID of the processor. 

In this document, all URLs are sub-paths of that base URL.

Additionally, as the Wiz Manager could likely schedule multiple Pipelines or Tasks in the same Pipeline that use the same Processor Task on one physical/logical Processor, the Processor MUST support partitioning all of its logic including Configuration and State by a Manager-scope globally unique **Run ID**

Metadata: this is information about how the processor identifies itself to the Manager

`GET /`
```json
{
    "name": "processor_name",
    "version": "0.1.0", //this is a semver version
    
    // we may allow generic metadata in the future, maybe under meta key?
    "tags": ["tag1", "tag2"],
    "author": "etc"
}
```
<!-- //"#/components/schemas/processor_metadata" -->

`GET /runs` - this is used for determining the number of runs on one processor, but keep in mind processors can be terminated after all runs have stopped
```json
{
    "runs": [
        {"runID": "aafhueiwlahufil"},
        {"runID": "aafhueiwlahufil"},
        {"runID": "aafhueiwlahufil"}
    ]
}
```

Getting state:

`GET /:runID/` returns
```json
{
    "state": "UNINITIALIZED" // or any of the other states
}

or

{
    "state": "CONFIGURED", // or any of the other states
    "configuration": { ... }
}
```

The only required fields for the processor are state.

## The State Sequence

A series of API calls that transition the Processor from one state to another

### `POST /:runID/configure`
Body:
```json
{
    "configuration": { ... }
}
```
Only the configuration field is required.

Response

```json
{
    "state": "configured",
    "configuration": { ... }
}

OR

{

    "state": "uninitialized",
    "error": "blabla"
}
```

### POST `/:runID/data/:chunkID`

Data is provided to the API in chunks, with globally unique **chunkID**s (these are tracked via Provenance)

Body:

```json
{
    "data": {
        "type": "raw",
        "value": { ... RAW JSON ... }
    }
}

OR

{
    "data": {
        "type": "file",
        "value": {
            // This is the only driver understood by any processor, and is simply using the local filesystem API and drivers
            // They can simply open the file with their open() syscalls or language file libraries
            // The executor, upon intercepting a request that does not use the local driver, is responsible for syncing it to be local.
            "driver": "local",

            // This location is relative to the locally provisioned Filesystem Storage Area, not absolute (actually, thats BS, it should be absolute -- e.g. k8s mounts)
            "location": "/bls/series.file"
        }
    }
}

OR

{
    "data": {
        "type": "external_api",

        "value": { 
            "api_type": "apache_arrow", //Kafka, etc,
            "configuration": {
                ... // API Specific configuration to either fetch the data or receive it on a port e.g.
                // "listener": "0.0.0.0:7894"

                // TODO: deal with the fact that most external API types will be streams that provide multiple chunks
                // Also, use watermill library
            }
         }
    }
}
```

Response:

this is assuming determining and validating are all synchronous which is not the case

we'll really just return the state of each chunk

```json
{
    "state": "running",
    "determined": true, //can also be false if this is a deterministic processor
    "validated": true
}
```

### GET `/:runID/data`

This returns all data chunks a given run has processed. Keep in mind after the run is GCed these will go away.

Should return list of chunk IDs and their states


## API IDLs and design choices

In the future, we want to build all of our APIs with JSON Schema and JSON HyperSchema. We may reevaluate Hydra as well (uses JSON-LD).

However the ecosystem for these tools is so far behind Swagger that it is pointless to use them effectively now.

**Notes for VSCode**: VSCode supports JSON Schema and Hyper-Schema version `draft-07`, as long as they are referenced using `$schema` over an `http://` not `https://` URI.