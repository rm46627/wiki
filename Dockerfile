FROM golang:1.17
WORKDIR /wiki-web-6937
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o dist .
CMD ./dist