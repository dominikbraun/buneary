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

## Contents

* [Installation](#installation)
    * [macOS/Linux](#macoslinux)
    * [Windows](#windows)
    * [Docker](#docker)
* [Usage](#usage)
    * [Create an exchange](#create-an-exchange)
    * [Create a queue](#create-a-queue)

## Installation

### macOS/Linux

Download the [latest release](https://github.com/dominikbraun/buneary/releases) for your platform. Extract the
downloaded binary into a directory like `/usr/local/bin`. Make sure the directory is in `PATH`.

### Windows

Download the [latest release](https://github.com/dominikbraun/buneary/releases), create a directory like
`C:\Program Files\buneary` and extract the executable into that directory.
[Add the directory to `Path`](https://www.computerhope.com/issues/ch000549.htm).

### Docker

Just append the actual `buneary` command you want to run after the image name.

Because `buneary` needs to dial the RabbitMQ server, the Docker container needs to be in the same network as the
RabbitMQ server. For example, if the server is running on your local machine, you could run a command as follows:

```
$ docker container run dominikbraun/buneary --network=host publish ... 
```

## Usage

### Create an exchange

**Syntax:**

```
$ buneary create exchange <ADDRESS> <NAME> <TYPE> [flags]
```

**Arguments:**

|Argument|Description|
|-|-|
|`ADDRESS`|The RabbitMQ server address. If no port is specified, `5672` is used.|
|`NAME`|The desired name of the new exchange.|
|`TYPE`|The exchange type. Has to be one of `direct`, `headers`, `fanout` and `topic`.|

**Flags:**

|Flag|Short|Description|
|-|-|-|
|`--user`|`-u`|The username to connect with. If not specified, you will be asked for it.|
|`--password`|`-p`|The password to authenticate with. If not specified, you will be asked for it.|
|`--auto-delete`||Automatically delete the exchange once there are no bindings left.|
|`--durable`||Make the exchange persistent, surviving server restarts.|
|`--internal`||Make the exchange internal.|

**Example:**

Create a direct exchange called `my-exchange` on a RabbitMQ server running on the local machine.

```
$ buneary create exchange localhost my-exchange direct
```

### Create a queue

**Syntax:**

```
$ buneary create queue <ADDRESS> <NAME> <TYPE> [flags]
```

**Arguments:**

|Argument|Description|
|-|-|
|`ADDRESS`|The RabbitMQ server address. If no port is specified, `5672` is used.|
|`NAME`|The desired name of the new queue.|
|`TYPE`|The queue type. Has to be one of `classic` and `quorum`.|

**Flags:**

|Flag|Short|Description|
|-|-|-|
|`--user`|`-u`|The username to connect with. If not specified, you will be asked for it.|
|`--password`|`-p`|The password to authenticate with. If not specified, you will be asked for it.|
|`--auto-delete`||Automatically delete the queue once there are no consumers left.|
|`--durable`||Make the queue persistent, surviving server restarts.|

**Example:**

Create a classic queue called `my-queue` on a RabbitMQ server running on the local machine.

```
$ buneary create queue localhost my-queue classic
```
