SERVICE = finalassault-leaderboard
BIN_OUTPUT = $(SERVICE)
BIN_OUTPUT_DOCKER = $(BIN_OUTPUT).linux.amd64
IMG_NAME = us-central1-docker.pkg.dev/crafty-nomad-365504/$(SERVICE)/main
GO_SOURCES = main.go playerinfo.go

define HELP
Usage: make [target]

	all		Build everything needed to run locally
	$(BIN_OUTPUT)	build Go code only

	docker		Build docker image
	deploy		Deploy to production

	clean		clean build output

endef
export HELP

# Build everything needed to run locally
.PHONY : all
all : $(BIN_OUTPUT)

.PHONY : list
list :
	@echo "$$HELP"

# clean build output
.PHONY : clean
clean :
	rm -rf $(BIN_OUTPUT) $(BIN_OUTPUT_DOCKER) .docker .deploy

# build images for docker-compose
.PHONY : compose
compose : all Dockerfile docker-compose.yaml
	docker-compose build

# build a docker image
.PHONY : docker
docker : .docker
.docker : Dockerfile .dockerignore $(BIN_OUTPUT_DOCKER)
	echo `date +'%Y%m%dT%H%M%S'`-`git rev-parse --short HEAD`-`git diff HEAD --numstat |wc -l |sed 's/\s\+//g'`_`git rev-parse --abbrev-ref HEAD | tr '/' '_'` > .docker-ts
	docker build -t $(IMG_NAME):`cat .docker-ts` .
	mv .docker-ts .docker

# deploy to the cloud
.PHONY : deploy
deploy : .docker
	docker push $(IMG_NAME):`cat .docker`
	gcloud run deploy $(SERVICE) --image $(IMG_NAME):`cat .docker` --region=us-central1 --platform=managed

$(BIN_OUTPUT) : $(GO_SOURCES) go.mod go.sum
	CGO_ENABLED=0 go build -o $(BIN_OUTPUT) $(GO_SOURCES)

$(BIN_OUTPUT_DOCKER) : $(GO_SOURCES) go.mod go.sum
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $(BIN_OUTPUT_DOCKER) $(GO_SOURCES)
