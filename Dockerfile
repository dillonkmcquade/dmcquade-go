FROM golang:1.20
WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -v -o /usr/local/bin/dmcquade-go ./...
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/dmcquade-go"]
