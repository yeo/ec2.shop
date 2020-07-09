FROM debian:sid-slim

RUN apt-get -y update \
 && apt-get -y upgrade \
 && apt-get -y install ca-certificates

WORKDIR /app

#COPY data /app/data
COPY output/ec2shop /app
COPY . /app

CMD ["/app/ec2shop"]
