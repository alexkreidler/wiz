# Wiz Manager Design

Goals: to allow a lightweight embedded manager along with a powerful multi-user daemon.

Since the embedded manager will only kick of the processor by
1. configuring all processors so the task graph is translated into configuration
2. each processor will then wait for other processors to provide it with the downstream data. In fat, processors should have a link_test data mode that sends a simple message that verifies the network link between separate processors to make sure it is configured properly

The lightweight manager will then only refresh state data from the respective processors when it is invoked again by the user. It will only store minimal information on-disk which should not be modified. 

It should store:

- the immutable version of the task graph provided to it by the user
- the placement of each processor on each node. (e.g. the node_id if we use that or the run_id)
- the configuration which it has sent to each processor (for debugging)

Again, once the manager starts the graph it should be able to run to completion on its own, with the manager only needed to query its state.

The daemon processor can store all this state in a database and possibly have systems for RBAC/other authentication mechanisms  