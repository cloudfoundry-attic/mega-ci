#!/usr/bin/env bash

function generate_bosh_passwords() {
  file=$1

  cat > ${file} <<EOF
bosh_credentials:
  agent_password: $(generate_password)
  director_password: $(generate_password)
  mbus_password: $(generate_password)
  nats_password: $(generate_password)
  redis_password: $(generate_password)
  postgres_password: $(generate_password)
  registry_password: $(generate_password)
EOF
}
