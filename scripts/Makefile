APP_NAME := cli_mysql_pipeline
GO       := go
GOFLAGS  := -trimpath -ldflags "-s -w"

.PHONY: build clean run full-sync

# 编译
build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p dist
	$(GO) build $(GOFLAGS) \
		-o dist/$(APP_NAME) \
		./instance/partystar/
	@cp configs/partystar.yaml dist/
	@echo "Done: dist/$(APP_NAME)"

# 清理
clean:
	@rm -rf dist/

# 本地运行（开发用）
run:
	$(GO) run ./instance/partystar/ --env=partystar

# 全量同步
# 用法: make full-sync INDEX=song
full-sync:
	$(GO) run ./instance/partystar/ \
		--env=partystar \
		--startWithFull=$(INDEX) \
		--dropIndexName=$(INDEX)
