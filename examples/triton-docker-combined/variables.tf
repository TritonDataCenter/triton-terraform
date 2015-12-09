variable "account" {}
variable "triton_url" {}
variable "key_path" {}
variable "key_id" {}
variable "cert_path" {}
variable "docker_host" {
 default = "tcp://us-sw-1.docker.joyent.com:2376"
}
