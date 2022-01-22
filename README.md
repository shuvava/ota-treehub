# Treehub


[![GoDoc](https://godoc.org/github.com/shuvava/ota-treehub?status.svg)](http://godoc.org/github.com/shuvava/ota-treehub)
[![Build Status](https://travis-ci.com/shuvava/ota-treehub.svg?branch=master)](https://travis-ci.com/github/shuvava/ota-treehub)
[![Coverage Status](https://coveralls.io/repos/github/shuvava/ota-treehub/badge.svg?branch=master)](https://coveralls.io/github/shuvava/ota-treehub?branch=master)


Treehub implements an `ostree` repository storage for over the air updates. This project is migration of [ota-community-edition/treehub][1] to golang.

This project implements an HTTP api that `ostree` can use to natively pull objects and revisions to update an `ostree` repository.

An HTTP api is provided to receive `ostree` repository objects and refs from command line tools such as `garage-push`, included with
[sota-tools](https://github.com/advancedtelematic/sota-tools).

This repo if forked from [Advancedtelematic TreeHub](https://github.com/advancedtelematic/treehub) application rewritten on golang with full support of original API.

## Links

* [Advancedtelematic TreeHub](https://github.com/advancedtelematic/treehub)

[1]: https://github.com/advancedtelematic/treehub
