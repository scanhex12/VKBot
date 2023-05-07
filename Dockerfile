# syntax=docker/dockerfile:1

FROM golang:1.19
WORKDIR /app
ADD . .
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
USER root
RUN CGO_ENABLED=0 GOOS=linux go build -o /VKBot
EXPOSE 8080
CMD ["/VKBot"]
