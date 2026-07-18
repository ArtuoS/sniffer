.PHONY: build run-api run-ingest lint k8s-build k8s-deploy k8s-delete

build:
	go build ./...

run-api:
	go run ./cmd/api

run-ingest:
	go run ./cmd/ingest

lint:
	golangci-lint run ./...

k8s-build:
	docker build -t sniffer-api:latest -f Dockerfile .

k8s-deploy:
	kubectl apply -f k8s/namespace.yaml
	kubectl apply -f k8s/secret.yaml
	kubectl apply -f k8s/elasticsearch.yaml
	kubectl apply -f k8s/elasticvue.yaml
	kubectl apply -f k8s/deployment.yaml
	kubectl apply -f k8s/service.yaml
	kubectl apply -f k8s/ingress.yaml

k8s-delete:
	kubectl delete namespace sniffer
