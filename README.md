![LINT](https://github.com/avbru/image-previewer/workflows/lint/badge.svg)
![TEST](https://github.com/avbru/image-previewer/workflows/test/badge.svg)
[![Coverage](https://coveralls.io/repos/github/avbru/image-previewer/badge.svg?branch=master)](https://coveralls.io/github/avbru/image-previewer?branch=master)

Specification: https://github.com/OtusGolang/final_project/blob/master/03-image-previewer.md

Router inspired by: https://github.com/kulti/otus_open_lesson/tree/v11052020/internal/router

Work in progress! Alpha stage...
To date only resizes images without a cache.

##Run, test:
```
make run
make run_integration_test
```
Example images: http://localhost/api/samples/

## TODO
* Add config
* Add image cache
* Add some integration tests
* Add more unit tests
* Switch to latest golangci-lint
* Add integration tests to workflows