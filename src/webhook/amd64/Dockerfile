# Our first FROM statement declares the build environment.
FROM rust

RUN rustup target add armv7-unknown-linux-gnueabihf
RUN apt-get update -yq && apt-get install -y gcc-arm*

WORKDIR /tmp
RUN git clone https://github.com/nukesor/webhook-server

WORKDIR /tmp/webhook-server
RUN mkdir -p /tmp/webhook-server/.cargo && mkdir -p /root/.cargo && mkdir /.cargo
COPY ./config /tmp/webhook-server/.cargo/config
COPY ./config /root/.cargo/config
COPY ./config /.cargo/config

RUN cargo build --release --target=armv7-unknown-linux-gnueabihf
#RUN find . -type f -executable

# Now, we need to build our _real_ Docker container, copying in `using-diesel`.
FROM alpine:latest

COPY --from=builder \
    /tmp/webhook-server/target/armv7-unknown-linux-gnueabihf/release/webhookserver /tmp/webhookserver

CMD /tmp/webhookserver
