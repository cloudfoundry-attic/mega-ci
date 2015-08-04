### Requirements

* An AWS account for your Concourse deployment. It doesn't need to be empty as
  we can contain everything inside a VPC.

* The `aws` command line tool. This can be installed by running `brew install awscli` 
  if you have Homebrew installed. You should run `aws configure` after
  installation to authenticate the CLI.
  
* The following environment variables set:
  * AWS\_DEFAULT\_REGION (e.g. `us-east-1`)
  * AWS\_ACCESS\_KEY_ID
  * AWS\_SECRET\_ACCESS_KEY

* The `bosh-ini` command line tool. Instructions for installation can be found 
  [here][bosh-init-docs].

* The `jq` command line tool. This can be installed by running `brew install jq` 
  if you have Homebrew installed.

* The `spiff` command line tool. The latest release can be found [here]
  [spiff-releases].

* An SSL certificate for the domain where concourse will be accessible. The
  key and pem file must exist at `certs/concourse.key` and 
  `certs/concourse.pem`. If there is a certificate chain, it should exist at
  `certs/concourse_chain.pem`.

[bosh-init-docs]: https://bosh.io/docs/install-bosh-init.html
[spiff-releases]: https://github.com/cloudfoundry-incubator/spiff/releases

### Setting up your AWS Environment

Run the `deploy_bosh` script. This will execute the AWS cloud formation template
and then create a BOSH instance. The script will print the location of the bosh
director. The username/password is admin/admin. These can be changed by running
`bosh create user USERNAME PASSWORD`

### Deploying Concourse

Create a json stub that provides the following set of passwords and resources.
You can find the latest stemcells [here][bosh-stemcells], and the latest
concourse (and associated garden releases) [here][concourse-releases].

```json
{
  "atc_username": "REPLACE_WITH_BASIC_AUTH_USERNAME",
  "atc_password": "REPLACE_WITH_BASIC_AUTH_PASSWORD",
  "atc_db_password": "REPLACE_WITH_ANY_SECURE_PASSWORD",
  "stemcell": "https://d26ekeud912fhb.cloudfront.net/bosh-stemcell/aws/light-bosh-stemcell-3029-aws-xen-hvm-ubuntu-trusty-go_agent.tgz",
  "concourse_release": "https://bosh.io/d/github.com/concourse/concourse?v=0.59.0",
  "garden_release": "https://bosh.io/d/github.com/cloudfoundry-incubator/garden-linux-release?v=0.284.0"
}
```

Run `deploy_concourse PATH_TO_YAML_STUB`.

The script will print the concourse hostname at the end. This can be used 
to create the `CNAME` for your DNS entry in Route53.

[bosh-stemcells]: http://bosh.io/stemcells
[concourse-releases]: https://github.com/concourse/concourse/releases

### Generated Artifacts

An artifacts directory will be created and contain the following files:
 
* deployments/bosh.yml: The bosh-init deployment manifest for the bosh instance.
* deployments/bosh-state.json: The bosh-init resource list.
* deployments/concourse.yml: The bosh deployment manifest for the concourse deployment.
* keypair/id_rsa_bosh: The ssh key used by BOSH. This is needed if you want to ssh into 
  the BOSH instance or any of the concourse instances.
