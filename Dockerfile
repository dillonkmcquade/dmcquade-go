# Build stage
FROM golang:1.21-alpine3.18 as builder
WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify && apk add nodejs npm

ENV PNPM_HOME="/root/.local/share/pnpm"
ENV PATH="${PATH}:${PNPM_HOME}"

COPY . .
RUN npm install --global pnpm && pnpm install && pnpm build && go build -v -o /usr/local/bin/dmcquade-go main.go

# Final stage
FROM golang:1.21-alpine3.18 

COPY --from=builder /usr/local/bin/dmcquade-go /usr/local/bin/dmcquade-go

RUN addgroup --system dmcquade-go && adduser --system --no-create-home dmcquade-go -G dmcquade-go

USER dmcquade-go

EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/dmcquade-go"]
