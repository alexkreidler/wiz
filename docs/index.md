# wiz

| Name        | Role                                                                                                                     |
| ----------- | ------------------------------------------------------------------------------------------------------------------------ |
| Wiz Core    | An extensible, functional, distributed package manager with a focus on Machine Learning projects                         |
| Wiz Tasks   | To run a DAG of generic tasks. Wiz ETL is built on this framework                                                        |
| Wiz Sources | To continously fetch data from updating sources                                                                          |
| Wiz Data    | To store data in whatever format available, either as unstrucutured/raw or in a database. Likely built on FoundationDB   |
| Wiz ETL     | To transform and harmonize the data across sources, preprocess, etc                                                      |
| Wiz Viz     | To perform effective server-side visualization of large data sets using a slightly modified version of the Vega standard |

As time goes on I'll add components which there is a need for

## Sources
Sources provide data which can be **refreshed**.

Declarative format -- in the package manifest
Multiple source plugins:
- HTTP scraping
- JSON/API requests
- Download via multiple protocols: 
  - HTTP
  - FTP
  - Git
- RSS
- WebSockets?

Specify a frequency at which to refresh the data.

Can also be used to input data which is not only being generated in the future but is from further in the past.
E.g. we can run the same web scraper on Archive.org pages to get stories from further back in time.
Sources would allow us to start training with an existing dataset and add to it.

Specify a **recording** or **persisting** protocol for transient sources that will go away -- e.g. persist the history of the RSS feed from when we listen to it

## Wiz Data

Data lake capabilities - unstructured and structured

Enable tracking in Git?? of all our data assets from creation to analysis

unstructured:
- binary data
- raw documents, pdfs, word files, images
- unformatted csvs,  json, etc that are not standard or cannot be parsed properly

structured:
- stata, spss, statistical files
- csv, tsv, json, yaml, etc
- parquet

All structured data should immediately be parsed and stored in database.
This then means wiz etl will need to work *in* the db


Database:
- multi-model: relational-ish, graph, tabular, key-value, timeseries, geographic
- FoundationBD achieves this
- Wiz ETL should have an interface for the structure of data storage

foundationdb allows for mixing and matching data strucutes
e.g. a time series graph db
or a tabular geographic db

optimized by the first accessible keys

interface should also have a system for indexes or running aggreates, averages, or summary stats of various things
These should be updated by the wiz wrapper around foundationdb (maybe called tachyonDB)

Benefits of storing earlier vs later after more processing
Earlier
- wiz ETL interfaces with one system
- already there, easier to transform into different data structures
- can hold the history of transformations directly in the DB (recorded through wrapper layer)

Later
- easier? to write basic python scripts etc
- simpler for smaller datasets?^

Also since FDB has limits, files etc must be stored in S3 or other object storage


All real data models are like graphs or at least Entity, Properties, Instances, References:
- relational: entities are each table, properties are columns, rows are instances, foreign keys are links to entities (references)
- document: entities are documents, fields are properties, object fields are nested entities/references
- graph: nodes are (instances of) entities, predicates and edge nodes are properties, predicates to nodes are references

They mainly differ only in the
1. number/relationship 
2. acces of data
of Entity, Properties, Instances, References

| Name       | Number of each         | Access of data                                                                                                       |
| ---------- | ---------------------- | -------------------------------------------------------------------------------------------------------------------- |
| Relational | n/a, many, n/a, few      | Mainly individual records or summaries of properties                                                                 |
| Document   | many, many, less, less | Multi-modal data, flexible structure. Either nesting or references/IDs                                               |
| Graph      | n/a, n/a, n/a, many       | Mainly interested in references between nodes: e.g. distance. Comparing different types of nodes with same relations |

The truth is we think way to much about data in terms of the way it is stored and its performance, vs how we actually want to access the data.

Build a system for multiple *views of the data* in each way???
<!-- The UniDB, ThorDB, olympusDB, enikodb, uniqdb -->
implements all three APIs: SQL, MongoApi, n-triples

Deal with:
- indexes
- caching
- sharding/partitioning/distributing
- search?
- algorithms:
  - statistical summaries/functions
  - relations operations e.g. graph algorithms
  - concurrent/atomic operations, transactions, many users
  - machine learning algorithms:
    - most just need batched access to some individual records
    - some need stat funcs.