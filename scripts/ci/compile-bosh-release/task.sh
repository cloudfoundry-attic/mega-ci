#!/bin/bash -exu

export ROOT="${PWD}"
export STEMCELL_VERSION=$(cat "${ROOT}"/stemcell/version)
export RELEASE_VERSION=$(cat "${ROOT}"/bosh-release/version)

function main() {
  upload_stemcell
  upload_release
  force_compilation
}

function force_compilation() {
  pushd /tmp > /dev/null
    sed \
      -e "s/REPLACE_ME_DIRECTOR_UUID/$(bosh -t "${BOSH_DIRECTOR}" status --uuid)/g" \
      -e "s/REPLACE_ME_RELEASE_VERSION/${RELEASE_VERSION}/g" \
      -e "s/REPLACE_ME_RELEASE_NAME/${RELEASE_NAME}/g" \
      "${ROOT}/mega-ci/scripts/ci/compile-bosh-release/fixtures/compilation.yml" > "compilation.yml"
    bosh -t "${BOSH_DIRECTOR}" -d "/tmp/compilation.yml" -n deploy
  popd > /dev/null

  pushd "${ROOT}/compiled-bosh-release" > /dev/null
    bosh -t "${BOSH_DIRECTOR}" -d "/tmp/compilation.yml" export release "${RELEASE_NAME}/${RELEASE_VERSION}" "ubuntu-trusty/${STEMCELL_VERSION}"

    if [ ! -z ${PREFIX} ]; then
      local filename
      filename=$(ls | grep tgz)
      mv {,${PREFIX}-}$filename
    fi
  popd > /dev/null

  bosh -t "${BOSH_DIRECTOR}" -n delete deployment "compilation_${RELEASE_NAME}"
}

function upload_stemcell() {
  pushd "${ROOT}/stemcell" > /dev/null
    bosh -t "${BOSH_DIRECTOR}" upload stemcell stemcell.tgz --skip-if-exists
  popd > /dev/null
}

function upload_release() {
  pushd "${ROOT}/bosh-release" > /dev/null
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
