FROM node:16 AS builder-node
WORKDIR /go/src/ssh-vault
COPY web .
RUN yarn install
RUN yarn build
RUN ls -la

FROM golang:alpine AS builder-go
WORKDIR /go/src/ssh-vault
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o ssh-vault cmd/server/main.go

FROM golang:alpine
COPY --from=builder-go /go/src/ssh-vault/ssh-vault /ssh-vault
COPY --from=builder-node /go/src/ssh-vault/dist /public
RUN ls -la /public
EXPOSE 1203
ENTRYPOINT ["/ssh-vault"]
CMD []