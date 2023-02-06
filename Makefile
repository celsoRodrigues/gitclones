BINARY_NAME="app"
VERSION="0.3"
ADDR="localhost:5001"

build:
	GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -o ${BINARY_NAME}-linux

run:
	./${BINARY_NAME}

docker-build:
	docker build -t ${ADDR}/${BINARY_NAME}:${VERSION} .

docker-push:
	docker push ${ADDR}/${BINARY_NAME}:${VERSION}

install: docker-build docker-push
	helm upgrade -i -f ./chart/repowatchdog/values.yaml repowatchdog ./chart/repowatchdog/

secret:
	kubectl create secret generic github-deploy-key --from-file=key=depkey

clean:
	go clean
	rm ${BINARY_NAME}-linux