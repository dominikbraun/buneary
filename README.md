<h1>buneary
<a href="https://circleci.com/gh/verless/verless"><img src="https://circleci.com/gh/verless/verless.svg?style=shield"></a>
<a href="https://www.codefactor.io/repository/github/verless/verless"><img src="https://www.codefactor.io/repository/github/verless/verless/badge" /></a>
<a href="https://github.com/verless/verless/releases"><img src="https://img.shields.io/github/v/release/verless/verless?sort=semver"></a>
<a href="LICENSE"><img src="https://img.shields.io/badge/license-Apache--2.0-brightgreen"></a>
</h1>

`buneary`, pronounced _bun-ear-y_, is an easy-to-use RabbitMQ command line client for managing exchanges, managing
queues and publishing messages to exchanges.

<p>
<br>
<img src="logo.png" alt="buneary">
<br>
<br>
</p>

---

## Installation

Download the [latest release](https://github.com/dominikbraun/buneary/releases) for your desired platform.

### macOS/Linux

Extract the downloaded binary into a directory like `/usr/local/bin`. Make sure the directory is in `PATH`.

### Windows

Create a directory like `C:\Program Files\buneary` and extract the executable into that directory.
[Add the directory to `Path`](https://www.computerhope.com/issues/ch000549.htm).

### Docker

Just run `docker container run dominikbraun/buneary` followed by the actual command you want to execute.

```
$ docker container run dominikbraun/buneary version
```

## Usage

### Create an exchange

```
$ buneary create exchange <ADDRESS> <NAME> <TYPE>
```

...

```
Usage:
  buneary create exchange <ADDRESS> <NAME> <TYPE> [flags]

Flags:
      --auto-delete   make the exchange auto-deleted
      --durable       make the exchange durable
  -h, --help          help for exchange
      --internal      make the exchange internal

Global Flags:
  -p, --password string   the password to authenticate with
  -u, --user string       the username to connect with
```


