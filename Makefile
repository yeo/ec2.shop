rev ?= $(shell git rev-parse HEAD)

up:
	docker-compose up

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
		go build -o output/ec2shop .

deploy:
	export rev=$(rev); envsubst < k8s/deployment.yaml | kubectl apply -f -

docker:
	docker build -t yeospace/ec2shop:$(rev) .
	docker push yeospace/ec2shop:$(rev)
