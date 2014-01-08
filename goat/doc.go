/*
Package goat provides an implementation of a BitTorrent tracker, written in Go.

Installation

goat can be built using Go 1.1+. It can be downloaded, built, and installed,
simply by running 'go get github.com/mdlayher/goat'.

In addition, goat depends on a MySQL server for data storage.  After creating a
database and user for goat, its database schema may be imported from the SQL
files located in 'res/'.  goat will not run unless MySQL is installed, and a
database and user are properly configured for its use.
*/
package goat
