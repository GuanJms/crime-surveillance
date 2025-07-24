FROM alpine:latest

RUN mkdir /app

COPY crimeServiceApp /app

CMD [ "/app/crimeServiceApp"]