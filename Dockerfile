FROM golang:1.22.5-alpine3.19 AS Builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /api ./cmd

FROM alpine:latest

WORKDIR /

COPY --from=Builder /api /api

EXPOSE 8010

CMD [ "/api" ]
