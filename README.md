# triton-terraform

[![wercker status](https://app.wercker.com/status/ceee1ebf9da101850ac92639e6e0711d/m "wercker status")](https://app.wercker.com/project/bykey/ceee1ebf9da101850ac92639e6e0711d)

<!-- markdown-toc start - Don't edit this section. Run M-x markdown-toc-generate-toc again -->
**Table of Contents**

- [triton-terraform](#triton-terraform)
    - [Installing](#installing)
        - [From a release](#from-a-release)
        - [From source](#from-source)
    - [Provider](#provider)
    - [Resources](#resources)
        - [`triton_key`](#tritonkey)
        - [`triton_machine`](#tritonmachine)
        - [`triton_firewall_rule`](#tritonfirewallrule)

<!-- markdown-toc end -->

## Installing

### From a release

Download a release from the
[releases page on Github](https://github.com/joyent/triton-terraform/releases),
then add the binary to your provider config in `~/.terraformrc` like so:

```hcl
providers {
  triton = "/path/where/you/expanded/triton-terraform"
}
```

### From source

If you want the latest version, install [go](https://golang.org/) and run `go
get github.com/joyent/triton-terraform`. The plugin should be built and added to
your Go binary folder. You'll then need to add the binary to your provider
config in `~/.terraformrc` like so:

```hcl
providers {
  triton = "triton-terraform"
}
```

## Provider

You'll need to set up the provider in a Terraform config file, like so:

```hcl
provider "triton" {
  account = "your-account-name"
  key = "~/.ssh/joyent.id_rsa" # the path to your key. If removed, defaults to ~/.ssh/id_rsa
  key_id = "50:87:72:54:cb:25:bf:af:b2:c9:61:19:59:93:fb:ab" # the corresponding key signature from your account page
}
```

## Resources

### `triton_key`

Creates and manages authentication keys in Triton. Do note that any change to
this resource, once created, will result in the old resource being destroyed and
recreated.

`name` is optional if there is a comment set on the key.

```hcl
resource "triton_key" "testkey" {
  name = "test key"
  key = "${file("some/other/id_rsa.pub")}"
}
```

### `triton_machine`

Creates and manages machines in Triton. Below is a fairly complete resource:

```hcl
resource "triton_machine" "testmachine" {
  name = "test-machine"
  package = "g3-standard-0.25-smartos"
  image = "842e6fa6-6e9b-11e5-8402-1b490459e334"
  tags = {
    test = "hello!"
  }
  networks = ["42325ea0-eb62-44c1-8eb6-0af3e2f83abc"]
}
```

### `triton_firewall_rule`

Creates and manages firewall rules in Triton. Note that the API currently
defaults rules to being disabled, so this provider does too.

```hcl
resource "triton_firewall_rule" "testrule" {
    rule = "FROM any TO tag www ALLOW tcp PORT 80"
    enabled = true
}
```

Notes:

- you can use the package UUID instead of the name, but Terraform will think
  that you want to change it and recreate the resource every time if you do.
- to use metadata keys, change dashes to underscores, and use them as top-level
  keys in the resource. For example, `user-script` becomes `user_script`.

## Using the Terraform Docker Provider

The [Terraform Docker provider](https://terraform.io/docs/providers/docker/index.html) needs to be configured with the address to the Docker API host and with the path to a directory that contains valid TLS certificates for authentication. The Docker [helper script](https://github.com/joyent/sdc-docker/tree/master/docs/api#the-helper-script) can be used to configure the Terraform provider. Terraform will read the values in the `DOCKER_HOST` and `DOCKER_CERT_PATH` environment variables that you generate from the script. Alternatively, you can explicitly configure the values in a Terraform provider block. For example:

```
provider "docker" {
  host = "tcp://us-east-1.docker.joyent.com:2376"
  cert_path = "/Users/localuser/.sdc/docker/jill"
}
```
