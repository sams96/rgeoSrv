# rgeoSrv
[![](https://img.shields.io/github/workflow/status/sams96/rgeoSrv/continuous-integration?style=for-the-badge)](https://github.com/sams96/rgeoSrv/actions?query=workflow%3Acontinuous-integration)
[![](https://goreportcard.com/badge/github.com/sams96/rgeoSrv?style=for-the-badge)](https://goreportcard.com/report/github.com/sams96/rgeoSrv)
[![Release](https://img.shields.io/github/tag/sams96/rgeoSrv.svg?label=release&color=24B898&logo=github&style=for-the-badge)](https://github.com/sams96/rgeoSrv/releases/latest)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=for-the-badge)](https://pkg.go.dev/github.com/sams96/rgeoSrv)

rgeoSrv wraps the package rgeo into a reverse geocoding microservice.

See [github.com/sams96/rgeo](https://github.com/sams96/rgeo) for more
information on rgeo.

### Installation

    go get github.com/sams96/rgeoSrv/..

or,

    docker pull docker.pkg.github.com/sams96/rgeosrv/rgeosrv

### Usage

    rgeoSrv -addr localhost:8080

or,

	docker run -p 8080:8080 docker.pkg.github.com/sams96/rgeosrv/rgeosrv and

then:

    curl "localhost:8080/query?0&52"

will yield:

    {"country":"United Kingdom","country_long":"United Kingdom of Great Britain and Northern Ireland","country_code_2":"GB","country_code_3":"GBR","continent":"Europe","region":"Europe","subregion":"Northern Europe"}
