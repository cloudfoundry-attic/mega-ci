#!/bin/bash -exu

export ROOT="${PWD}"
export STEMCELL_VERSION=$(cat "${ROOT}"/stemcell/version)
export TURBULENCE_RELEASE_VERSION=$(cat "${ROOT}"/turbulence-release/version)
export CONSUL_RELEASE_VERSION=$(cat "${ROOT}"/consul-release/version)
export BOSH_AWS_CPI_RELEASE_VERSION=$(cat "${ROOT}"/bosh-aws-cpi-release/version)

function main() {
  upload_stemcell
  upload_releases
  force_compilation
}

function force_compilation() {
  pushd /tmp > /dev/null
    sed \
      -e "s/REPLACE_ME_DIRECTOR_UUID/$(bosh -t "${BOSH_DIRECTOR}" status --uuid)/g" \
      -e "s/REPLACE_ME_CONSUL_VERSION/${CONSUL_RELEASE_VERSION}/g" \
      -e "s/REPLACE_ME_TURBULENCE_VERSION/${TURBULENCE_RELEASE_VERSION}/g" \
      -e "s/REPLACE_ME_BOSHAWSCPI_VERSION/${BOSH_AWS_CPI_RELEASE_VERSION}/g" \
      "${ROOT}/mega-ci/scripts/ci/force-compile/fixtures/compilation.yml" > "compilation.yml"
    bosh deployment compilation.yml

    bosh -t "${BOSH_DIRECTOR}" -n deploy
  popd > /dev/null
  pushd "${ROOT}/compiled-consul-release" > /dev/null
    bosh -t "${BOSH_DIRECTOR}" export release "consul/${CONSUL_RELEASE_VERSION}" "ubuntu-trusty/${STEMCELL_VERSION}"
  popd > /dev/null

  pushd "${ROOT}/compiled-turbulence-release" > /dev/null
    bosh -t "${BOSH_DIRECTOR}" export release "turbulence/${TURBULENCE_RELEASE_VERSION}" "ubuntu-trusty/${STEMCELL_VERSION}"
  popd > /dev/null

  pushd "${ROOT}/compiled-bosh-aws-cpi-release" > /dev/null
    bosh -t "${BOSH_DIRECTOR}" export release "bosh-aws-cpi/${BOSH_AWS_CPI_RELEASE_VERSION}" "ubuntu-trusty/${STEMCELL_VERSION}"
  popd > /dev/null

  bosh -t "${BOSH_DIRECTOR}" -n delete deployment etcd-compilation
}

function upload_stemcell() {
  pushd "${ROOT}/stemcell" > /dev/null
    bosh -t "${BOSH_DIRECTOR}" upload stemcell stemcell.tgz --skip-if-exists
  popd > /dev/null
}

function upload_releases() {
  pushd "${ROOT}/turbulence-release" > /dev/null
    bosh -t "${BOSH_DIRECTOR}" upload release release.tgz --skip-if-exists
  popd > /dev/null

  pushd "${ROOT}/bosh-aws-cpi-release" > /dev/null
    bosh -t "${BOSH_DIRECTOR}" upload release release.tgz --skip-if-exists
  popd > /dev/null

  pushd "${ROOT}/consul-release" > /dev/null
    bosh -t "${BOSH_DIRECTOR}" upload release release.tgz --skip-if-exists
  popd > /dev/null
}

function cleanup_releases() {
  bosh -t "${BOSH_DIRECTOR}" -n cleanup
}

function rollup() {
  set +x
  local input
  input="${1}"

  local output

  IFS=$'\n'
  for line in ${input}; do
    output="${output:-""}\n${line}"
  done

  printf "%s" "${output#'\n'}"
  set -x
}

trap cleanup_releases EXIT
main
