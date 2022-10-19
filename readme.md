# pod debug工具

用于日常定位pod中应用无法启动或者想自己copy一个pod,支持Deployment/StatefulSet/Job/Pod资源,可以从文件或控制台输入json/yaml格式并指定输出json/yaml格式用于`kubectl apply -f -`
,生成的pod会被`sh -c /usr/sbin/sshd -D & touch debug && tail -f debug`阻塞住从而不退出

# 使用

```shell
kubectl get deployment/xxx -o json |pdebug | kubectl apply -f -
kubectl get statesulset/xxx -o json |pdebug | kubectl apply -f -
kubectl get pod/xxx -o json |pdebug | kubectl apply -f -
kubectl get job/xxx -o json |pdebug | kubectl apply -f -
cat pod.yaml |pdebug | kubectl apply -f -
pdebug -f pod.json| kubectl apply -f -
pdebug -f pod.json -t yaml
```

# 打包

```shell
go mod tidy
go build -o pdebug -ldflags="-s -w" main.go
go build -mod=mod

```
