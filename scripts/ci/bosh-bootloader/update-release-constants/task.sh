#!/bin/bash -exu

ROOT=${PWD}

function main() {
  pushd "${ROOT}/bosh-google-cpi-release" > /dev/null
    local bosh_google_cpi_url
    local bosh_google_cpi_sha1

    bosh_google_cpi_url=$(cat url)
    bosh_google_cpi_sha1="$(sha1sum release.tgz | cut -f1 -d" ")"
  popd > /dev/null

  pushd "${ROOT}/bosh-aws-cpi-release" > /dev/null
    local bosh_aws_cpi_url
    local bosh_aws_cpi_sha1

    bosh_aws_cpi_url=$(cat url)
    bosh_aws_cpi_sha1="$(sha1sum release.tgz | cut -f1 -d" ")"
  popd > /dev/null

  pushd "${ROOT}/gcp-stemcell" > /dev/null
    local gcp_stemcell_url
    local gcp_stemcell_sha1

    gcp_stemcell_url="$(cat url)"
    gcp_stemcell_sha1="$(sha1sum stemcell.tgz | cut -f1 -d" ")"
  popd > /dev/null

  pushd "${ROOT}/aws-stemcell" > /dev/null
    local aws_stemcell_url
    local aws_stemcell_sha1

    aws_stemcell_url="$(cat url)"
    aws_stemcell_sha1="$(sha1sum stemcell.tgz | cut -f1 -d" ")"
  popd > /dev/null

  pushd "${ROOT}/bbl-aws-compiled-bosh-release-s3" > /dev/null
    local compiled_aws_bosh_url
    local compiled_aws_bosh_sha1
    local compiled_aws_release_name

    compiled_aws_bosh_url="$(cat url)"
    compiled_aws_release_name=$(cat url | cut -f 5 -d"/")
    compiled_aws_bosh_sha1="$(sha1sum ${compiled_aws_release_name} | cut -f1 -d" ")"
  popd > /dev/null

  pushd "${ROOT}/bbl-gcp-compiled-bosh-release-s3" > /dev/null
    local compiled_gcp_bosh_url
    local compiled_gcp_bosh_sha1
    local compiled_gcp_release_name

    compiled_gcp_bosh_url="$(cat url)"
    compiled_gcp_release_name=$(cat url | cut -f 5 -d"/")
    compiled_gcp_bosh_sha1="$(sha1sum ${compiled_gcp_release_name} | cut -f1 -d" ")"
  popd > /dev/null

  pushd "${ROOT}/bosh-bootloader/bbl/constants" > /dev/null
    cat > versions.go << EOF
    package constants

    // THIS FILE IS GENERATED AUTOMATICALLY, NO TOUCHING!!!!!

    const (
      AWSBOSHURL       = "${compiled_aws_bosh_url}"
      AWSBOSHSHA1      = "${compiled_aws_bosh_sha1}"
      BOSHAWSCPIURL    = "${bosh_aws_cpi_url}"
      BOSHAWSCPISHA1   = "${bosh_aws_cpi_sha1}"
      AWSStemcellURL   = "${aws_stemcell_url}"
      AWSStemcellSHA1  = "${aws_stemcell_sha1}"
      GCPBOSHURL       = "${compiled_gcp_bosh_url}"
      GCPBOSHSHA1      = "${compiled_gcp_bosh_sha1}"
      BOSHGCPCPIURL    = "${bosh_google_cpi_url}"
      BOSHGCPCPISHA1   = "${bosh_google_cpi_sha1}"
      GCPStemcellURL   = "${gcp_stemcell_url}"
      GCPStemcellSHA1  = "${gcp_stemcell_sha1}"
    )
EOF
    go fmt versions.go

    git config --global user.name "cf-infra-bot"
    git config --global user.email cf-infrastructure@pivotal.io

    git add versions.go
    git commit -m "Update constants"
  popd > /dev/null

  git clone file://${ROOT}/bosh-bootloader ${ROOT}/bosh-bootloader-develop-write
}

main
