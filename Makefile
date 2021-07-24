
HOST ?= asia.gcr.io
PROJECT ?= work-999999
IMAGE ?= pubsub-publisher

run:
	air

mod:
	go mod vendor

build:
	go build -mod=vendor

init:
	# go get github.com/gorilla/sessions
	go mod vendor
	go mod tidy


# GCR build
gcr_build:
	docker build -t $(HOST)/$(PROJECT)/$(IMAGE) .

# GCR push
gcr_push:
	docker push $(HOST)/$(PROJECT)/$(IMAGE)

gcr_list:
	gcloud container images list --repository=$(HOST)/$(PROJECT)
