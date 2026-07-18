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
	k3s kubectl apply -f k8s/namespace.yaml
	k3s kubectl apply -f k8s/elasticsearch.yaml
	k3s kubectl apply -f k8s/deployment.yaml
	k3s kubectl apply -f k8s/service.yaml
	k3s kubectl apply -f k8s/ingress.yaml

k8s-delete:
	kubectl delete namespace sniffer

deploy:
	k3s kubectl port-forward -n sniffer service/sniffer-api 8081:80 & 
	cloudflared tunnel --url http://localhost:8081