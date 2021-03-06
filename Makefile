# I'm sure there is better way. But I would need to find it first
MK_FILE_PATH = $(lastword $(MAKEFILE_LIST))
PRJ_DIR      = $(abspath $(dir $(MK_FILE_PATH)))
GO           ?= go
GOPATH_LOCAL = $(PRJ_DIR)/build/
GOPATH_DIR   = src/github.com/henrydcase/nobs
VENDOR_DIR   = tls_vendor
OPTS         ?=
NOASM        ?=
TEST_PATH    ?= ./...
GOCACHE      ?= off
BENCH_OPTS   ?= -v -bench=. -run="NonExistingTest" -benchmem
TEST_PATH    ?= ./...

ifeq ($(NOASM),1)
	OPTS+=$(OPTS_TAGS)
endif

ifeq ($(PPROF),1)
	BENCH_OPTS+= -cpuprofile=cpu.out -memprofile=mem0.out
endif

TARGETS ?= \
	dh   \
	drbg \
	ec \
	hash \
	kem \
	utils

prep-%:
	mkdir -p $(GOPATH_LOCAL)/$(GOPATH_DIR)
	cp -rf $* $(GOPATH_LOCAL)/$(GOPATH_DIR)/$*

make_dirs:
	mkdir -p $(GOPATH_LOCAL)/$(GOPATH_DIR)
	cp -rf etc $(GOPATH_LOCAL)/$(GOPATH_DIR)

test: clean make_dirs $(addprefix prep-,$(TARGETS))
	cd $(GOPATH_LOCAL); GOPATH=$(GOPATH_LOCAL) go test $(OPTS) -v $(TEST_PATH)

cover:
	cd $(GOPATH_LOCAL); GOPATH=$(GOPATH_LOCAL) go test \
		-race -coverprofile=coverage_$(NOASM).txt -covermode=atomic $(OPTS) -v $(TEST_PATH)
	cat $(GOPATH_LOCAL)/coverage_$(NOASM).txt >> coverage.txt

bench: clean $(addprefix prep-,$(TARGETS))
	cd $(GOPATH_LOCAL); GOCACHE=$(GOCACHE) GOPATH=$(GOPATH_LOCAL) $(GO) test \
		$(BENCH_OPTS) $(TEST_PATH)

clean:
	rm -rf $(GOPATH_LOCAL)
	rm -rf $(VENDOR_DIR)

vendor-sidh-for-tls: clean
	mkdir -p $(VENDOR_DIR)/github_com/henrydcase/nobs/
	rsync -a . $(VENDOR_DIR)/github_com/henrydcase/nobs/ --exclude=$(VENDOR_DIR) --exclude=.git --exclude=.travis.yml --exclude=README.md
	find $(VENDOR_DIR) -type f -print0 -name "*.go" | xargs -0 sed -i 's/github\.com/github_com/g'

pprof-cpu:
	$(GO) tool pprof $(GOPATH_LOCAL)/cpu.out

pprof-mem:
	$(GO) tool pprof $(GOPATH_LOCAL)/mem0.out
