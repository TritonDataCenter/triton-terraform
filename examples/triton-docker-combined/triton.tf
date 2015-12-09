# See http://www.joyent.com/blog/introducing-hashicorp-terraform-for-joyent-triton

provider "triton" {
  account = "${var.account}"
  key = "${var.key_path}"
  key_id = "${var.key_id}"
  url = "${var.triton_url}"
}

provider "docker" {
  host = "${var.docker_host}"
  cert_path = "${var.cert_path}"
}

resource "docker_image" "nginx" {
  name = "nginx:latest"
  keep_updated = true
}

resource "docker_container" "nginx" {
  count = 1
  name = "nginx-terraform-${format("%02d", count.index+1)}"
  image = "${docker_image.nginx.latest}"
  must_run = true

  env = ["env=test", "role=test"]

  ports {
    internal = 80
    external = 80
  }
}

resource "triton_machine" "testmachine" {
  name = "test-machine"
  package = "t4-standard-512M"
  image = "ffe82a0a-83d2-11e5-b5ac-f3e14f42f12d"

  count = 1
}

resource "triton_machine" "windowsmachine" {
  name = "win-test-terraform"
  package = "333814c2-b4a7-481f-9f86-73bb4122c7a3"
  image = "66810176-4011-11e4-968f-938d7c9edfa2"

  count = 1
}
