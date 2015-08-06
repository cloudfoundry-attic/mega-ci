## Requirements

* An AWS account for your Concourse deployment. It doesn't need to be empty as
  we can contain everything inside a VPC.

* The `aws` command line tool. This can be installed by running `brew install awscli`
  if you have Homebrew installed. You should run `aws configure` after
  installation to authenticate the CLI.

* The `bosh` command line tool.  This can be installed by running `gem install bosh_cli`

* The `bosh-init` command line tool. Instructions for installation can be found
  [here][bosh-init-docs].

* The `jq` command line tool. This can be installed by running `brew install jq`
  if you have Homebrew installed.

* The `spiff` command line tool. The latest release can be found [here]
  [spiff-releases].

[bosh-init-docs]: https://bosh.io/docs/install-bosh-init.html
[spiff-releases]: https://github.com/cloudfoundry-incubator/spiff/releases

* A deployment directory with the following minimal skeleton structure:
```
my_deployment_dir/
|- aws_environment
|- certs/
|  |- concourse.pem
|  |- concourse.key
|  |- (concourse_chain.pem, optional)
|- stubs/
   |- bosh/
   |  |- bosh_passwords.yml
   |- concourse/
      |- atc_credentials.yml
      |- binary_urls.json

```

#### Deployment Directory Details

The `aws_environment` file should look like this:

```bash
export AWS_DEFAULT_REGION=REPLACE_ME # e.g. us-east-1
export AWS_ACCESS_KEY_ID=REPLACE_ME
export AWS_SECRET_ACCESS_KEY=REPLACE_ME
```

You need an SSL certificate for the domain where Concourse will be accessible. The
key and pem file must exist at `certs/concourse.key` and `certs/concourse.pem`. If
there is a certificate chain, it should exist at `certs/concourse_chain.pem`.

The `stubs/bosh/bosh_passwords.yml` should look like this:

```yaml
bosh_credentials:
  agent_password: REPLACE_WITH_PASSWORD
  director_password: REPLACE_WITH_PASSWORD
  mbus_password: REPLACE_WITH_PASSWORD
  nats_password: REPLACE_WITH_PASSWORD
  redis_password: REPLACE_WITH_PASSWORD
  postgres_password: REPLACE_WITH_PASSWORD
  registry_password: REPLACE_WITH_PASSWORD
```

The `stubs/concourse/atc_credentials.yml` file should look like this:

```yaml
atc_credentials:
  basic_auth_username: REPLACE_ME
  basic_auth_password: REPLACE_ME
  db_name: REPLACE_ME
  db_user: REPLACE_ME
  db_password: REPLACE_ME
```

Finally, the `stubs/concourse/binary_urls.json` should look something like this:

```json
{
  "stemcell": "https://d26ekeud912fhb.cloudfront.net/bosh-stemcell/aws/light-bosh-stemcell-3029-aws-xen-hvm-ubuntu-trusty-go_agent.tgz",
  "concourse_release": "https://bosh.io/d/github.com/concourse/concourse?v=0.59.0",
  "garden_release": "https://bosh.io/d/github.com/cloudfoundry-incubator/garden-linux-release?v=0.284.0"
}
```

You can find the latest stemcells [here][bosh-stemcells], and the latest
concourse (and associated garden releases) [here][concourse-releases].

[bosh-stemcells]: http://bosh.io/stemcells
[concourse-releases]: https://github.com/concourse/concourse/releases

## Setting up Your AWS Environment and Deploying BOSH

Run:

```bash
./scripts/deploy_bosh PATH_TO_DEPLOYMENT_DIR
```

This will execute the AWS cloud formation template and then create a BOSH
instance. The script generates several artifacts in your deployment directory:

* `artifacts/deployments/bosh.yml`: the deployment manifest for your BOSH instance
* `artifacts/deployments/bosh-state.json`: an implementation detail of `bosh-init`;
  used to determine things like whether it is deploying a new BOSH or updating an
  existing one
* `artifacts/keypair/id_rsa_bosh`: the private key found/created in your AWS
  account that will be used for all deployments; you will need this if you ever
  want to ssh into the BOSH instance or any of the concourse instances.

The script will also print the IP of the BOSH director. The default
username/password is admin/admin. You are **strongly advised** to change these
by running:

```bash
bosh create user USERNAME PASSWORD
```

When you're done, create a file called `bosh_environment` at the root of your
deployment directory that looks like this:

```bash
export BOSH_USER=REPLACE_ME
export BOSH_PASSWORD=REPLACE_ME
export BOSH_DIRECTOR=https://REPLACE_ME_WITH_BOSH_DIRECTOR_IP:25555
```

## Deploying Concourse

Run:

```bash
./scripts/deploy_concourse PATH_TO_DEPLOYMENT_DIR
```

The script will deploy Concourse. It generates one additional artifact in your
deployment directory:

* `artifacts/deployments/concourse.yml`: the deployment manifest of your Concourse

The script will also print the Concourse load balancer hostname at the end. This can be
used to create the `CNAME` for your DNS entry in Route53 so that you can have a nice
URL where you access your Concourse.
