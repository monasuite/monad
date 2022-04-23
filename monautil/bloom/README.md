bloom
=====

[![Build Status](http://img.shields.io/travis/monasuite/monautil.svg)](https://travis-ci.org/monasuite/monautil)
[![ISC License](http://img.shields.io/badge/license-ISC-blue.svg)](http://copyfree.org)
[![GoDoc](http://img.shields.io/badge/godoc-reference-blue.svg)](https://pkg.go.dev/github.com/monasuite/monautil/bloom)


Package bloom provides an API for dealing with monacoin-specific bloom filters.

A comprehensive suite of tests is provided to ensure proper functionality.  See
`test_coverage.txt` for the gocov coverage report.  Alternatively, if you are
running a POSIX OS, you can run the `cov_report.sh` script for a real-time
report.

## Installation and Updating

```bash
$ go get -u github.com/monasuite/monautil/bloom
```

## Examples

=======
* [NewFilter Example]
  (http://godoc.org/github.com/monasuite/monautil/bloom#example-NewFilter)  
  Demonstrates how to create a new bloom filter, add a transaction hash to it,
  and check if the filter matches the transaction.

## License

Package bloom is licensed under the [copyfree](http://copyfree.org) ISC
License.
