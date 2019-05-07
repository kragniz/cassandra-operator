.DEFAULT_GOAL := travis
.PHONY: travis clean-all build-all setup setup-all

TEST_REGISTRY ?= localhost:5000
POD_START_TIMEOUT ?= 120s
KUBE_CONTEXT ?= dind
USE_MOCK ?= true
FAKE_CASSANDRA_IMAGE ?= $(TEST_REGISTRY)/fake-cassandra:v$(gitRev)
CASSANDRA_BOOTSTRAPPER_IMAGE ?= $(TEST_REGISTRY)/cassandra-bootstrapper:v$(gitRev)
CASSANDRA_SNAPSHOT_IMAGE ?= $(TEST_REGISTRY)/cassandra-snapshot:v$(gitRev)
NAMESPACE ?= test-cassandra-operator
GINKGO_NODES ?= 0
GINKGO_COMPILERS ?= 0

gitRev := $(shell git rev-parse --short HEAD)
projectDir := $(realpath $(dir $(firstword $(MAKEFILE_LIST))))

travis: make-all

setup: check-system-dependencies setup-all recreate-dind-cluster

install: install-all

clean: clean-all

release: release-all

check-system-dependencies:
	@echo "== check-system-dependencies"
ifeq (, $(shell which go))
	@echo "golang not found, please install it - see https://golang.org/doc/install"
endif
ifeq (, $(shell which docker))
	@echo "docker not found, please install it - see https://docs.docker.com/install"
endif
ifeq (, $(shell which dgoss))
	@echo "dgoss not found, please install it - see https://github.com/aelsabbahy/goss#installation"
endif

make-all:
	@echo "== make-all"
	$(MAKE) -C cassandra-bootstrapper
	$(MAKE) -C fake-cassandra-docker
	GINKGO_NODES=$(GINKGO_NODES) GINKGO_COMPILERS=$(GINKGO_COMPILERS) KUBE_CONTEXT=$(KUBE_CONTEXT) TEST_REGISTRY=$(TEST_REGISTRY) FAKE_CASSANDRA_IMAGE=$(FAKE_CASSANDRA_IMAGE) USE_MOCK=$(USE_MOCK) $(MAKE) -C cassandra-snapshot
	GINKGO_NODES=$(GINKGO_NODES) GINKGO_COMPILERS=$(GINKGO_COMPILERS) KUBE_CONTEXT=$(KUBE_CONTEXT) TEST_REGISTRY=$(TEST_REGISTRY) FAKE_CASSANDRA_IMAGE=$(FAKE_CASSANDRA_IMAGE) CASSANDRA_BOOTSTRAPPER_IMAGE=$(CASSANDRA_BOOTSTRAPPER_IMAGE) CASSANDRA_SNAPSHOT_IMAGE=$(CASSANDRA_SNAPSHOT_IMAGE) USE_MOCK=$(USE_MOCK) POD_START_TIMEOUT=$(POD_START_TIMEOUT) $(MAKE) -C cassandra-operator

install-all:
	@echo "== install-all"
	$(MAKE) -C fake-cassandra-docker install
	$(MAKE) -C cassandra-bootstrapper install
	$(MAKE) -C cassandra-snapshot install
	$(MAKE) -C cassandra-operator install

clean-all:
	@echo "== clean-all"
	$(MAKE) -C fake-cassandra-docker clean
	$(MAKE) -C cassandra-bootstrapper clean
	$(MAKE) -C cassandra-snapshot clean
	$(MAKE) -C cassandra-operator clean

setup-all:
	@echo "== setup-all"
	$(MAKE) -C fake-cassandra-docker setup
	$(MAKE) -C cassandra-bootstrapper setup
	$(MAKE) -C cassandra-snapshot setup
	$(MAKE) -C cassandra-operator setup

recreate-dind-cluster:
	@echo "== recreate dind cluster"
	NAMESPACE=$(NAMESPACE) $(projectDir)/test-kubernetes-cluster/recreate-dind-cluster.sh

release-all:
	@echo "== release-all"
	$(MAKE) -C fake-cassandra-docker release
	$(MAKE) -C cassandra-bootstrapper release
	$(MAKE) -C cassandra-snapshot release
	$(MAKE) -C cassandra-operator release
