# Streamy

[![Godoc][doc-image]][doc-url] [![Release][release-image]][release-url] [![Build][build-image]][build-url]

A Golang library that provides functions for streams.

## Usage

See [streamy_test.go](streamy_test.go), [reader_test.go](reader_test.go), [progress_test.go](progress_test.go) and [have_test.go](have_test.go).

## Test

```shell
# Run tests
make test

# Continuous testing
make test-ui

# Benchmarks
make test-benchmarks
```

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md)

## License

Licensed under The MIT License (MIT)  
For the full copyright and license information, please view the LICENSE.txt file.

[doc-url]: https://pkg.go.dev/github.com/devfacet/streamy
[doc-image]: https://pkg.go.dev/badge/github.com/devfacet/streamy

[release-url]: https://github.com/devfacet/streamy/releases/latest
[release-image]: https://img.shields.io/github/release/devfacet/streamy.svg?style=flat-square

[build-url]: https://github.com/devfacet/streamy/actions/workflows/test.yaml
[build-image]: https://github.com/devfacet/streamy/workflows/test.yaml/badge.svg
