variable "triton_account" {}
variable "triton_url" {
  default = "https://us-sw-1.api.joyentcloud.com"
}
variable "triton_key_path" {
  default = "path to your ssh key for triton"
}
variable "triton_key_id" {}
variable "docker_cert_path" {
  default = "/Users/<USER>/.sdc/docker/<TRITON_ACCOUNT>"
}
variable "docker_host" {
  default = "tcp://us-sw-1.docker.joyent.com:2376"
}
