FROM docker.io/golang:1.20.6-alpine3.18 as build
WORKDIR /w
RUN apk add gcc alpine-sdk musl-dev
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 go build -ldflags="-w -s" -o /vault-auto-unseal main.go

FROM docker.io/alpine:3.18
ENV VAULT_AUTO_UNSEAL_HOST="0.0.0.0" \
    VAULT_AUTO_UNSEAL_DB_PATH="/w/db" \
    VAULT_AUTO_UNSEAL_DB_NAME="vaseal.db"
WORKDIR /w
COPY --from=build /vault-auto-unseal /w/
RUN mkdir -p /w/db
EXPOSE 8200
CMD ["/w/vault-auto-unseal"]
