default: all

version := 0.1-alpha

unamestr := $(shell uname)

ifeq ($(unamestr),Linux)
  platform := linux
else
  ifeq ($(unamestr),Darwin)
    platform := darwin
  endif
endif

tags := $(platform) sqlite_omit_load_extension
test_tags := $(tags) testbinary
test_flags := -coverpkg github.com/SQLess/SQLess/... -cover -race -c

ldflags_role_client := -X main.version=$(version) -X github.com/SQLess/SQLess/conf.RoleTag=C
ldflags_role_client_simple_log := $(ldflags_role_client) -X github.com/SQLess/SQLess/utils/log.SimpleLog=Y

GOTEST := CGO_ENABLED=1 go test $(test_flags) -tags "$(test_tags)"
GOBUILD := CGO_ENABLED=1 go build -tags "$(tags)"

bin/cql.test: stamp-submodule stamp-gen
	$(GOTEST) \
		-ldflags "$(ldflags_role_client)" \
		-o bin/cql.test \
		github.com/SQLess/SQLess/cmd/cql

bin/cql: stamp-submodule stamp-gen
	$(GOBUILD) \
		-ldflags "$(ldflags_role_client_simple_log)" \
		-o bin/cql \
		github.com/SQLess/SQLess/cmd/cql

stamp-gen: bin/hsp
	PATH=`pwd`/bin:$(PATH) go generate ./...
	touch $@

bin/hsp:
	go get -u github.com/SQLess/HashStablePack/hsp
	go build \
		-o bin/hsp \
		github.com/SQLess/HashStablePack/hsp

stamp-submodule: crypto/secp256k1/libsecp256k1/src/ecdsa.h
	touch $@

crypto/secp256k1/libsecp256k1/src/ecdsa.h:
	git submodule update --init crypto/secp256k1/libsecp256k1

client: bin/cql bin/cql.test

all: client

clean:
	find . -name '*_gen*.go' -exec rm -vf {} \;
	$(RM) bin/cql* bin/hsp
	$(RM) *.cover.out
	$(RM) coverage.txt
	$(RM) stamp-submodule stamp-gen

testnet-clean:
	$(RM) conf/testnet/.dsn
	$(RM) conf/testnet/observer.db3
	$(RM) conf/testnet/public.keystore*
	$(RM) conf/testnet/observer.db3

.PHONY: bin/cql bin/cql.test
