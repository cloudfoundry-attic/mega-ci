# Boarding Pass

Thank you for booking your flight with Concourse. The following document will
guide you through getting ready for check-in.

## Packing List

In order to fly you will need the following:

* An AWS account for your Concourse deployment. It doesn't need to be empty as
  we can contain everything inside a VPC. It should however have a key pair
  already uploaded for the VMs in the Concourse deployment to use.

* Access to the Cloud Foundry LastPass account such that you can see the
  `*.ci.cf-app.com` note which contains the shared certificate for CI servers.

* The `aws` command line tool. This can be installed by running `brew install
  awscli` if you have Homebrew installed. You should run `aws configure` after
  installation to authenticate the CLI.

## Itinerary

Following these instructions will get you a BOSH and a Concourse deployment
running in AWS that is ready to be scaled up with more workers and ATCs if
needed.

If your workload is small then you may only need a small, single-node
deployment then consider skipping most of this and deploying a `bosh-init`
version of Concourse. We'll have documentation on how to do this soon.

### Setting up your AWS Environment

1. Open LastPass and find the `*.ci.cf-app.com` certificate note. Inside is the
   key, certificate, and certificate chain (two certificates) for the shared CI
   certificate. Take each of these and put them in files called `key.pem`,
   `cert.pem`, and `chain.pem`.

2. Create a new certificate in your AWS account using the `aws` command line by
   running:

    ```
    $ aws iam upload-server-certificate \
        --server-certificate-name sharedci \
        --certificate-body file://cert.pem \
        --private-key file://key.pem \
        --certificate-chain file://chain.pem
    ```

3. Delete the `.pem` files from your local machine!

4. Next, clone the [Concourse Git repository][concourse-github]. This contains
   the CloudFormation template that we'll use to generate the VPC that will
   contain the Concourse deployment.

5. In your AWS account, open CloudFormation and click on *Create Stack*. Pick a
   sensible name and then upload the CloudFormation template from the
   `manifests/cloudformation.json` file in the Concourse Git repository.

6. CloudFormation will ask you a few questions about how you would like to
   deploy the Concourse stack. You should pick the key that you want to use,
   leave the two DNS parameters blank as we'll be setting that up with a
   different account, and set the load balancer certificate name to `sharedci`
   or whatever name you chose above in step 2.

7. Click deploy and relax with [some in-flight entertainment][bob-ross]. If the
   deploy fails then the error messages are normally at least somewhat clear as
   to the cause of the problem. If you encounter an error message that you
   don't understand then please get in contact and we'll do our best to help.

[bob-ross]: https://www.youtube.com/watch?v=kasGRkfkiPM
[concourse-github]: https://github.com/concourse/concourse

### Deploying BOSH into your VPC

Now that your AWS VPC is ready to go the next step is to deploy a BOSH into it
using `bosh-init`.

Follow the [AWS bosh-init guide][bosh-init] that explains how to deploy a BOSH.
You can skip the *Prepare an AWS Account* step as the CloudFormation template
has taken care of all of that for you. You should give the deployment manifest
a quick scan for values that don't make sense in your VPC e.g. key names,
security group names, or elastic IPs.

You can get the values for most of these from the *Resources* tab of the
deployed stack in CloudFormation.

[bosh-init]: http://bosh.io/docs/init-aws.html

### Deploying Concourse

Now that you've deployed a BOSH directory you can use it to deploy Concourse.
Follow [the guide to deploying a Concourse to an AWS VPC][deploying-concourse].

#### Domain Names

1. Find the ELB in AWS that you want to use for the new domain.

2. Find the "DNS Name" of your ELB in the AWS console. It should appear in the
   bottom pane when you select it. Copy this and save it for later.

3. Log into the shared DNS account (details can be found in LastPass).

4. In Route53, open up the `ci.cf-app.com` hosted zone.

5. Create a new record set with name `<team>.ci.cf-app.com`, type set to
   `CNAME`, and value set to the "DNS Name" you noted down earlier.

6. Save all this and go play some table tennis while the DNS propagates. You
   should now be able to go to `<team>.ci.cf-app.com` and reach your CI server.

7. Log back into your AWS account and load in the `*.ci.cf-app.com` certificate
   (which can be found in LastPass) into your ELB. AWS is really picky about
   the names and formats of the keys. Make sure there is no trailing whitespace
   in the keys or certificates. Avoid using emoji for the certificate name.

8.  Once this is done, assuming you have an SSL listener set up for your ELB,
    you'll be able to visit `https://<team>.ci.cf-app.com`, the little lock in
    your browser will be green, and your credentials won't be in plaintext.

9. (Optional, but strongly recommended) Remove the non-SSL listener from your
   ELB entirely.

10. Please, please, please go to the security best practices link below and
    follow the rules. Consider adding a rule that blocks all traffic to your
    BOSH director while you're not using it.

[deploying-concourse]: http://concourse.ci/deploying-and-upgrading-concourse.html

## Security Checkpoint

CI servers have a lot of credentials and secrets flowing through them as they
need access to your source code, deployment environments, etc. It is therefore
critical that they are secured appropriately.

There's a [Security Best Practices page][security-best-practices] on the
internal Pivotal wiki detailing steps you can take to improve the security of
your Concourse deployment.

[security-best-practices]: https://sites.google.com/a/pivotal.io/cloud-foundry/engineering/security-best-practices

## Connecting Flights

If you would like more information about what to do next or how to extend
Concourse further then the following links should be helpful:

* The [pipeline that builds Concourse itself][concourse-pipeline] is open
  source and a good example of a large pipeline that has all different kinds of
  inputs and outputs (check out the *publish* group for our release fan-out).

* The [Concourse documentation][concourse-docs] should be a one stop reference
  for working with Concourse. If something isn't clear or information is
  missing then please let us know so that we can fix it for everyone.
  
[concourse-pipeline]: https://github.com/concourse/concourse/blob/master/ci/pipelines/concourse.yml
[concourse-docs]: http://concourse.ci/
