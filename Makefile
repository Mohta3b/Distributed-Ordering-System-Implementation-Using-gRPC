SRC_DIR := src

OUT_DIR := bin

PROTO_DIR := $(SRC_DIR)/proto

PACKAGE_DIR = pkg

PROTO_OUT_DIR := $(OUT_DIR)/proto_bin

GO_BUILD := go build

PROTOC := protoc

PROTOC_PLUGIN_GO := $(shell which protoc-gen-go)

OUT_GO_DIR := $(OUT_DIR)

.DEFAULT_GOAL := all

.PHONY: all clean

all: proto_compile go_build

go_build:
	@mkdir -p $(OUT_GO_DIR)
	@$(GO_BUILD) $(PACKAGE_DIR)/*.go
	@for dir in $(wildcard $(SRC_DIR)/*); do \
		$(GO_BUILD) -o $(OUT_GO_DIR)/$$(basename $$dir) $$dir/*.go; \
	done

# create_bin_directory:
# 	@mkdir -p $(OUT_DIR)

$(OUT_GO_DIR): $(wildcard $(SRC_DIR)/*/) $(wildcard $(SRC_DIR)/*/*/)
	@mkdir -p $(OUT_DIR)
	@$(GO_BUILD) -o $@ $<

# proto_compile: $(PROTO_OUT_DIR)

proto_compile:
	@mkdir -p $(PROTO_OUT_DIR)
	@$(PROTOC) --go_out=$(PROTO_OUT_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(PROTO_OUT_DIR) --go-grpc_opt=paths=source_relative \
		--proto_path=$(PROTO_DIR) $(wildcard $(PROTO_DIR)/*.proto)

clean:
	@rm -rf $(OUT_DIR)
