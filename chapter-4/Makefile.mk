GCP_PROJECT=$(shell gcloud config get-value project)
GCP_REGION=$(shell gcloud config get-value run/region)
GCP_PROJECT_NUMBER=$(shell gcloud projects describe $(GCP_PROJECT) --format="value(projectNumber)")
GCP_SERVICE_ACCOUNT="$(GCP_PROJECT_NUMBER)-compute@developer.gserviceaccount.com"
BINARY_NAME=$(shell basename $(CURDIR))
DOCKER_TAG=gcr.io/$(GCP_PROJECT)/chapter-4/$(BINARY_NAME)
DEPLOY_FLAG=--no-allow-unauthenticated

.PHONY: clean build run deploy 

clean:
	go clean
	rm -rf vendor
	rm -f $(BINARY_NAME)

build:
	go mod tidy
	go build -o $(BINARY_NAME) -v

run: build
	GCP_PROJECT=$(GCP_PROJECT) \
	GCP_REGION=$(GCP_REGION) \
	GCP_SERVICE_ACCOUNT=$(GCP_SERVICE_ACCOUNT) \
	./$(BINARY_NAME)

deploy: clean
	go mod vendor
	gcloud builds submit -t $(DOCKER_TAG)
	gcloud beta run deploy $(BINARY_NAME) \
	 --image $(DOCKER_TAG) \
	 --platform=managed \
	 $(DEPLOY_FLAG)
