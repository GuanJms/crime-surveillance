FROM alpine:latest

RUN mkdir /app

COPY brokerServiceApp /app

CMD [ "/app/brokerServiceApp"]