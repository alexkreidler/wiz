# Wiz Processors

This document lists the default Wiz Processors and their implementation status.

Each Wiz Processor defines the structure of its input data, along with its configuration.

Additionally, each processor should describe its caching/hashing immutability mechanisms. Generally this is serialized in the **state** object. If a state object exists, each processor compares the most recently fetched version with it. The state object can also provide rules for comparing an old state object to a new object and determining if they have diverged in their versions. Or this can be implemented in the processor itself.

The worst-case scenario if the processor does not implement immutability mechanisms is that the output of the processor will automatically be hashed and stored in the processors state by the Executor. The executor determines which of these strategies to use: 
1. fully immutable/functional (the processor specifies this in its setup) -- the processor's output is a deterministic function of its input, no state is necessary. If this is done incorrectly then the pipeline's functional garuntees are destroyed
2. immutability state mechanisms - each processor stores enough state to be able to determine with some reasonable degree of certainty whether versions have changed. For additional flexibility, it could give a degree of certainty as a numeric value and then strategy 3 could be followed if it is below a threshold
3. Hashing - performed by the executor, stored in state

### idea
Each type:input processor is effectively loading a file from a remote storage system. For DBs, it is loading a record.
Thus, each can be represented as a remote reference to either a file or a record. 

## HTTP

The goal of the HTTP processor is to be able to access files, APIs, web pages, etc for processing in the pipeline.

```yaml
data:
    url: http://cnn.com/article # (required)
    # handle file writing (e.g. what should the file be named, just the document name?)
    method: GET # | POST | other HTTP verbs
    body: |
        http body content for POST requests, etc
    headers: # todo: add default headers
        Accept: text/html
```

The data object provides for all configuration about the specific request.

Effectively modeling a single file, how to model folders.

```yaml
configuration:
    concurrent: # concurrent enables concurrent download of a single file if the server provides the Accept-Ranges header
        enabled: true
        num: 3 # num_workers
```

How to distribute file parallelism, e.g. each request object: this should be transparent to the task graph.

The user should only be concerned with 1. creating the task graph 2. configuring each processor 3. providing the initial data

for file parallelism: data object representing a file is 1 request to a processor backend. This backend could be load balanced and distributed as long as fetches file and returns state. 

In fact, default Wiz Executor runs a single processor goroutine for each download request, the HTTP manager routine will launch other routines for actual data fetching (e.g. with the number of concurrent workers), others for computing state? this shouldn't take to long

However, with 10,000 files and 1 executor machine, it would launch 10,000 downloads at once which could clog the network. Thus we need Executor-level (actually really node-level) processor concurrency limits to not clog network, storage, etc.

If each processor was run as a Pod on k8s with set request quotas for CPU and memory, they'd automatically be limited. Add network as another resource to k8s?

This clogging seems like an environment/orchestration problem and not applicable for us rn.


## Git

Access files via Git

```yaml
data:
    url: https://github.com/alexkreidler/wiz # (required)
    branch: master # anything other than a commit or tag that can be fully determined to a specific commit will result in state contaning the current branch's latest commit. 
    origin: origin
    depth: 1
```

As written above, this will get an entire repository, aka a folder. How to specify specific files?

If writing simply `master`, or `branch-name`, the Git processor will store a simple state object only containing the latest commit on that branch.

If the branch is updated to a newer current commit, two things could happen:
1. `wiz get <package name>` - this will simply pull the state recorded in the lock file, but will print a warning saying which processors have changed state
2. `wiz get <package name> --update` - this will pull the latest versions of all processors and record the new state to the lock file, printing all processors that have changed state and asking the user for the version bump (default 0.0.1)

This allows the user to execute consistent versioning strategies. In fact, in a later version, we may allow complex versioning strategies, which bump by different amounts depending on which processors have changed. Each processor will have a `versionBump` which is the minimum applied version if that processor has changed. Thus, we will apply the highest version bump out of the `versionBump`s for all changed processors. We may even allow a more complex language for specifying bumps based on specific combinations of processors. `--auto` will use these automatic versioning strategies

Additionally, `wiz status, update, etc <package>` will return the status

```yaml
configuration:
    redirects: true # follow to the git repo
    # should depth be put here instead of in data? most use cases will just have depth one, want to access file at specific version
```

## Decompression/Extraction

TODO

## SQL in/out

## FoundationDB in/out