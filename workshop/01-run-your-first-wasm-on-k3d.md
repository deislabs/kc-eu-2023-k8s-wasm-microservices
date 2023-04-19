# Run your first spin application on k3d

## What is runwasi?

[Runwasi](https://github.com/containerd/runwasi) is governed by the CNCF Containerd project. It operates as a containerd shim and allows you to run WASI applications on Kubernetes. It is a great way to get started with WASI and WebAssembly on Kubernetes.

## To get started, we will use k3d

The shell script below will create a k3d cluster locally with the Wasm shims installed and containerd configured.

```
k3d cluster create wasm-cluster --image ghcr.io/deislabs/containerd-wasm-shims/examples/k3d:v0.5.1 -p "8081:80@loadbalancer" --agents 2
```

Let's apply the runtime class and workloads for the shim: 
```
kubectl apply -f https://github.com/deislabs/containerd-wasm-shims/raw/main/deployments/workloads/runtime.yaml
kubectl apply -f https://github.com/deislabs/containerd-wasm-shims/raw/main/deployments/workloads/workload.yaml
```

> If you have trouble with the Internet connection, you can find the files in the `apps/01` directory of this repo.

Now let's test it out:
```
echo "waiting 5 seconds for workload to be ready"
sleep 5
curl -v http://127.0.0.1:8081/spin/hello
```