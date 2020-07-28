![Test](https://github.com/avbru/image-previewer/workflows/Test/badge.svg?branch=master)
![integration-tests](https://github.com/avbru/image-previewer/workflows/integration-tests/badge.svg)
![golangci-lint](https://github.com/avbru/image-previewer/workflows/golangci-lint/badge.svg)

Specification: https://github.com/OtusGolang/final_project/blob/master/03-image-previewer.md

Router inspired by: https://github.com/kulti/otus_open_lesson/tree/v11052020/internal/router

Work in progress! Alpha stage...
To date only resizes images without a cache.

##Run, test:
```
make run
make test
make run_integration_test
```
Example images: http://localhost/api/samples/

## TODO
* Add config
* Add image cache
* Add some integration tests
* Add more unit tests

## Done TODOs
* Switch to latest golangci-lint
* Add integration tests to workflows