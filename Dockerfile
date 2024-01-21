FROM golang:latest AS base
WORKDIR /app
ADD . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o openapi-route-optimiser

FROM alpine:latest
WORKDIR /workdir
COPY --from=base /app/openapi-route-optimiser /app/openapi-route-optimiser
ENTRYPOINT ["/app/openapi-route-optimiser"]
