.PHONY: mock
mock:
	@mockgen -source=webook/gin/internal/service/user.go -package=svcmocks -destination=webook/gin/internal/service/mock/user.mock.gen.go
	@mockgen -source=webook/gin/internal/service/code.go -package=svcmocks -destination=webook/gin/internal/service/mock/code.mock.gen.go
	@go mod tidy

