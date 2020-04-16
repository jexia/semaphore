# Flow manager

Flow managers manage nodes within a given flow. Each node represents an action to be executed (ex: request, rollback).
Nodes are executed concurrently from one another. Dependencies are based on references within the given node or if a node dependency is defined.
A manager keeps track of all the processes being executed and tracks all the nodes which have been called.
If an error is thrown inside one of the processes during the execution of the flow stopped and a rollback initiated.

## Branches

Nodes are executed concurrently from one another.
When a node is executed a check is performed to check whether the dependencies have been met.
Only if all of the dependencies have been met is the node executed.

```
+------------+
|            |
|    Node    +------------+
|            |            |
+------+-----+            |
       |                  |
       |                  |
+------v-----+     +------v-----+
|            |     |            |
|    Node    |     |    Node    |
|            |     |            |
+------+-----+     +------+-----+
       |                  |
       |                  |
+------v-----+            |
|            |            |
|    Node    <------------+
|            |
+------------+
```