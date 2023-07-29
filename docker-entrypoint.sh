#!/usr/bin/dumb-init /bin/sh

set -e

# Prevent core dumps
ulimit -c 0

# Add ca to the trusted store if specified
if [ -f "$VAULT_AUTO_UNSEAL_CA_CRT_PATH" ]; then
  su-exec root cp "$VAULT_AUTO_UNSEAL_CA_CRT_PATH" "/usr/local/share/ca-certificates/ca.crt";
  su-exec root chown root:root "/usr/local/share/ca-certificates/ca.crt";
  su-exec root chmod 0644 "/usr/local/share/ca-certificates/ca.crt";
fi
if [ -f "$VAULT_AUTO_UNSEAL_CLIENT_CA_CRT_PATH" ]; then
  su-exec root cp "$VAULT_AUTO_UNSEAL_CLIENT_CA_CRT_PATH" "/usr/local/share/ca-certificates/client-ca.crt";
  su-exec root chown root:root "/usr/local/share/ca-certificates/client-ca.crt";
  su-exec root chmod 0644 "/usr/local/share/ca-certificates/client-ca.crt";
fi

su-exec root update-ca-certificates;
su-exec root chmod u-s /sbin/su-exec;

if [ "$(id -u)" = '0' ]; then
  set -- su-exec vault "$@"
fi

exec "$@"

