# Stats Store

<a href="https://gitpod.io/#https://github.com/dracory/statsstore" style="float:right:"><img src="https://gitpod.io/button/open-in-gitpod.svg" alt="Open in Gitpod" loading="lazy"></a>

[![Tests Status](https://github.com/dracory/statsstore/actions/workflows/tests.yml/badge.svg?branch=main)](https://github.com/dracory/statsstore/actions/workflows/tests.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/dracory/statsstore)](https://goreportcard.com/report/github.com/dracory/statsstore)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/dracory/statsstore)](https://pkg.go.dev/github.com/dracory/statsstore)

Stats Store - a visitor stats storage implementation for Go.

## License
This project is licensed under the GNU Affero General Public License v3.0 (AGPL-3.0). You can find a copy of the license at [https://www.gnu.org/licenses/agpl-3.0.en.html](https://www.gnu.org/licenses/agpl-3.0.txt)

For commercial use, please use my [contact page](https://lesichkov.co.uk/contact) to obtain a commercial license.

## Installation
```
go get -u github.com/dracory/statsstore
```

## Setup

```golang
store, err := NewStore(NewStoreOptions{
	VisitorTableName:     "stats_visitor",
	DB:                 databaseInstance,
	AutomigrateEnabled: true,
})

```
