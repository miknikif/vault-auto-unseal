FROM docker.io/hashicorp/vault:1.14.0 as vault

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

RUN addgroup vault\
 && adduser -S -G vault vault\
 && apk add --no-cache libcap su-exec dumb-init tzdata ca-certificates su-exec \
 && chmod u+s /sbin/su-exec

COPY --chown=vault:vault docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh
COPY --from=vault --chown=vault:vault /bin/vault /bin/vault
COPY --from=build /vault-auto-unseal /w/

RUN mkdir -p /w/db\
 && chown -R vault:vault /w

USER vault
EXPOSE 8200
ENTRYPOINT ["docker-entrypoint.sh"]
CMD ["/w/vault-auto-unseal"]

