# Treehub

Treehub implements an `ostree` repository storage for over the air
updates. This project is migration of [ota-community-edition/treehub][1] to golang.

This project implements an HTTP api that `ostree` can use to natively
pull objects and revisions to update an `ostree` repository.

An HTTP api is provided to receive `ostree` repository objects and
refs from command line tools such as `garage-push`, included with
[sota-tools](https://github.com/advancedtelematic/sota-tools).

## Links

* [Advancedtelematic TreeHub](https://github.com/advancedtelematic/treehub)

[1]: https://github.com/advancedtelematic/treehub
