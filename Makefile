build:
	GOPROXY=https://goproxy.cn,direct \
    go build --trimpath --ldflags="-w -s" -o release/bin/gitlab_bot cmd/main.go