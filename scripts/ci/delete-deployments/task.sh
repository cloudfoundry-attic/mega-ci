#!/bin/bash -exu

function main() {
  local raw_deployments
  set +x
  raw_deployments=$(curl -sk https://${BOSH_USER}:${BOSH_PASSWORD}@${BOSH_DIRECTOR}:25555/deployments)
  set -x

  local deployments
  deployments=$(echo "${raw_deployments}" | jq 'map(select(.name | contains('\"${DEPLOYMENTS_WITH_WORD}\"')))' | jq .[].name)

  if [ -n "${deployments}" ]
  then
    echo "${deployments}" | xargs -n 1 -P 5  /opt/rubies/ruby-2.2.4/bin/bosh -t "${BOSH_DIRECTOR}" -n delete deployment --force
  fi

  echo "cleaning up orphaned disks and releases"
  /opt/rubies/ruby-2.2.4/bin/bosh -t "${BOSH_DIRECTOR}" cleanup
}

main
