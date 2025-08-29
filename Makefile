BINARY_NAME=cli
SRC=./cmd/cli

BIN_DIR=bin

PLATFORMS = \
    darwin/amd64 \
    darwin/arm64 \
    linux/amd64 \
    linux/arm64 \
    windows/amd64 \
    windows/arm64

define build_template
$(BIN_DIR)/$(BINARY_NAME)-$(1)-$(2)$(3): 
    GOOS=$(1) GOARCH=$(2) go build -o $(BIN_DIR)/$(BINARY_NAME)-$(1)-$(2)$(3) $(SRC)
endef

$(foreach plat,$(PLATFORMS),\
  $(eval \
    $(call build_template,\
      $(word 1,$(subst /, ,$(plat))),\
      $(word 2,$(subst /, ,$(plat))),\
      $(if $(findstring windows,$(plat)),.exe,)\
    )\
  )\
)

build-all: $(foreach plat,$(PLATFORMS),$(BIN_DIR)/$(BINARY_NAME)-$(word 1,$(subst /, ,$(plat)))-$(word 2,$(subst /, ,$(plat)))$(if $(findstring windows,$(plat)),.exe,)))

.PHONY: clean
clean:
    rm -rf $(BIN_DIR)

.PHONY: all
all: build-all

# Individual targets for convenience
darwin-amd64: $(BIN_DIR)/$(BINARY_NAME)-darwin-amd64
darwin-arm64: $(BIN_DIR)/$(BINARY_NAME)-darwin-arm64
linux-amd64: $(BIN_DIR)/$(BINARY_NAME)-linux-amd64
linux-arm64: $(BIN_DIR)/$(BINARY_NAME)-linux-arm64
windows-amd64: $(BIN_DIR)/$(BINARY_NAME)-windows-amd64.exe
windows-arm64: $(BIN_DIR)/$(BINARY_NAME)-windows-arm64.exe