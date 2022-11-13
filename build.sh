env GOOS=linux GOARCH=amd64 go build k8sOp
docker build --no-cache -t mondegreen/scheduler:latest .
docker push mondegreen/scheduler:latest