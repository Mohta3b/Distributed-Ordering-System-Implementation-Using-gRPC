SRC_DIR := src

OUT_DIR := bin

DATA_DIR = data

GO_BUILD := go build


SUBDIRS := $(wildcard $(SRC_DIR)/*)
CLIENT_DIRS := $(wildcard $(SUBDIRS:%=%/client))
SERVER_DIRS := $(wildcard $(SUBDIRS:%=%/server))
HANDLER_DIR = handler

MAINS := $(patsubst $(SRC_DIR)/%/client,%,$(CLIENT_DIRS))

.DEFAULT_GOAL := all

.PHONY: all clean

all: $(MAINS)

$(MAINS): %: $(OUT_DIR)/%_server $(OUT_DIR)/%_client

$(OUT_DIR)/%_server: $(SRC_DIR)/%/server/server.go
	@mkdir -p $(OUT_DIR)
	@$(GO_BUILD) -o $@ $<

$(OUT_DIR)/%_client: $(SRC_DIR)/%/client/client.go
	@mkdir -p $(OUT_DIR)
	@$(GO_BUILD) -o $@ $<


$(OUT_DIR)/data: $(DATA_DIR)/dataset.go
	@mkdir -p $(OUT_DIR)
	@$(GO_BUILD) -o $@ $<


$(OUT_DIR)/handler: $(SRC_DIR)/$(HANDLER_DIR)/handler.go
	@mkdir -p $(OUT_DIR)
	@$(GO_BUILD) -o $@ $<

clean:
	@rm -rf $(OUT_DIR)

