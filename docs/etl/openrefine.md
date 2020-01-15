# OpenRefine inspiration

Each processor should have optional 

UI:

- using a sequence of multiple facets to filter data
- store all changes to data in the history of the dataset
- allow changing either the filter or the processor based on a previously configured step in history -- creates a history branch
- serialize filters and processors into graph

Mapping from openrefine data model to graph data model:

Filters are applied in sequential graph order to select data
Processor is applied on output and merged into original dataset, replacing specified rows with IDs
