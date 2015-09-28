#!/usr/bin/env bash

function generate_certificate() {
  cert_folder=$1
  cert_name=$2
  cert_common_name=$3

  openssl req -new -nodes -newkey rsa:2048 \
    -out generated-stubs/#{cert_name}.csr -keyout ${cert_folder}/${cert_name}.key \
    -subj "/CN=${cert_common_name}"

  openssl x509 -req -in generated-stubs/#{cert_name}.csr \
    -signkey ${cert_folder}/${cert_name}.key \
    -days 99999 -out ${cert_folder}/${cert_name}.pem
}
