.PHONY:
.SILENT:
.DEFAULT_GOAL:= run

help: ## Отобразить описание справки
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' Makefile | sort | awk \
	'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

init: gen ## Инициализация приложения
	go mod tidy
	go mod vendor

gen: ## Генерация proto-файлов
		mkdir -p pkg/blog.v1
		protoc 	--proto_path=api/blog.v1 \
	            --proto_path=proto \
				-I api/blog.v1\
				--go_out=pkg/blog.v1 --go_opt=paths=import \
				--go-grpc_out=pkg/blog.v1 --go-grpc_opt=paths=import \
				--grpc-gateway_out=pkg/blog.v1 \
				--grpc-gateway_opt=logtostderr=true \
				--grpc-gateway_opt=paths=import \
				--swagger_out=allow_merge=true,merge_file_name=proto:docs \
				api/blog.v1/blog.system.proto
		mv pkg/blog.v1/github.com/Shemistan/blog/pkg/blog.v1/* pkg/blog.v1/
		rm -rf pkg/blog.v1/github.com

lint:
	golangci-lint run
