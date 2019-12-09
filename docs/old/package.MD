## Describes the package file

```yaml
package: 
    name: mnist
    type: dataset
    version: "1.0.3"
    revised: "2/3/2019"
    # publishing organization
    publisher: "stanford"
    # maintainer usernames
    maintainers: 
      - yann.lecun
      - user.name
source:
    directory: "mnist_data"
    providers: 
        - gcp
        - aws
        - microsoft
        - url: "http://src.mnist.com/mnist.tar.gz"
transform:
    - download
    - tar
    - move
    - process "script.sh"
# tabular dataset verification standard
# hash of the hash of each row
```

```yaml
package:
  name: tensorflow
  type: library
  version: 2.0.0
```

```yaml
package:
    name: dnc
    type: module
    version: "1.2.0"
    publisher: "deepmind"
in:
source:
    git:
      repo: "github.com/deepmind/dnc"
      # path: 
      commit: abcabcabcabcabcabc
      tag: <automatic from version>
binaries:
    github_releases:
      repo: "github.com/deepmind/dnc"
deps:
    solver: SMT
    lockfile: wiz.lock
    packages:
        - python: "3.7"
        - tensorflow: "2.1"
        - mnist: "4.1"
out:
    executables: 
      - tensorboard
    linkables:
      - tensorflow
    python_packages:
      - tensorflow, tf
install: 
    - pip
    - wheel
    

```