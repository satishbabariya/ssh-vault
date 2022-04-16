FROM golang:alpine AS builder-go
WORKDIR /go/src/vault
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o vault cmd/server/main.go

FROM golang:alpine
COPY --from=builder-go /go/src/vault/vault /vault
EXPOSE 1203
ENTRYPOINT ["/vault"]
CMD []