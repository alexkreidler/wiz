## Wiz Data Packages Specification

Wiz Data Packages are packages that publish data to be accessed. They have a set of unique requirements and features.

In general, data goes through a few steps before being useful. 

1. Data access (from HTTP, FTP, SQL, other databases, APIs, mirrors, etc)
2. Data reading/processing
3. Data transformations

There are many open source and other products available which are designed to achieve each of these goals, however, Wiz incorporates its umbrella projects to achieve these goals. 

Firstly, it is important to recognize that Wiz cannot and will not be able to solve all of these problems. We cannot cater to all usecases. However, we are committed to achieving the broadest and best feature set possible. 

To this end, we have a set of reasonable setttings to allow this. 

In general, all of these steps combined can be thought of as a DAG, but it is important to have a specification which represents this in an easy to understand way. 

The Wiz Package Manager is built on the Wiz Tasks Framework.

Specification: 

```yaml
package:
  name: mnist
  type: data
  version: "1.0.3" # this is associated with a specific hash of the data
  hash: "76773e85cc8fe22abce6fd7606c80a1f57c5c630"
  access: http
  access_options:
    url: http://yann.lecun.com/exdb/mnist/train-images-idx3-ubyte.gz

```

The `access_options` key represents a single **processor** in the Wiz Task Framework, and can be configured exactly as such

<!-- Rename `access` to `source`? -->

These keys are just syntactic sugar over a full WTF graph. In fact, some datasets come from multiple file/folder sources and need the full complexity of the graph to expose their data. 

This is shown below: 


```yaml
package:
  name: mnist
  type: data
  version: "1.0.3" # this is associated with a specific hash of the data
  hash: "76773e85cc8fe22abce6fd7606c80a1f57c5c630"
  sources:
    - access: http
      access_options:
        url: http://yann.lecun.com/exdb/mnist/train-images-idx3-ubyte.gz
    - access: http
      access_options:
        url: http://yann.lecun.com/exdb/mnist/train-labels-idx1-ubyte.gz
```

As seen above, this is quite repetitive, especially for many files.

t10k-images-idx3-ubyte.gz:   test set images (1648877 bytes)
t10k-labels-idx1-ubyte.gz:   test set labels (4542 bytes) 

```yaml
package:
  name: mnist
  type: data
  version: "1.0.3" # this is associated with a specific hash of the data
  hash: "76773e85cc8fe22abce6fd7606c80a1f57c5c630"
  sources:
    access: http
    access_options:
        url: http://yann.lecun.com/exdb/mnist/
        files: 
            - t10k-images-idx3-ubyte.gz
            - t10k-labels-idx1-ubyte.gz
```

However, the Wiz logical `type: data` package is effectively just a folder. Thus, if no further processing were applied to these datasets, then the structure would look like

```
.
├── train-images-idx3-ubyte.gz
├── train-labels-idx1-ubyte.gz
```

Thus, we have more advanced options:

```yaml
package:
  name: mnist
  type: data
  version: "1.0.3" # this is associated with a specific hash of the data
  hash: "76773e85cc8fe22abce6fd7606c80a1f57c5c630"
  sources:
    - access: http
      access_options:
        url: http://yann.lecun.com/exdb/mnist/train-images-idx3-ubyte.gz
      compression: .gz
      format: custom | binary
    - access: http
      access_options:
        url: http://yann.lecun.com/exdb/mnist/train-labels-idx1-ubyte.gz
```

This effectively represents the following graph

- Download
- Decompress
- Provide binary file

which are specified as follows:


```yaml
package:
  name: mnist
  type: data
  version: "1.0.3" # this is associated with a specific hash of the data
  hash: "76773e85cc8fe22abce6fd7606c80a1f57c5c630"
  sources:
    - access: http
      access_options:
        url: http://yann.lecun.com/exdb/mnist/train-images-idx3-ubyte.gz
      processors: 
        - type: decompress
          options:
            format: .gz
        - type: move
        - type: output
          options:
            format: custom | binary
    - access: http
      access_options:
        url: http://yann.lecun.com/exdb/mnist/train-labels-idx1-ubyte.gz
```

For some other formats including CSV, TSV, JSON, etc, which can be parsed and loaded into a database, this is enabled as an option through the CLI


```yaml
name: bls
description: Bureau of Labor Statistics
source:
  type:
    proprietary_db: LABSTAT
    format: tsv
  access: http
  url: https://download.bls.gov/pub/time.series/
```

```bash
wiz get bls --db=db-config.yaml --local=/data
```



```yaml
outputs:
    - type: sql
      options:
        uri: sql://mysql
        table_name: bls-data
    - type: foundationdb
      options:
        config_file: docker.config
        namespace_name: bls-data
    - type: local
      options: 
        dir: /data/bla
```


This `output.yaml` file effectively makes up the output processors of the DAG, and will be run in parallel by default as the data is all the same.

Depending on the output steps and the steps defined in the package manifest, Wiz will optimize the graph in various ways. Some processors, like HTTP and FTP downloads of folders can be parallelized so multiple files can be downloaded at once. If the data does not need to be output in the `type: local` format, then Wiz will only write a temporary file as it is sent to the database.

Some datasets are streaming which means that their outputs are not explicit and static. For example, running `wiz get cnn-rss --output=foundationdb` will result in a long-running task that ingests CNN RSS feeds and inputs it into the databse.

However, since all `wiz` cli task related commands return instantly and then provide references to running tasks on the cluster, except in the case of the **local** node, it will be immediately abstracted onto the cluster.
