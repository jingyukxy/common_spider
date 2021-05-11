TARGET="./bin/tag_engine"
MAIN_PATH="src/main/main.go"
VERSION="v0.0.1"
DATE= `date +%FT%T%z`

ifeq (${VERSION}, "v0.0.1")
    VERSION=VERSION = "v0.0.1"
endif

.PHONY: version
version:
	@echo ${VERSION}

# .PHONY 有 build 文件，不影响 build 命令执行
.PHONY: build
build:
	@echo version: ${VERSION} date: ${DATE}
	@go  build -o ${TARGET} ${MAIN_PATH}

install:
	@echo download package
	@go mod download

# 交叉编译运行在linux系统环境
build-linux:
	@echo version: ${VERSION} date: ${DATE} os: linux-centOS
	@GOOS=linux go build -o ${PROJECT} ${MAIN_PATH}

run: build
	@./${TARGET}

clean:
	rm -rf logs/*
	rm -f ${TARGET}

