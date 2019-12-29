# The Wiz Project

Wiz is a package manager for AI and ML projects and a set of associated libraries for dealing with common problems.

This is an ambitious project, so I've outlined a set of goals for v0.1.0:

* [x] build out the core Wiz Pipelines API (see [here](../tasks/tasks.go))
* [x] build the Wiz Processor API and implement the following basic builtin processors:
* [x] Git download/checkout
* [ ] Build the go-getter Processor using the go-getter package. This satisfies all basic fetching reqs and more. Main issue is that not all config options are exposed/supported, but it is really good for quick fetching
* [ ] Build out Wiz Manager and setting up the initial configuration of the processors.
* [ ] Build the auto-forwarding capabilities of the processors so the manager doesn't have to be running constantly
* [ ] SQL and
* [ ] FoundationDB loaders
* [ ] Build a basic Wiz Packages API that incorporates Wiz Pipelines for data packages