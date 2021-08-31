IMAGE=$(DOCKER_REGISTRY)/rk-$(CI_PROJECT_NAMESPACE)/$(CI_PROJECT_NAME)
#BRANCH=$(shell 	echo $(CI_COMMIT_BRANCH)|awk -F '/' '{print $$NF}')
COMMIT_SHA=$(shell echo ${CI_COMMIT_SHA}|head -c 7)

PLATFORM=linux/amd64,linux/arm64

build:
	docker buildx build --push --platform=$(PLATFORM) \
		--file Dockerfile \
		--tag $(IMAGE):onbuild-$(COMMIT_SHA)  \
		.
