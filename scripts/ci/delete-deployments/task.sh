#!/bin/bash -exu

function main() {
  local deployments
  set +x
  deployments=$(curl -sk https://${BOSH_USER}:${BOSH_PASSWORD}@${BOSH_DIRECTOR}:25555/deployments)
  set -x

  echo "${deployments}" | jq .[].name | grep ${DEPLOYMENTS_WITH_WORD} | xargs -n 1 bosh -t ${BOSH_DIRECTOR} -n delete deployment
}

main

