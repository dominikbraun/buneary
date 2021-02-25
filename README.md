<h1>buneary
<a href="https://circleci.com/gh/dominikbraun/buneary"><img src="https://circleci.com/gh/dominikbraun/buneary.svg?style=shield"></a>
<a href="https://www.codefactor.io/repository/github/dominikbraun/buneary"><img src="https://www.codefactor.io/repository/github/dominikbraun/buneary/badge" /></a>
<a href="https://github.com/dominikbraun/buneary/releases"><img src="https://img.shields.io/github/v/release/dominikbraun/buneary?sort=semver"></a>
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

* [Example](#example)
* [Installation](#installation)
    * [macOS/Linux](#macoslinux)
    * [Windows](#windows)
    * [Docker](#docker)
* [Usage](#usage)
    * [Create an exchange](#create-an-exchange)
    * [Create a queue](#create-a-queue)
    * [Create a binding](#create-a-binding)
    * [Get all exchanges](#get-all-exchanges)
    * [Get an exchange](#get-an-exchange)
    * [Get all queues](#get-all-queues)
    * [Get a queue](#get-a-queue)
    * [Get all bindings](#get-all-bindings)
    * [Get a binding](#get-a-binding)
    * [Get messages in a queue](#get-messages-in-a-queue)
    * [Publish a message](#publish-a-message)
    * [Delete an exchange](#delete-an-exchange)
    * [Delete a queue](#delete-a-queue)
* [Credits](#credits)
    
## Example

In the following example, a message `Hello!` is published and sent to an exchange called `my-exchange`. The RabbitMQ
server is running on the local machine, and we'll use a routing key called `my-routing-key` for the message.

```
$ buneary publish localhost my-exchange my-routing-key "Hello!"
```

Since the RabbitMQ server listens to the default port, the port can be omitted here. The above command will prompt you
to type in the username and password, but you could do this using command options as well.

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
$ docker container run --network=host dominikbraun/buneary version
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
|`ADDRESS`|The RabbitMQ HTTP API address. If no port is specified, `15672` is used.|
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
|`ADDRESS`|The RabbitMQ HTTP API address. If no port is specified, `15672` is used.|
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

### Create a binding

**Syntax:**

```
$ buneary create binding <ADDRESS> <NAME> <TARGET> <BINDING KEY> [flags]
```

**Arguments:**

|Argument|Description|
|-|-|
|`ADDRESS`|The RabbitMQ HTTP API address. If no port is specified, `15672` is used.|
|`NAME`|The desired name of the new binding.|
|`TARGET`|The name of the target queue or exchange. If it is an exchange, use `--to-exchange`.|
|`BINDING KEY`|The binding key.|

**Flags:**

|Flag|Short|Description|
|-|-|-|
|`--user`|`-u`|The username to connect with. If not specified, you will be asked for it.|
|`--password`|`-p`|The password to authenticate with. If not specified, you will be asked for it.|
|`--to-exchange`||Denote that the binding target is another exchange.|

**Example:**

Create a binding from `my-exchange` to `my-queue` on a RabbitMQ server running on the local machine.

```
$ buneary create binding localhost my-exchange my-queue my-binding-key
```

### Get all exchanges

**Syntax:**

```
$ buneary get exchanges <ADDRESS> [flags]
```

**Arguments:**

|Argument|Description|
|-|-|
|`ADDRESS`|The RabbitMQ HTTP API address. If no port is specified, `15672` is used.|

**Flags:**

|Flag|Short|Description|
|-|-|-|
|`--user`|`-u`|The username to connect with. If not specified, you will be asked for it.|
|`--password`|`-p`|The password to authenticate with. If not specified, you will be asked for it.|

**Example:**

Get all exchanges from a RabbitMQ server running on the local machine - this particular example also shows the output.

```
$ buneary get exchanges localhost
User: guest
Password:
+--------------------+---------+---------+-------------+----------+
|        NAME        |  TYPE   | DURABLE | AUTO-DELETE | INTERNAL |
+--------------------+---------+---------+-------------+----------+
|                    | direct  | yes     | no          | no       |
| amq.direct         | direct  | yes     | no          | no       |
| amq.fanout         | fanout  | yes     | no          | no       |
| amq.headers        | headers | yes     | no          | no       |
| amq.match          | headers | yes     | no          | no       |
| amq.rabbitmq.trace | topic   | yes     | no          | yes      |
| amq.topic          | topic   | yes     | no          | no       |
+--------------------+---------+---------+-------------+----------+

```

### Get an exchange

**Syntax:**

```
$ buneary get exchange <ADDRESS> <NAME> [flags]
```

**Arguments:**

|Argument|Description|
|-|-|
|`ADDRESS`|The RabbitMQ HTTP API address. If no port is specified, `15672` is used.|
|`NAME`|The name of the exchange.|

**Flags:**

|Flag|Short|Description|
|-|-|-|
|`--user`|`-u`|The username to connect with. If not specified, you will be asked for it.|
|`--password`|`-p`|The password to authenticate with. If not specified, you will be asked for it.|

**Example:**

Get an exchange called `my-exchange` from a RabbitMQ server running on the local machine.

```
$ buneary get exchange localhost my-exchange
```

### Get all queues

**Syntax:**

```
$ buneary get queues <ADDRESS> [flags]
```

**Arguments:**

|Argument|Description|
|-|-|
|`ADDRESS`|The RabbitMQ HTTP API address. If no port is specified, `15672` is used.|

**Flags:**

|Flag|Short|Description|
|-|-|-|
|`--user`|`-u`|The username to connect with. If not specified, you will be asked for it.|
|`--password`|`-p`|The password to authenticate with. If not specified, you will be asked for it.|

**Example:**

Get all queues from a RabbitMQ server running on the local machine.

```
$ buneary get queues localhost
```

### Get a queue

**Syntax:**

```
$ buneary get queue <ADDRESS> <NAME> [flags]
```

**Arguments:**

|Argument|Description|
|-|-|
|`ADDRESS`|The RabbitMQ HTTP API address. If no port is specified, `15672` is used.|
|`NAME`|The name of the queue.|

**Flags:**

|Flag|Short|Description|
|-|-|-|
|`--user`|`-u`|The username to connect with. If not specified, you will be asked for it.|
|`--password`|`-p`|The password to authenticate with. If not specified, you will be asked for it.|

**Example:**

Get a queue called `my-queue` from a RabbitMQ server running on the local machine.

```
$ buneary get queue localhost my-exchange
```

### Get all bindings

**Syntax:**

```
$ buneary get bindings <ADDRESS> [flags]
```

**Arguments:**

|Argument|Description|
|-|-|
|`ADDRESS`|The RabbitMQ HTTP API address. If no port is specified, `15672` is used.|

**Flags:**

|Flag|Short|Description|
|-|-|-|
|`--user`|`-u`|The username to connect with. If not specified, you will be asked for it.|
|`--password`|`-p`|The password to authenticate with. If not specified, you will be asked for it.|

**Example:**

Get all bindings from a RabbitMQ server running on the local machine.

```
$ buneary get bindings localhost
```

### Get a binding

**Syntax:**

```
$ buneary get binding <ADDRESS> <EXCHANGE NAME> <TARGET NAME> [flags]
```

**Arguments:**

|Argument|Description|
|-|-|
|`ADDRESS`|The RabbitMQ HTTP API address. If no port is specified, `15672` is used.|
|`EXCHANGE NAME`|The name of the source exchange.|
|`TARGET NAME`|The name of the target.|

**Flags:**

|Flag|Short|Description|
|-|-|-|
|`--user`|`-u`|The username to connect with. If not specified, you will be asked for it.|
|`--password`|`-p`|The password to authenticate with. If not specified, you will be asked for it.|

**Example:**

Get the binding or bindings between `my-exchange` and `my-queue` from a RabbitMQ server running on the local machine.

```
$ buneary get binding localhost my-exchange my-queue
```

### Get messages in a queue

**Syntax:**

```
$ buneary publish <ADDRESS> <QUEUE NAME> [flags]
```

**Arguments:**

|Argument|Description|
|-|-|
|`ADDRESS`|The RabbitMQ AMQP address. If no port is specified, `5672` is used.|
|`QUEUE NAME`|The name of the queue to read messages from.|

**Flags:**

|Flag|Short|Description|
|-|-|-|
|`--user`|`-u`|The username to connect with. If not specified, you will be asked for it.|
|`--password`|`-p`|The password to authenticate with. If not specified, you will be asked for it.|
|`--max`||The maximum amount of messages to read from the queue.|
|`--requeue`||Reading messages will de-queue them. Re-queue the messages after reading them.|
|`--force`|`-f`|Skip the manual confirmation and force reading the messages.|

**Example:**

Read up to 10 messages from the `my-queue` queue on a RabbitMQ server running on the local machine.

```
$ buneary get messages --max 10 localhost my-queue
```

### Publish a message

**Syntax:**

```
$ buneary publish <ADDRESS> <EXCHANGE> <ROUTING KEY> <BODY> [flags]
```

**Arguments:**

|Argument|Description|
|-|-|
|`ADDRESS`|The RabbitMQ AMQP address. If no port is specified, `5672` is used.|
|`EXCHANGE`|The name of the target exchange.|
|`ROUTING KEY`|The routing key of the message.|
|`BODY`|The actual message body.|

**Flags:**

|Flag|Short|Description|
|-|-|-|
|`--user`|`-u`|The username to connect with. If not specified, you will be asked for it.|
|`--password`|`-p`|The password to authenticate with. If not specified, you will be asked for it.|
|`--headers`||Comma-separated message headers in the form `--headers key1=val1,key2=val2`.|

**Example:**

Publish a message `Hello!` to `my-exchange` on a RabbitMQ server running on the local machine.

```
$ buneary publish localhost my-exchange my-routing-key "Hello!"
```

### Delete an exchange

**Syntax:**

```
$ buneary delete exchange <ADDRESS> <NAME> [flags]
```

**Arguments:**

|Argument|Description|
|-|-|
|`ADDRESS`|The RabbitMQ HTTP API address. If no port is specified, `15672` is used.|
|`NAME`|The name of the exchange to be deleted.|

**Flags:**

|Flag|Short|Description|
|-|-|-|
|`--user`|`-u`|The username to connect with. If not specified, you will be asked for it.|
|`--password`|`-p`|The password to authenticate with. If not specified, you will be asked for it.|

**Example:**

Delete an exchange called `my-exchange` on a RabbitMQ server running on the local machine.

```
$ buneary delete exchange localhost my-exchange
```

### Delete a queue

**Syntax:**

```
$ buneary delete queue <ADDRESS> <NAME> [flags]
```

**Arguments:**

|Argument|Description|
|-|-|
|`ADDRESS`|The RabbitMQ HTTP API address. If no port is specified, `15672` is used.|
|`NAME`|The name of the queue to be deleted.|

**Flags:**

|Flag|Short|Description|
|-|-|-|
|`--user`|`-u`|The username to connect with. If not specified, you will be asked for it.|
|`--password`|`-p`|The password to authenticate with. If not specified, you will be asked for it.|

**Example:**

Delete a queue called `my-queue` on a RabbitMQ server running on the local machine.

```
$ buneary delete queue localhost my-queue
```

## Credits

* [michaelklishin/rabbit-hole](https://github.com/michaelklishin/rabbit-hole) is used as RabbitMQ client library.
* [streadway/amqp](https://github.com/streadway/amqp) is used as AMQP client library.
* For all third-party packages used, see [go.mod](go.mod).
* The Buneary graphic is made by [dirocha](https://imgbin.com/user/dirocha).
