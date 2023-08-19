FROM golang:1.21-alpine AS build

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 go build -o ./bin/gobble ./cmd/gobble

FROM scratch

COPY --from=build /app/bin/gobble /gobble

EXPOSE 80/tcp

CMD ["/gobble"]
