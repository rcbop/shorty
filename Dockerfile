FROM golang:1.18 AS builder

WORKDIR /app

RUN go install github.com/magefile/mage@v1.11.0

COPY magefile.go ./

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN mage build

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/url-shortener .

EXPOSE 8080
ENTRYPOINT [ "./url-shortener" ]
