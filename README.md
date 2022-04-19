# Streamy

[![Godoc][doc-image]][doc-url] [![Release][release-image]][release-url] [![Build][build-image]][build-url]

A Golang library that provides functions for streams.

## Usage

See [streamy_test.go](streamy_test.go), [reader_test](reader_test), [progress_test](progress_test) and [have_test.go](have_test.go).

## Test

```shell
# Test everything:
make test

# For BDD development:
# It will open a new browser window. Make sure:
#   1. There is no errors on the terminal window.
#   2. There is no other open GoConvey page.
make test-ui

# Benchmarks
make test-benchmarks
```

## Release

```shell
# Update and commit CHANGELOG.md first (i.e. git add CHANGELOG.md && git commit -m "v1.0.0").
# Set GIT_TAG using semver (i.e. GIT_TAG=v1.0.0)
make release GIT_TAG=
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
[build-image]: https://github.com/devfacet/streamy/workflows/Test/badge.svg