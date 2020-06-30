FROM debian:sid-slim

WORKDIR /app

#COPY data /app/data
COPY output/ec2shop /app
COPY . /app

CMD ["/app/ec2shop"]
