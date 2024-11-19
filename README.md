# achilles-sdk-api

API types consumed by [`reddit/achilles-sdk`](https://github.com/reddit/achilles-sdk).

This repo should minimize dependencies on external Go modules and if it must import external modules, it _must not_ use
any runtime logic.

This is to ensure that the structs, interfaces, and types exported by this module can be used in a variety of consuming
projects without causing dependency conflicts and thus forcing particular versions of commonly imported Kubernetes modules
(e.g. `github.com/kubernetes-sigs/controller-runtime`, `github.com/kubernetes/client-go`).
