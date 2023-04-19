# Slight 101

[Slight](https://github.com/deislabs/spiderlightning) is an experimental runtime implementation for running WebAssembly modules with `wasi-cloud-core` capabilities.

`wasi-cloud-core` is a wasm world that defines a core set of interfaces for distributed applications. It includes interfaces to interfact with key-value stores, blob storage, sql databases, message brokers and more. 

You can read more about `wasi-cloud-core` [here](https://github.com/WebAssembly/WASI/issues/520)

Here is a list of all wasi-cloud-core proposals
- [wasi-keyvalue](https://github.com/WebAssembly/wasi-keyvalue)
- [wasi-blob-store](https://github.com/WebAssembly/wasi-blob-store)
- [wasi-distributed-lock-service](https://github.com/WebAssembly/wasi-distributed-lock-service)
- [wasi-messaing](https://github.com/WebAssembly/wasi-messaging)
- [wasi-http](https://github.com/WebAssembly/wasi-http)
- [wasi-runtime-config](https://github.com/WebAssembly/wasi-runtime-config)
- [wasi-distributed-lock-service](https://github.com/WebAssembly/wasi-distributed-lock-service)
- [wasi-sql](https://github.com/WebAssembly/wasi-sql)

### Create your first slight application

To get started, we will install slight on your local machine.
```
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/deislabs/spiderlightning/main/install.sh)"
```

Make sure that you have the latest version of `slight` installed:

```
slight --version
slight 0.4.1
```

Next, create a new rust project with slight:

```
slight new -n spidey@v0.4.1 rust && cd spidey
```

This will create a new rust project with slight. You can find the source code for the project in the `src` directory.

### Install Rust

If you don't have Rust installed, you can install it by running the following command:

```
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
```

Next, you will need to add `wasm32-wasi` target to your rust toolchain. You can do this by running the following command:

```
rustup target add wasm32-wasi
```

### Build and Run
Then you can build the project by running the following command:

```
cargo build --target wasm32-wasi
```

Finally, you can run the project by running the following command:

```
slight -c slightfile.toml run target/wasm32-wasi/debug/spidey.wasm

Hello, SpiderLightning!
```

### Slight keyvalue with Redis

Add http server capability
```
cd wit && slight add http-server@v0.4.1
cd ..
```

Modify the `src/main.rs` file to `src/lib.rs` and change to the following code:

```rust
use anyhow::Result;

use http_server::*;
use keyvalue::*;
use slight_http_handler_macro::register_handler;
use slight_http_server_macro::on_server_init;

wit_bindgen_rust::import!("wit/http-server.wit");
wit_bindgen_rust::export!("wit/http-server-export.wit");
wit_bindgen_rust::import!("wit/keyvalue.wit");
wit_error_rs::impl_error!(http_server::HttpRouterError);
wit_error_rs::impl_error!(keyvalue::KeyvalueError);

#[on_server_init]
fn main() -> Result<()> {
    let router = Router::new()?;
    let router_with_route = router
        .get("/hello", "handle_hello")?
        .get("/get", "handle_get")?
        .put("/set", "handle_set")?;

    println!("Server is running on port 3000");
    let _ = Server::serve("0.0.0.0:3000", &router_with_route)?;
    Ok(())
}

#[register_handler]
fn handle_hello(req: Request) -> Result<Response, HttpError> {
    println!("I just got a request uri: {} method: {}", req.uri, req.method);
    Ok(Response {
        headers: Some(req.headers),
        body: Some("hello".as_bytes().to_vec()),
        status: 200,
    })
}

#[register_handler]
fn handle_get(request: Request) -> Result<Response, HttpError> {
    let keyvalue =
        Keyvalue::open("my-container").map_err(|e| HttpError::UnexpectedError(e.to_string()))?;

    match keyvalue.get("key") {
        Err(KeyvalueError::KeyNotFound(_)) => Ok(Response {
            headers: Some(request.headers),
            body: Some("Key not found".as_bytes().to_vec()),
            status: 404,
        }),
        Ok(value) => Ok(Response {
            headers: Some(request.headers),
            body: Some(value),
            status: 200,
        }),
        Err(e) => Err(HttpError::UnexpectedError(e.to_string())),
    }
}

#[register_handler]
fn handle_set(request: Request) -> Result<Response, HttpError> {
    assert_eq!(request.method, Method::Put);
    if let Some(body) = request.body {
        let keyvalue = Keyvalue::open("my-container")
            .map_err(|e| HttpError::UnexpectedError(e.to_string()))?;
        keyvalue
            .set("key", &body)
            .map_err(|e| HttpError::UnexpectedError(e.to_string()))?;
    }
    Ok(Response {
        headers: Some(request.headers),
        body: None,
        status: 204,
    })
}
```

Next, change Cargo.toml to the following code:

```toml
[package]
name = "spidey"
version = "0.1.0"
edition = "2021"

# See more keys and their definitions at https://doc.rust-lang.org/cargo/reference/manifest.html
[lib]
crate-type = ["cdylib"]

[dependencies]
anyhow = "1"
wit-bindgen-rust = { git = "https://github.com/fermyon/wit-bindgen-backport" }
wit-error-rs = { git = "https://github.com/danbugs/wit-error-rs", rev = "05362f1a4a3a9dc6a1de39195e06d2d5d6491a5e" }
slight-http-handler-macro = { git = "https://github.com/deislabs/spiderlightning", tag = "0.4.1" }
slight-http-server-macro = { git = "https://github.com/deislabs/spiderlightning", tag = "0.4.1" }

```

Last, add the following capability the slightfile.toml

```toml
[[capability]]
resource = "keyvalue.redis"
name = "my-container"
    [capability.configs]
    REDIS_ADDRESS = "redis://127.0.0.1:6379"
```

### Build and Run

```bash
cargo build --target wasm32-wasi
```

Start a local redis server

```bash
redis-server --port 6379
```

Run the slight

```bash
slight -c slightfile.toml run target/wasm32-wasi/debug/spidey.wasm
```

### Test

```bash
curl -X PUT http://localhost:3000/set -d "hello"

curl http://localhost:3000/get
```

### Try out other slight capabilities

1. [Keyvalue](https://github.com/deislabs/spiderlightning/tree/main/examples/keyvalue-demo)
2. [Messaging](https://github.com/deislabs/spiderlightning/tree/main/tests/messaging-test)
3. [Blob Store](https://github.com/deislabs/spiderlightning/tree/main/examples/blob-store-demo)
4. [HTTP Server](https://github.com/deislabs/spiderlightning/tree/main/examples/http-server-demo)