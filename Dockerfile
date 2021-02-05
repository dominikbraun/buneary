# This Dockerfile builds a lightweight distribution image for Docker Hub.
# It only contains the application without any source code.
FROM alpine:3.11.5 AS downloader

# The buneary release to be downloaded from GitHub.
ARG VERSION

RUN apk add --no-cache \
    curl \
    tar

RUN curl -LO https://github.com/dominikbraun/buneary/releases/download/${VERSION}/verless-linux-amd64.tar && \
    tar -xvf buneary-linux-amd64.tar -C /bin && \
    rm -f buneary-linux-amd64.tar

# The final stage. This is the image that will be distrubuted.
FROM alpine:3.11.5 AS final

LABEL org.label-schema.schema-version="1.0"
LABEL org.label-schema.name="buneary"
LABEL org.label-schema.description="An easy-to-use CLI client for RabbitMQ."
LABEL org.label-schema.url="https://github.com/dominikbraun/buneary"
LABEL org.label-schema.vcs-url="https://github.com/dominikbraun/buneary"
LABEL org.label-schema.version=${VERSION}

COPY --from=downloader ["/bin/buneary", "/bin/buneary"]

# Create a symlink for musl, see https://stackoverflow.com/a/35613430.
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

WORKDIR /project

ENTRYPOINT ["/bin/buneary"]