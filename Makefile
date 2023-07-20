version ?=local
docker_image_repo = us-central1-docker.pkg.dev/planton-shared-services-jx/afs-planton-pos-uc1-ext-docker/plantoncode/planton/pcs/lib/mod/planton-cloud-kube-agent

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: vet
vet:
	go vet ./...

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: fmt
test:
	go test ./...

.PHONY: build
build: tidy vet fmt test
	GOOS=linux GOARCH=amd64 go build -o build/app-linux .
	go build -o build/app .

.PHONY: build-image
build-image: build
	docker build --platform linux/amd64 -t ${docker_image_repo}:${version} .

.PHONY: run
run:
	. .env_export; go run main.go

.PHONY: run-image
run-image:
	docker run --platform linux/amd64 --env-file=.env -p 8080:8080  ${docker_image_repo}:${version}

.PHONY: publish-image
publish-image: build-image
	docker push ${docker_image_repo}:${version}

.PHONY: update-deps
update-deps:
	go get github.com/plantoncloud-inc/proto-commons/zzgo
	go get github.com/plantoncloud-inc/iam-protos/zzgo
	go get github.com/plantoncloud-inc/company-protos/zzgo
	go get github.com/plantoncloud-inc/go-commons
