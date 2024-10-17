<img src="https://raw.githubusercontent.com/google/mtail/main/logo.png" alt="mtail" title="mtail" align="right" width="140">

# mtail - extract internal monitoring data from application logs for collection into a timeseries database

[![ci](https://github.com/jaqx0r/mtail/workflows/CI/badge.svg)](https://github.com/jaqx0r/mtail/actions?query=workflow%3ACI+branch%3main)
[![GoDoc](https://godoc.org/github.com/jaqx0r/mtail?status.png)](http://godoc.org/github.com/jaqx0r/mtail)
[![Go Report Card](https://goreportcard.com/badge/github.com/jaqx0r/mtail)](https://goreportcard.com/report/github.com/jaqx0r/mtail)
[![OSS-Fuzz](https://oss-fuzz-build-logs.storage.googleapis.com/badges/mtail.svg)](https://bugs.chromium.org/p/oss-fuzz/issues/list?sort=-opened&can=1&q=proj:mtail)
[![codecov](https://codecov.io/gh/jaqx0r/mtail/branch/main/graph/badge.svg)](https://codecov.io/gh/jaqx0r/mtail)

`mtail` is a tool for extracting metrics from application logs to be exported
into a timeseries database or timeseries calculator for alerting and
dashboarding.

It fills a monitoring niche by being the glue between applications that do not
export their own internal state (other than via logs) and existing monitoring
systems, such that system operators do not need to patch those applications to
instrument them or writing custom extraction code for every such application.

The extraction is controlled by [mtail programs](https://jaqx0r.github.io/mtail/Programming-Guide)
which define patterns and actions:

    # simple line counter
    counter lines_total
    /$/ {
      lines_total++
    }

Metrics are exported for scraping by a collector as JSON or Prometheus format
over HTTP, or can be periodically sent to a collectd, StatsD, or Graphite
collector socket.

Read the [programming guide](https://jaqx0r.github.io/mtail/Programming-Guide) if you want to learn how
to write mtail programs.

Ask general questions on the users mailing list: https://groups.google.com/g/mtail-users

## Installation

There are various ways of installing **mtail**.

### Precompiled binaries

Precompiled binaries for released versions are available in the
[Releases page](https://github.com/jaqx0r/mtail/releases) on Github. Using the
latest production release binary is the recommended way of installing **mtail**.

Windows, OSX and Linux binaries are available.

### Building from source

The simplest way to get `mtail` is to `go get` it directly.

`go get github.com/jaqx0r/mtail/cmd/mtail`

This assumes you have a working Go environment with a recent Go version.  Usually mtail is tested to work with the last two minor versions  (e.g. Go 1.12 and Go 1.11).

If you want to fetch everything, you need to turn on Go Modules to succeed because of the way Go Modules have changed the way go get treats source trees with no Go code at the top level.

```
GO111MODULE=on go get -u github.com/jaqx0r/mtail
cd $GOPATH/src/github.com/jaqx0r/mtail
make install
```

If you develop the compiler you will need some additional tools
like `goyacc` to be able to rebuild the parser.

See the [Build instructions](https://jaqx0r.github.io/mtail/Building) for more details.

A `Dockerfile` is included in this repository for local development as an
alternative to installing Go in your environment, and takes care of all the
build dependency installation, if you don't care for that.


## Deployment

`mtail` works best when paired with a timeseries-based calculator and
alerting tool, like [Prometheus](http://prometheus.io).

> So what you do is you take the metrics from the log files and
> you bring them down to the monitoring system?

[It deals with the instrumentation so the engineers don't have
to!](http://www.imdb.com/title/tt0151804/quotes/?item=qt0386890)  It has the
extraction skills!  It is good at dealing with log files!!

Learn more about [interoperability with other tools](https://jaqx0r.github.io/mtail/Interoperability)

## Read More

Full documentation at https://jaqx0r.github.io/mtail/

Read more about writing `mtail` programs:

* [Programming Guide](https://jaqx0r.github.io/mtail/Programming-Guide)
* [Language Reference](https://jaqx0r.github.io/mtail/Language)
* [Metrics](https://jaqx0r.github.io/mtail/Metrics)
* [Managing internal state](https://jaqx0r.github.io/mtail/state)
* [Testing your programs](https://jaqx0r.github.io/mtail/Testing)

Read more about hacking on `mtail`

* [Building from source](https://jaqx0r.github.io/mtail/Building)
* [Contributing](CONTRIBUTING.md)
* [Style](https://jaqx0r.github.io/mtail/style)

Read more about deploying `mtail` and your programs in a monitoring environment

* [Deploying](https://jaqx0r.github.io/mtail/Deploying)
* [Interoperability](https://jaqx0r.github.io/mtail/Interoperability) with other systems
* [Troubleshooting](https://jaqx0r.github.io/mtail/Troubleshooting)
* [FAQ](https://jaqx0r.github.io/mtail/faq)


## Getting more help and reporting defects

If you have any questions, please use the [GitHub Discussions Q&A](https://github.com/jaqx0r/mtail/discussions/new?category=q-a).

We also have an email list : https://groups.google.com/forum/#!forum/mtail-users

For any defects please [file a new issue](https://github.com/jaqx0r/mtail/issues/new).
