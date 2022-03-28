 SET CGO_ENABLED=0
 SET GOOS=linux
 SET GOARCH=amd64
 go mod tidy
 go build -o cessctl ../cessctl/main.go