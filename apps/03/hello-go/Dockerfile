# syntax=docker/dockerfile:1

FROM tinygo/tinygo:0.28.1 as builder

WORKDIR /src

RUN mkdir -p /home/tinygo/.cache/tinygo

COPY . .

RUN --mount=type=cache,target=/home/tinygo/.cache/tinygo,uid=1000,gid=1000 \
    tinygo build -target=wasi -gc=leaking -no-debug -o main.wasm main.go

FROM scratch
COPY --from=builder /src/main.wasm /
COPY --from=builder /src/spin.toml /
