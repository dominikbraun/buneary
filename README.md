<p align="center">
<br>
<img src="logo.png" alt="buneary">
<br>
<br>
</p>

`buneary`, pronounced _bun-ear-y_, is an easy-to-use RabbitMQ command line client for managing exchanges, managing
queues and publishing messages to exchanges.

## Installation

### macOS/Linux

Download the [latest release](https://github.com/dominikbraun/buneary/releases) for your platform. Extract the
downloaded binary into a directory like `/usr/local/bin`. Make sure the directory is in `PATH`.

### Windows

Download the [latest release](https://github.com/dominikbraun/buneary/releases), create a directory like
`C:\Program Files\buneary` and extract the executable into that directory.
[Add the directory to `Path`](https://www.computerhope.com/issues/ch000549.htm).

### Docker

Just run `docker container run dominikbraun/buneary` followed by the actual command you want to execute.

```
$ docker container run dominikbraun/buneary publish localhost my-exchange my-routingkey "Hello!"
```

## Usage

### Publishing messages

...

### Creating exchanges

...

## Creating queues

...