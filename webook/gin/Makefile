.PHONY: docker
docker:
#先删掉旧的本地文件
	@rm webook-gin || true
#编译构建
	@GOOS=linux GOARCH=arm go build -o webook-gin .
	@docker rmi -f xiaobao/webook-gin:v0.0.1 .
	@docker build -t xiaobao/webook-gin:v0.0.1 .
