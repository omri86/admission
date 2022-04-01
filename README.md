# Kubernetes Admission Control Demo

This repo contains files that are used as a demo of using k8s admission controllers.

**Do not use this example in production!**



First run these commands to generate certificates:

```
cd certs
chmod +x generate_keys.sh
./generate_keys.sh
```



Then apply the webhook config to your cluster:

```
kubectl apply -f webhook.yaml
```



Next build and run the demo server:

```
go build main.go 
./main
```



To validate the server works as expected - the next command should fail:

```
kubectl apply -f objects/nginx-pod.yaml
```

and this should succeed:

```
kubectl apply -f objects/busbox-pod.yaml
```

