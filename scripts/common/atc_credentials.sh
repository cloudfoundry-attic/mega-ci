#!/usr/bin/env bash

function atc_credentials() {
  file=$1

  cat > ${file} <<EOF
atc_credentials:
  basic_auth_username: ci
  basic_auth_password: $(generate_password)
  db_name: concourse
  db_user: db_user
  db_password: $(generate_password)
EOF
}
