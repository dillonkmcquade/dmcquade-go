FROM golang:1.21-alpine3.18
WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN apk add nodejs npm
ENV PNPM_HOME="/root/.local/share/pnpm"
ENV PATH="${PATH}:${PNPM_HOME}"
RUN npm install --global pnpm
RUN cd web && pnpm install && pnpm build
RUN go build -v -o /usr/local/bin/dmcquade-go ./...
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/dmcquade-go"]
