variable "triton_account" {}
variable "triton_url" {
  default = "https://us-sw-1.api.joyentcloud.com"
}
variable "triton_key_path" {
  default = "path to your ssh key for Triton"
}
variable "triton_key_id" {
  default = "Run 'ssh-keygen -l -E md5 -f <TRITON_KEY_PATH> | cut -d' ' -f2 | cut -d: -f2-'"
}
variable "docker_cert_path" {
  default = "/Users/<USER>/.sdc/docker/<TRITON_ACCOUNT>"
}
variable "docker_host" {
  default = "tcp://us-sw-1.docker.joyent.com:2376"
}
