.PHONY: build

NODE_VERSION := $(shell ./version.sh)
NODE_COMMIT := $(shell [ -z "${NODE_COMMIT}" ] && git log -1 --format='%H' || echo ${NODE_COMMIT} )

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=musecore \
	-X github.com/cosmos/cosmos-sdk/version.ServerName=musecored \
	-X github.com/cosmos/cosmos-sdk/version.ClientName=museclientd \
	-X github.com/cosmos/cosmos-sdk/version.Version=$(NODE_VERSION) \
	-X github.com/cosmos/cosmos-sdk/version.Commit=$(NODE_COMMIT) \
	-X github.com/cosmos/cosmos-sdk/types.DBBackend=pebbledb \
	-X github.com/RWAs-labs/muse/pkg/constant.Name=musecored \
	-X github.com/RWAs-labs/muse/pkg/constant.Version=$(NODE_VERSION) \
	-X github.com/RWAs-labs/muse/pkg/constant.CommitHash=$(NODE_COMMIT) \
	-buildid= \
	-s -w

BUILD_FLAGS := -ldflags '$(ldflags)' -tags pebbledb,ledger

###############################################################################
###                          Install commands                               ###
###############################################################################

install: go.sum
		@echo "--> Installing musecored, museclientd, and museclientd-supervisor"
		@go install -mod=readonly $(BUILD_FLAGS) ./cmd/musecored
		@go install -mod=readonly $(BUILD_FLAGS) ./cmd/museclientd
		@go install -mod=readonly $(BUILD_FLAGS) ./cmd/museclientd-supervisor

install-museclient: go.sum
		@echo "--> Installing museclientd"
		@go install -mod=readonly $(BUILD_FLAGS) ./cmd/museclientd

install-musecore: go.sum
		@echo "--> Installing musecored"
		@go install -mod=readonly $(BUILD_FLAGS) ./cmd/musecored

# running with race detector on will be slow
install-museclient-race-test-only-build: go.sum
		@echo "--> Installing museclientd"
		@go install -race -mod=readonly $(BUILD_FLAGS) ./cmd/museclientd

install-musetool: go.sum
		@echo "--> Installing musetool"
		@go install -mod=readonly $(BUILD_FLAGS) ./cmd/musetool

###############################################################################
###                           Generation commands  		                    ###
###############################################################################
DOCKER := $(shell which docker)

protoVer=latest
protoImageName=ghcr.io/cosmos/proto-builder:$(protoVer)
protoImage=$(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace --user $(shell id -u):$(shell id -g) $(protoImageName)

proto-format:
	@echo "--> Formatting Protobuf files"
	@$(protoImage) find ./ -name "*.proto" -exec clang-format -i {} \;
.PHONY: proto-format

typescript: proto-format
	@echo "--> Generating TypeScript bindings"
	@bash ./scripts/protoc-gen-typescript.sh
.PHONY: typescript

proto-gen: proto-format
	@echo "--> Removing old Go types "
	@find . -name '*.pb.go' -type f -delete
	@echo "--> Generating Protobuf files"
	@$(protoImage) sh ./scripts/protoc-gen-go.sh

openapi: proto-format
	@echo "--> Generating OpenAPI specs"
	@bash ./scripts/protoc-gen-openapi.sh
.PHONY: openapi

THIRD_PARTY_DIR=./third_party
DEPS_COSMOS_SDK_VERSION := $(shell cat go.sum | grep -E 'github.com/cosmos/cosmos-sdk\s' | grep -v -e 'go.mod' | tail -n 1 | awk '{ print $$2; }')
DEPS_IBC_GO_VERSION := $(shell cat go.sum | grep 'github.com/cosmos/ibc-go' | grep -v -e 'go.mod' | tail -n 1 | awk '{ print $$2; }')
DEPS_COSMOS_PROTO := $(shell cat go.sum | grep 'github.com/cosmos/cosmos-proto' | grep -v -e 'go.mod' | tail -n 1 | awk '{ print $$2; }')
DEPS_COSMOS_GOGOPROTO := $(shell cat go.sum | grep 'github.com/cosmos/gogoproto' | grep -v -e 'go.mod' | tail -n 1 | awk '{ print $$2; }')
DEPS_COSMOS_ICS23 := go/$(shell cat go.sum | grep 'github.com/cosmos/ics23/go' | grep -v -e 'go.mod' | tail -n 1 | awk '{ print $$2; }')

proto-download-deps:
	mkdir -p "$(THIRD_PARTY_DIR)/cosmos_tmp" && \
	cd "$(THIRD_PARTY_DIR)/cosmos_tmp" && \
	git init && \
	git remote add origin "https://github.com/cosmos/cosmos-sdk.git" && \
	git config core.sparseCheckout true && \
	printf "proto\nthird_party\n" > .git/info/sparse-checkout && \
	git pull origin "$(DEPS_COSMOS_SDK_VERSION)" && \
	rm -f ./proto/buf.* && \
	mv ./proto/* ..
	rm -rf "$(THIRD_PARTY_DIR)/cosmos_tmp"

	mkdir -p "$(THIRD_PARTY_DIR)/ibc_tmp" && \
	cd "$(THIRD_PARTY_DIR)/ibc_tmp" && \
	git init && \
	git remote add origin "https://github.com/cosmos/ibc-go.git" && \
	git config core.sparseCheckout true && \
	printf "proto\n" > .git/info/sparse-checkout && \
	git pull origin "$(DEPS_IBC_GO_VERSION)" && \
	rm -f ./proto/buf.* && \
	mv ./proto/* ..
	rm -rf "$(THIRD_PARTY_DIR)/ibc_tmp"

	mkdir -p "$(THIRD_PARTY_DIR)/cosmos_proto_tmp" && \
	cd "$(THIRD_PARTY_DIR)/cosmos_proto_tmp" && \
	git init && \
	git remote add origin "https://github.com/cosmos/cosmos-proto.git" && \
	git config core.sparseCheckout true && \
	printf "proto\n" > .git/info/sparse-checkout && \
	git pull origin "$(DEPS_COSMOS_PROTO_VERSION)" && \
	rm -f ./proto/buf.* && \
	mv ./proto/* ..
	rm -rf "$(THIRD_PARTY_DIR)/cosmos_proto_tmp"

	mkdir -p "$(THIRD_PARTY_DIR)/gogoproto" && \
	curl -SSL "https://raw.githubusercontent.com/cosmos/gogoproto/$(DEPS_COSMOS_GOGOPROTO)/gogoproto/gogo.proto" > "$(THIRD_PARTY_DIR)/gogoproto/gogo.proto"

	mkdir -p "$(THIRD_PARTY_DIR)/google/api" && \
	curl -sSL https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/annotations.proto > "$(THIRD_PARTY_DIR)/google/api/annotations.proto"
	curl -sSL https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/http.proto > "$(THIRD_PARTY_DIR)/google/api/http.proto"

	mkdir -p "$(THIRD_PARTY_DIR)/cosmos/ics23/v1" && \
	curl -sSL "https://raw.githubusercontent.com/cosmos/ics23/$(DEPS_COSMOS_ICS23)/proto/cosmos/ics23/v1/proofs.proto" > "$(THIRD_PARTY_DIR)/cosmos/ics23/v1/proofs.proto"