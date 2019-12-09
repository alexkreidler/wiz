# Wiz Task Framework Specification

aka. Wiz ETL

Overview: the ability to easily configure and perform a DAG of data transformations on big data with streaming features and multiple representations.

it appears Beam solves this: https://beam.apache.org/documentation/programming-guide

## Specification components

1. A spec for a the generic architecture of Wiz ETL
2. A spec for the configuration of individual processors
3. A set of basic processors and their configuration

## Implementation

The Wiz ETL implementation is written in Golang and is fully compliant with the specification. We also have a process for developing new processors and adopting them into the specification.

<!-- ## Architecture and Components -->

## Architecture
Wiz ETL represents its data transformation tasks as a Directed Acyclic Graph (DAG). 

It can take a plurality of **data sources**. If the data is a streaming data source, it can be **refreshed**. <!-- ? -->

Each step in the DAG is a **processor**, which processes or transforms the data and can return an output. Each processor can have special configuration that modifies its behavior.

The power of Wiz ETL comes from the ability to do this at a large scale effectively, across a distributed cluster of nodes.

A **node** is a logical computing unit that can run a **processor**. It may be a physical or virtual machine, or a logical node in an orchestration environment like Kubernetes.

Most existing orchestration tools provide **schedulers** which determine where a given workload gets provisioned in the cluster. For the sake of simplicity, Wiz ETL will use those default schedulers by building on top of their existing APIs and working with the workload placement of each of the processors.

However, when designing a big data system of this scale, there are often much more complex considerations to take into account, including: network and data access topolgoies (e.g. bandwidth between nodes, and storage devices), the availability and location of specialized compute accelerators like GPUs, etc. High performance computing (HPC) systems take all of this into account. 

Therefore, in the future, it may be possible that Wiz ETL will:
1. provide a mechanism for specifying additional topology information and associating that with existing systems (physical machines/clusters or orchestrating systems) and
2. provide an algorithm for using APIs to schedule the proccessors to take advantage of that information

For now, this is left to the future

### Wiz ETL data philosophy
For many big data systems, data is a first class citizen and is often serialized, stored or manipulated as a first class citizen of the system itself. (https://flink.apache.org/flink-applications.html)

However, Wiz ETL fundamentally believes that the next era of data solutions will need multi-modal data that is transformed and eventually stored in a mult data-model database like FoundationDB.

Wiz ETL supports a wide array of data access mechanisms and databases for a variety of use-cases. However, the benefit of storing data in Wiz Data is evident.