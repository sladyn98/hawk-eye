
<p align="center">
    <img width="200px" src="img/logo.png">
</p>



[![License: GPL v3](https://img.shields.io/badge/License-GPLv3+-blue.svg)](http://www.gnu.org/licenses/gpl-3.0)
![Go](https://github.com/sladyn98/hawk-eye/workflows/Go/badge.svg?branch=master)

hawk-eye is a continuous integration status reporter built to watch over the Github CI statuses.

:construction: This is just  a proof of concept,and not fully stable. Expect dragons and unfinished business. :construction:

## Contribute
PRs accepted.

git clone git@github.com:sladyn98/hawk-eye.git
You can now run `make` to build the project, or `make install` to install the binary in $GOPATH/bin/.

## To-Do

a) Enable commands to get the CI status from Github Actions. eg: `hawk-eye getCIStatus` which would then return a one or zero enabling users to combine it with different steps like `hawk-eye getCIStatus && npm publish`.

b) Enable Travis CI status support.(Future Support)
