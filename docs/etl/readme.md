# Wiz ETL Framework

an end-to-end framework built on top of the Wiz Tasks Framework to assist with data harmonization and processing.

Ideas: 
How to deal with the Library  of Congress scraper: has both configuration (e.g. which API endpoints to use) and logic (a.k.a should use this detail level for these races not others, and should loop this many times).
How to balance declarative configuration (easily allows customization) vs code/logic which is more flexible and powerful. Most templating engines provide a level of logic builtin to the configuration, but its crappy and not worth it, learning a DSL.

Goal: put as much in configuration as possible, then have Hooks or something to call into code in any language and determine the other parts of configuration or dynamic config/next steps for the processors.

UI: 
The ETL UI should be solely focused on manipulating the data, not the underlying internals of the Task framework.

## Data Model Views
4 Tabs/Views, one for each data model:
- Relational/Tabular
- Structured --> e.g. how to transform structured data like JSON into the most sensible KV system. actually seems very simple just do: key: value, and key.subkey: value. But allow for building indexes,etc
- Graph
- Custom data models: timeseries, combinations, etc. Have a online repository of example data models which are their format, "[idx1, idx2, idx3, etc] = val2:val2:etc" and documentation on patterns to use it with and access performance etc. Then user can just plug in the keys and vars to the model based on what they need.

Example custom: timeseries
Will have a UI to choose the detail/resolution of the time. E.g. a simple single Unix timestamp, or an encoded format like "2017-3-19" or a set of [year, month, day, hour, minute, second, etc]
Then that generates the beginning of the KV format, the user can add what else they need: e.g. a tabular format after that

Think about composing multiple formats.

Each view is the **main view** that represents data that is meant to be accessed or primarily matches that data model


## Data Sources
Should have a menu like File->Open to open Wiz Package Manager spec files or an Add button that adds packages by name (fetching from a configurable registry)

Not thought about realtime features yet.

## Raw Data viewing

Goal: we want high-performance viewing of the raw data, which could be text files multiple Gigs in size. UX like glogg

We also want to be able to quickly preview it or sections of it in a data format provided by the data models. E.g. preview a 3GB CSV file as 1. raw and 2. in a tabular format


For more specific formats, like health records or images/scans, we prolly can't view them effectively here, so other tools should be used. Think about adding generic API to view different types of data. 

The Wiz Data View API is a performant API for a server to serve a raw text file and for a client to view different parts of it on-demand. Basically just one endpoint /view/:fileID/:locationSpecifier

where locationSpecifier is a range of bytes or lines.

TODO: think about performance of this approach --  would streaming via WebSockets be better? how to deal with slow scrolling, should keep old content, maybe overscan twice to allow good slow scrolling.

In terms of viewing in a data model format, should the parser be implemented on  Client or Server side. To be snappy, prolly server side, but then needs to be encoded anyway to be sent to client. Do tests to figure out