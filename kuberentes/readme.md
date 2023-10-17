## Redis on Kubernetes

Using kind to create the cluster with namespace as redis
for kin [installation](https://kind.sigs.k8s.io/)

```
kind create cluster --name redis --image kindest/node:v1.23.5
```

### creating namespace

```
kubectl create ns redis
```

### Storage class

```
kubectl get storageclass
NAME                 PROVISIONER             RECLAIMPOLICY   VOLUMEBINDINGMODE      ALLOWVOLUMEEXPANSION   AGE
standard (default)   rancher.io/local-path   Delete          WaitForFirstConsumer   false                  84s
```

### Deployment : redis nodes 

```
cd storage/redis/kubernetes/
kubectl apply -n redis -f ./redis/redis-configmap.yaml
kubectl apply -n redis -f ./redis/redis-statefulset.yaml

kubectl -n redis get pods
kubectl -n redis get pv

kubectl -n redis logs redis-0
kubectl -n redis logs redis-1
kubectl -n redis logs redis-2
```
### replication status

```
kubectl -n redis exec -it redis-0 -- sh
redis-cli 
auth a-very-complex-password-here
info replication
```

### Deploying sentinels

```
cd storage/redis/kubernetes/
kubectl apply -n redis -f ./sentinel/sentinel-statefulset.yaml

kubectl -n redis get pods
kubectl -n redis get pv
kubectl -n redis logs sentinel-0
```

### Cleanup

```
kind delete cluster --name redis
```
