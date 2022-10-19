# Pod debug tool

For daily positioning pods in which the application cannot start or want to copy a pod by itself, support
Deployment/StatefulSet/Job/Pod resources, you can input json/yaml format from a file or console and specify the output
json/yaml format for 'kubectl apply -f -' , the generated pod will be 'sh -c /usr/sbin/sshd -D & touch debug && tail -f
debug' blocks and does not exit

# use

```shell
kubectl get deployment/xxx -o json |pdebug | kubectl apply -f -
kubectl get statesulset/xxx -o json |pdebug | kubectl apply -f -
kubectl get pod/xxx -o json |pdebug | kubectl apply -f -
kubectl get job/xxx -o json |pdebug | kubectl apply -f -
cat pod.yaml |pdebug | kubectl apply -f -
pdebug -f pod.json| kubectl apply -f -
pdebug -f pod.json -t yaml
```

# Packaging

```shell
go mod tidy
go build -o pdebug -ldflags="-s -w" main.go
go build -mod=mod

```