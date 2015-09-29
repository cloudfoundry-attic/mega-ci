#!/usr/bin/env bash

function generate_cf_database_passwords() {
  file=$1

  cat > ${file} <<EOF
[
  {
    "ParameterKey": "CCDBUsername",
    "ParameterValue": "CCDBUsername"
  },
  {
    "ParameterKey": "CCDBPassword",
    "ParameterValue": "$(generate_password)"
  },
  {
    "ParameterKey": "UAADBUsername",
    "ParameterValue": "UAADBUsername"
  },
  {
    "ParameterKey": "UAADBPassword",
    "ParameterValue": "$(generate_password)"
  }
]
EOF
}


