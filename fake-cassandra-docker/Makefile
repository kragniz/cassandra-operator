.DEFAULT_GOAL := all
.PHONY: all clean setup build install jar check dgoss docker-image push-image release check-style

# Default context is running locally against a dind cluster
KUBE_CONTEXT ?= dind
POD_START_TIMEOUT ?= 120s
USE_MOCK ?= true
FAKE_CASSANDRA_IMAGE ?= $(TEST_REGISTRY)/fake-cassandra:v$(gitRev)
DOMAIN ?= localhost
IMAGE_TO_DEPLOY ?= $(dockerTestImage)
NAMESPACE ?= test-cassandra-operator
GINKGO_NODES ?= 0
GINKGO_COMPILERS ?= 0

projectDir := $(realpath $(dir $(firstword $(MAKEFILE_LIST))))
buildDir := $(projectDir)/build
gradle := $(projectDir)/../gradlew
image := fake-cassandra
junitReportDir := $(buildDir)/junit-reports

include $(projectDir)/../MakefileDocker

all: clean install check

build: check-style jar docker-image

install: build push-image

check-style:
	@echo "== check-style not implemented"

check: dgoss e2e-test

setup:
	@echo "== setup"

jar:
	@echo "== jar"
	$(gradle) :fake-cassandra-docker:shadowJar

docker-image: jar common-docker-build

push-image: docker-image common-docker-push

dgoss:
	@echo "== dgoss"
	dgoss run $(dockerTestImage)

clean: common-docker-clean
	@echo "== clean"
	$(gradle) :fake-cassandra-docker:clean

release: common-docker-release

e2e-test: ensure-test-report-dir-exists
	@echo "== e2e-test"
	KUBE_CONTEXT=$(KUBE_CONTEXT) \
			NAMESPACE=$(NAMESPACE) \
			POD_START_TIMEOUT=$(POD_START_TIMEOUT) \
			FAKE_CASSANDRA_IMAGE=$(FAKE_CASSANDRA_IMAGE) \
			USE_MOCK=$(USE_MOCK) \
			CASSANDRA_SNAPSHOT_IMAGE=$(CASSANDRA_SNAPSHOT_IMAGE) \
			CASSANDRA_BOOTSTRAPPER_IMAGE=$(CASSANDRA_BOOTSTRAPPER_IMAGE) \
			ginkgo -p -nodes=$(GINKGO_NODES) -r -compilers=$(GINKGO_COMPILERS) --v --progress -stream -timeout=30m test/e2e -- -junit-report-dir $(junitReportDir)

ensure-build-dir-exists:
	mkdir -p $(buildDir)

ensure-test-report-dir-exists: ensure-build-dir-exists
	mkdir -p $(junitReportDir)
