# triton-terraform

[![wercker status](https://app.wercker.com/status/ceee1ebf9da101850ac92639e6e0711d/m "wercker status")](https://app.wercker.com/project/bykey/ceee1ebf9da101850ac92639e6e0711d)

<!-- markdown-toc start - Don't edit this section. Run M-x markdown-toc-generate-toc again -->
**Table of Contents**

- [triton-terraform](#triton-terraform)
    - [Provider](#provider)
    - [Resources](#resources)
        - [`triton_key`](#tritonkey)
        - [`triton_machine`](#tritonmachine)

<!-- markdown-toc end -->

## Provider

You can set up the Triton provider for development by adding the following to
your terraform RC after `go get`ing this repo:

```hcl
providers {
  triton = "triton-terraform"
}
```

Then you'll need to set up the provider in a Terraform config file, like so:

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

Notes:

- you can use the package UUID instead of the name, but Terraform will think
  that you want to change it and recreate the resource every time if you do.
- changing metadata is a little odd. In order to prevent Terraform from thinking
  the metadata needs to be changed every time you run `terraform plan`, known
  keys have been promoted to the top level. So, if you need to access the
  `user-script` key, use `user_script` (just replace the dashes with
  underscores, in general).
- due to a bug in the SDC API for Go (which this tool is implemented against),
  you will have to run `terraform apply` twice to get your tags and metadata to
  apply.
