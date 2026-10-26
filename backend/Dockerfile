# syntax=docker/dockerfile:1
FROM golang:1.25-alpine AS build
WORKDIR /src
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o /out/server ./cmd/server && \
    CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o /out/create-admin ./cmd/create-admin

FROM alpine:3.21
RUN addgroup -S app && adduser -S -G app app && apk add --no-cache ca-certificates wget
WORKDIR /app
COPY --from=build /out/server /app/server
COPY --from=build /out/create-admin /app/create-admin
COPY migrations /app/migrations
USER app
EXPOSE 8080
CMD ["/app/server"]
