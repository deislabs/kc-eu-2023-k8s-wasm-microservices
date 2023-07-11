# Deploy your Spin application to a WASM enabled Kubernetes cluster using kwasm

[kwasm](https://kwasm.sh) is a Kubernetes operator to install the necessary
runwasi shim on existing Kubernetes clusters. At the moment, kwasm supports
these Kubernetes distributions:

- Microsoft AKS
- Google GKE (only with Ubuntu nodes)
- Amazon EKS (only with Ubuntu nodes)
- DigitalOcean DOKS
- Civo Kubernetes
- Kind
- Minikube
- MicroK8s

## Requirements

To install the kwasm-operator into a supported cluster, you will need to have
`helm` installed on your machine, as well as access to a supported cluster. If you do not have one yet, you can create a new cluster with KinD:

```
kind create cluster
```

## Installation

The following commands will install the kwasm operator and set up the `wasmtime-spin` runtime class: 

```
# Add HELM repository if not already done
helm repo add kwasm http://kwasm.sh/kwasm-operator/

# Install KWasm operator
helm install -n kwasm --create-namespace kwasm-operator kwasm/kwasm-operator

# Create the runtime class
kubectl apply -f ../apps/05/runtimeclass.yaml
```

Ensure that the kwasm-operator is running correctly and that the runtimeclass has been created:

```
kubectl get pods -n kwasm
NAME                             READY   STATUS    RESTARTS   AGE
kwasm-operator-88788fd4c-nqmfm   1/1     Running   0          46s

kubectl get runtimeclass
NAME            HANDLER   AGE
wasmtime-spin   spin      7s
```

You can now select which nodes you want your Spin workloads to run on by
annotating the nodes with `kwasm.sh/kwasm-node=true`. To annotate all nodes,
you can run the following command:

```
kubectl annotate node --all kwasm.sh/kwasm-node=true
```

This will trigger a job on each matching node in the cluster to install the containerd shims. Make sure the jobs are finished before proceeding.

```
kubectl get jobs -n kwasm
NAME                                 COMPLETIONS   DURATION   AGE
kind-control-plane-provision-kwasm   1/1           10s        76s
```

## Deploy a Spin workload to the cluster

You are now ready to run Spin workloads on your cluster. To deploy the demo
application, execute this command:

```
kubectl apply -f ../apps/05/spin.yaml
```

Let's verify that the service has been created and that the pod is in the `Running` state.

```
kubectl get svc
NAME         TYPE           CLUSTER-IP     EXTERNAL-IP   PORT(S)        AGE
kubernetes   ClusterIP      10.96.0.1      <none>        443/TCP        7m33s
wasm-spin    LoadBalancer   10.96.197.81   <pending>     80:31820/TCP   2m23s

kubectl get pods
NAME                         READY   STATUS    RESTARTS   AGE
wasm-spin-74cb5d68c8-rt7h6   1/1     Running   0          5s
```

### Call the Spin application

Notice that the wasm-spin service does not have an external IP. This is due to
the KinD cluster running locally without any LoadBalancer provider installed.
For the purpose of this demo, we will port-forward this service. If your
cluster was able to assign an external IP, feel free to skip the
port-forwarding part and use the external IP instead.

```
kubectl port-forward svc/wasm-spin 8080:80
Forwarding from 127.0.0.1:8080 -> 80
Forwarding from [::1]:8080 -> 80
```

Now that we are forwarding port 8080 to the Spin service, let's make a request using curl:

```
curl localhost:8080/hello
Hello world from Spin!
```

Great! We can see that we got a response from our Spin workload!

Let's check the logs of the pod:

```
kubectl logs wasm-spin-74cb5d68c8-rt7h6
Hello, world! You should see me in pod logs
```

Congrats! You now know how to run a WASM workload on a standard Kubernetes cluster using Spin and kwasm-operator!

Feel free to delete the KinD cluster by running

```
kind delete cluster
```