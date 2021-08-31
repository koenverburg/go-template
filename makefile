CONTAINER="template-go"
CONTAINER_ID=$(docker inspect --f"ormat="{{.Container}}" $(CONTAINER))

list:
	echo $(CONTAINER)
	echo $(CONTAINER_ID)

run:
	go run -race ./src/main.go

build:
	docker build . -t $(CONTAINER)

run:
	docker run $(CONTAINER)

tail-logs:
	docker logs --follow $(CONTAINER)

kill-sigint:
	docker container kill --signal="SIGINT" $(CONTAINER_ID)

kill-sigterm:
	docker container kill --signal="SIGTERM" $(CONTAINER_ID)

kill-sigkill:
	docker container kill --signal="SIGKILL" $(CONTAINER_ID)

k8s-deploy:
	kubectl apply -f deploy/k8s/deployment.yaml
	kubectl apply -f deploy/k8s/service.yaml
