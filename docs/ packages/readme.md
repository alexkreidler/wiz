# Wiz packages

The wiz package manager's goals are to provide a speedy development experience for everyone. The idea is to easily compartmentalize all dependencies irrespective of the language that they are in. We also aim to eliminate dependency hell by making all packages 100% functional. This also provides additional opportunities for caching, etc.

Functional requirements

Ideas:

The package can have a manifest file, which has 3 key sections

```yaml
metadata:
    name:
    author:
    etc:
dependencyCrafter:
    - name: test-package
      version: ~0.1.0 # A semver version range
dependencies:
    - name: test-package
      version: 0.1.0 # a fully resolved semver version that only represents a single version of a package
```

The dependencyCrafter key, or similar, it solely there to help the user design the versions of the package that it requires. The user should check the final outputed dependency map to make sure it is generated as expected. 

This is different than most package manager which allow packages to specify version ranges and then use a special version solver algorithm which varies. NPM for example simply fetches the most recent version of any given package recursively, and then has an algorithm that can hoist packages up so that the replication is not complete. 

Conda uses an SMT solver algorithm to try and find packages that are compatible with all other packages. This can result in downgrading existing packages to make them work with newly installed ones.

Although all package managers try to make their algorithms as transparent, explainable, etc as possible, it is still a source of a lot of confusion and unexpected behavior across the board. Thus Wiz borrows from Nix in terms of using 100% functional and deterministic packages where each package must specify the exact downstream version that it requires.

This can result in a certain amount of duplication on disk, but generally this is not a problem for developer machines and the benefits significantly outweigh the costs.

There are many other benefits to Functional PMs as well. 