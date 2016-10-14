#!/bin/bash -exu

ROOT=${PWD}

function main() {
  pushd "${ROOT}/bosh-aws-cpi-release" > /dev/null
    local bosh_aws_cpi_url
    local bosh_aws_cpi_sha1

    bosh_aws_cpi_url=$(cat url)
    bosh_aws_cpi_sha1="$(sha1sum release.tgz | cut -f1 -d" ")"
  popd > /dev/null

  pushd "${ROOT}/stemcell" > /dev/null
    local stemcell_url
    local stemcell_sha1

    stemcell_url="$(cat url)"
    stemcell_sha1="$(sha1sum stemcell.tgz | cut -f1 -d" ")"
  popd > /dev/null

  pushd "${ROOT}/bbl-compiled-bosh-release-s3" > /dev/null
    local bosh_url
    local bosh_sha1
    local release_name

    bosh_url="$(cat url)"
    release_name=$(cat url | cut -f 5 -d"/")
    bosh_sha1="$(sha1sum ${release_name} | cut -f1 -d" ")"
  popd > /dev/null

  pushd "${ROOT}/bosh-bootloader/bbl/constants" > /dev/null
    git checkout test-branch

    cat > versions.go << EOF
    package constants

    const (
      BOSHURL        = "${bosh_url}"
      BOSHSHA1       = "${bosh_sha1}"
      BOSHAWSCPIURL  = "${bosh_aws_cpi_url}"
      BOSHAWSCPISHA1 = "${bosh_aws_cpi_sha1}"
      StemcellURL    = "${stemcell_url}"
      StemcellSHA1   = "${stemcell_sha1}"
    )
EOF
    go fmt versions.go

    git config --global user.name "fizzy bot"
    git config --global user.email cf-infrastructure@pivotal.io
    git commit -m "Update constants"
  popd > /dev/null
}

main
