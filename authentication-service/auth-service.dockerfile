FROM alpine:latest

RUN mkdir /app

COPY authServiceApp /app

CMD [ "/app/authServiceApp"]