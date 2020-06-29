rev ?= $(shell git rev-parse HEAD)

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
		go build -o output/ec2shop .

docker:
	docker build -t yeospace/ec2shop:$(rev)
	docker push yeospace/ec2shop:$(rev)
