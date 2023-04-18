FROM scratch AS build
WORKDIR kc-eu-2023-k8s-wasm-microservices/apps/02/hello-rust
COPY . .

FROM scratch
COPY --from=build kc-eu-2023-k8s-wasm-microservices/apps/02/hello-rust/spin.toml .
COPY --from=build kc-eu-2023-k8s-wasm-microservices/apps/02/hello-rust/target/wasm32-wasi/release/hello_rust.wasm ./target/wasm32-wasi/release/hello_rust.wasm
