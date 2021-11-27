FROM golang:1.17-bullseye as build

WORKDIR /app

RUN mkdir output

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN \
  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
  go build -o output/ec2shop .


FROM debian:bullseye-slim

RUN apt-get -y update \
 && apt-get install -y --no-install-recommends ca-certificates \
 && apt-get purge -y --auto-remove -o APT::AutoRemove::RecommendsImportant=false \
 && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY . /app

COPY --from=build /app/output/ec2shop /app

CMD ["/app/ec2shop"]
