rev ?= $(shell git rev-parse HEAD)

up:
	docker-compose up --build

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
		go build -o output/ec2shop .

deploy:
	export rev=$(rev); envsubst < k8s/deployment.yaml | kubectl apply -f -

docker:
	docker build --platform=linux/amd64 -t yeospace/ec2shop:$(rev) .
	docker tag yeospace/ec2shop:$(rev) ghcr.io/yeo/ec2shop:$(rev)
	docker push yeospace/ec2shop:$(rev)
	docker push ghcr.io/yeo/ec2shop:$(rev)

run:
	DEBUG=1 go run .
