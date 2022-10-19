# Pod debug tool

For daily positioning pods in which the application cannot start or want to copy a pod by itself, support
Deployment/StatefulSet/Job/Pod resources, you can input json/yaml format from a file or console and specify the output
json/yaml format for 'kubectl apply -f -' , the generated pod will be 'sh -c /usr/sbin/sshd -D & touch debug && tail -f
debug' blocks and does not exit

# Packaging

go build -o pdg -ldflags="-s -w" main.go
