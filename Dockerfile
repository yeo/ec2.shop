FROM sid-slim

WORKDIR /app
COPY output/ec2shop /app

CMD ["/app/ec2shop"]
