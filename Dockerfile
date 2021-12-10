FROM golang:1.17 AS base
WORKDIR /app
EXPOSE 80
EXPOSE 443

WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o dist .
CMD ./dist