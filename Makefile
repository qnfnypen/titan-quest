export GOPROXY=https://goproxy.cn

all: gen build
.PHONY: all

GOCC?=go

titan-quest:
	rm -f titan-quest
	$(GOCC) build $(GOFLAGS) -o titan-quest .
.PHONY: titan-explorer

gen:
	sqlc generate
.PHONY: gen

build: titan-quest
.PHONY: build
