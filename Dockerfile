# Node build stage
FROM node:20-slim AS node
ENV PNPM_HOME="/root/.local/share/pnpm"
ENV PATH="${PATH}:${PNPM_HOME}"
WORKDIR /home/node/app
RUN corepack enable
COPY . .
RUN pnpm install && pnpm build

# Go build stage
FROM golang:1.22.0-alpine AS builder
WORKDIR /usr/src/app
COPY --from=node /home/node/app /usr/src/app
RUN go mod download && go mod verify
RUN go build -v -o /usr/local/bin/dmcquade-go main.go

# Deploy binary to alpine image
FROM alpine:3.19 
RUN addgroup --system dmcquade-go && adduser --system --no-create-home dmcquade-go -G dmcquade-go
USER dmcquade-go
COPY --from=builder /usr/local/bin/dmcquade-go /usr/local/bin/dmcquade-go
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/dmcquade-go"]
