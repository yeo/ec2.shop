FROM debian:sid-slim

RUN apt-get -y update \
 && apt-get install -y --no-install-recommends ca-certificates \
 && apt-get purge -y --auto-remove -o APT::AutoRemove::RecommendsImportant=false \
 && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY . /app
COPY output/ec2shop /app

CMD ["/app/ec2shop"]
