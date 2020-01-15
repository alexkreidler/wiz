# Wiz ETL/Data API

General idea: apache arrow + (WASM) on frontend --> HTTP, WebSockets --> apache arrow in-mem DB backend with data faceting and modifying (openrefine functionality) --> FoundationDB key-value store

eventually, may need to build own highly available but with better consistency/HA/clustering Ordered KV Store