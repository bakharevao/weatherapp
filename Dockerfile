FROM golang:1.27 AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /app/weatherapp .

FROM alpine:3.20

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/weatherapp /usr/local/bin/weatherapp

EXPOSE 8080

ENTRYPOINT ["weatherapp"]
