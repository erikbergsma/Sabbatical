# -*- mode: dockerfile -*-
#
#ARG BASE_IMAGE=ekidd/rust-musl-builder:latest

# Our first FROM statement declares the build environment.
#FROM ${BASE_IMAGE} AS builder

FROM rust

WORKDIR /tmp
RUN git clone https://github.com/nukesor/webhook-server

WORKDIR /tmp/webhook-server
RUN cargo build --release --target=armv7-unknown-linux-gnueabihf
#RUN find . 

# Now, we need to build our _real_ Docker container, copying in `using-diesel`.
FROM ubuntu:20.04

COPY --from=builder \ 
    /tmp/webhook-server/target/armv7-unknown-linux-gnueabihf/release/webhookserver /tmp/webhookserver

CMD /tmp/webhookserver
