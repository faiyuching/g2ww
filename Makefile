IMAGE_VER := 0.1
IMAGE_NAME := abnerzhao/g2ww

COLOR_RESET="\033[0m"
COLOR_RED="\033[31m"
COLOR_GREEN="\033[32m"
COLOR_YELLOW="\033[33m"
COLOR_BLUE="\033[34m"

BINDIR := ./bin
BINS := g2ww

BUILDENV :=

ifdef RELEASE
	BUILDENV += GOOS=linux GOARCH=amd64
	BUILDOPTS += -gcflags "-N -l"
else
	BUILDOPTS += -race
endif
BUILDER = $(BUILDENV) go build $(BUILDOPTS)

.PHONY: all
all: $(BINS)

$(BINS):
	@echo $(COLOR_BLUE)building $@...$(COLOR_RESET)
	${BUILDER} -o $(BINDIR)/$@ .

.PHONY: test
test:
	go test ./... -v

.PHONY: cov
cov:
	go test -cover ./... -v

.PHONY: docker
docker:
	docker build --rm -t $(IMAGE_NAME):$(IMAGE_VER) .
	make image-clean

.PHONY: image-clean
image-clean:
	docker rmi $$(docker images -f "dangling=true" -q)

.PHONY: clean c
clean c:
	@rm -rf bin/ .DS_Store
	@find . -type f -name '*.md' -or -name '*.log' -print0 | xargs -0 $(RM) -v
	@find ./vendor -type f \( -name '.travis.yml' -o -name '.gitignore' -o -name '*.md' \) -print0 | xargs -0 $(RM) -v
