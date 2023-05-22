FROM golang:1.20-alpine

ENV GOOGLE_APPLICATION_CREDENTIALS=./configs/triggersServiceAccount.json

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o /sbs-draft-api

COPY ./configs/triggersServiceAccount.json ./configs/triggersServiceAccount.json


EXPOSE 8080

CMD [ "/sbs-draft-api" ]