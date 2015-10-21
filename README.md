# triton-terraform

[![wercker status](https://app.wercker.com/status/ceee1ebf9da101850ac92639e6e0711d/m "wercker status")](https://app.wercker.com/project/bykey/ceee1ebf9da101850ac92639e6e0711d)

<!-- markdown-toc start - Don't edit this section. Run M-x markdown-toc-generate-toc again -->
**Table of Contents**

- [triton-terraform](#triton-terraform)
    - [Provider](#provider)
    - [Resources](#resources)
        - [`triton_key`](#tritonkey)

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
