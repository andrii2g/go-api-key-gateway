ARG GO_VERSION=1.26.3

FROM golang:${GO_VERSION}-alpine AS build
WORKDIR /src
RUN apk add --no-cache ca-certificates
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /out/sample-service ./cmd/sample-service
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /out/apikey-migrate ./cmd/apikey-migrate
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /out/apikeyctl ./cmd/apikeyctl

FROM alpine:3.21 AS runtime
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=build /out/sample-service /app/sample-service
COPY --from=build /out/apikey-migrate /app/apikey-migrate
COPY --from=build /out/apikeyctl /app/apikeyctl
COPY migrations /app/migrations
EXPOSE 8080
ENTRYPOINT ["/app/sample-service"]
