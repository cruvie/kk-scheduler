# brew install just

set shell := ["zsh", "-ic"]

default: fmt

which-shell:
	echo ${SHELL}
	pwd

fmt:
	go tool gofumpt -l -w .

lint:
	go tool golangci-lint version
	go tool golangci-lint run

deadcode:
	go tool deadcode ./...

govulncheck:
	go tool govulncheck ./...

capslock:
    # 要检查的包
    cd kk_nats
    go tool capslock

race_detector:
  #https://zhuanlan.zhihu.com/p/78655582
  go run -race ./main/main.go

fix:
	go fix ./...

# 仅做参考，运行前需要删除vendor，目前无法排除https://github.com/uber-go/nilaway/issues/99
nilaway:
	go tool nilaway ./...

update-dep:
	go get -u ./...
	go get tool
	go mod tidy

cache-clean:
    go clean -cache

proto-gen:
    buf generate

test:
	go test ./...
