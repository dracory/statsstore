# statsstore

<a href="https://gitpod.io/#https://github.com/gouniverse/statsstore" style="float:right:"><img src="https://gitpod.io/button/open-in-gitpod.svg" alt="Open in Gitpod" loading="lazy"></a>

[![Tests Status](https://github.com/gouniverse/statsstore/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/gouniverse/statsstore/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/gouniverse/statsstore)](https://goreportcard.com/report/github.com/gouniverse/statsstore)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/gouniverse/statsstore)](https://pkg.go.dev/github.com/gouniverse/statsstore)

Vault - a secure value storage (data-at-rest) implementation for Go.

## License

This project is licensed under the GNU General Public License version 3 (GPL-3.0). You can find a copy of the license at https://www.gnu.org/licenses/gpl-3.0.en.html

For commercial use, please use my [contact page](https://lesichkov.co.uk/contact) to obtain a commercial license.

## Installation
```
go get -u github.com/gouniverse/statsstore
```

## Setup

```golang
store, err := NewStore(NewStoreOptions{
	VisitorTableName:     "stats_visitor",
	DB:                 databaseInstance,
	AutomigrateEnabled: true,
})

```
