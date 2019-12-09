# Wiz Data

The data flowing between components of Wiz is one of the most important components of the framework, and arguably the most difficult to make decisions for.

Unlike other frameworks like Apache Beam, Wiz does not provide explicit structures that the data must fit into, and does not provide any additional features on top of the data except for Provenance, which provides detailed auditing data about every record and file in the system.

Wiz has two types of data
1. records
2. folders/files

Records are processable logical records, which originate from files in any the supported formats including:
- csv, tsv, etc
- json
- avro
- parquet
- protobuf

The record-level features allow Wiz to parallelize operations effectively based on heuristics about each processor.

Records do not specifically mean only row-based data. They can also represent unstructured documents and columns of a dataset. 

