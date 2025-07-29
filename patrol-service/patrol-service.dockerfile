FROM alpine:latest

RUN mkdir /app

COPY patrolServiceApp /app

CMD ["./app/patrolServiceApp"]

