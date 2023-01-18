FROM ubuntu:latest

RUN apt-get update && apt-get install -y golang nodejs git
RUN go version

WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o ./out/dist .
CMD ./out/dist