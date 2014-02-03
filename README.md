goat [![Build Status](https://travis-ci.org/mdlayher/goat.png?branch=master)](https://travis-ci.org/mdlayher/goat)
====

goat: __Go__ __A__wesome __T__racker.  BitTorrent tracker implementation, written in Go.  MIT Licensed.

Full documentation for goat may be found on [GoDoc](http://godoc.org/github.com/mdlayher/goat).

For convenience, these links are provided to quickly jump to specific areas of documentation:

- [API](http://godoc.org/github.com/mdlayher/goat#hdr-API)
- [Configuration](http://godoc.org/github.com/mdlayher/goat#hdr-Configuration)
- [Listeners](http://godoc.org/github.com/mdlayher/goat#hdr-Listeners)

Installation
============

To download, build, and install goat, simply run:

`go get github.com/mdlayher/goat`

If using MySQL, the SQL schema files for goat can be found in [`res/mysql/`](https://github.com/mdlayher/goat/tree/master/res/mysql).
The database tables must be created manually before goat will run.

To build goat for use with ql, you can run:

`go get -tags='ql' github.com/mdlayher/goat`

The ql schema files and a tool to build the database can be found in [`res/ql`](https://github.com/mdlayher/goat/tree/master/res/mysql).
The ql database will be automatically copied from `res/ql/goat.db` to `~/.config/goat/goat.db`.

Contributing
============

If you'd like to contribute patches to goat, we'd love to have your help!  There are a few requirements
which we ask you follow when contributing:

- Ensure your code is properly formatted, linted, and error checked.  Running `make fmt` will call `go fmt`,
`golint`, and `errcheck` to handle this automatically.
- Ensure your code passes all tests, by running `make test` to verify each package works properly.
- Document your code.  Heavily.  Make it very clear exactly what the program is doing, and why.

Special thanks to the following for their help with the project:

- [NickPresta](https://github.com/NickPresta) - general guidance, code linting, cleanup
- [ChimeraCoder](https://github.com/ChimeraCoder) - project structure, remote deployment fixes
- [sdgoij](https://github.com/sdgoij) - ql storage backend, database abstraction
- [toqueteos](https://github.com/toqueteos) - constant time string comparison
