# pod debug工具

用于日常定位pod中应用无法启动或者想自己copy一个pod,支持Deployment/StatefulSet/Job/Pod资源,可以从文件或控制台输入json/yaml格式并指定输出json/yaml格式用于`kubectl apply -f -`
,生成的pod会被`sh -c /usr/sbin/sshd -D & touch debug && tail -f debug`阻塞住从而不退出

# 打包

go build -o pdg -ldflags="-s -w" main.go
