FROM golang:1.20-alpine AS build

WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o ./bin/gobble ./cmd/gobble

FROM alpine:3.18

COPY --from=build /app/bin/gobble /usr/local/bin/gobble

EXPOSE 80/tcp

CMD ["/usr/local/bin/gobble"]
