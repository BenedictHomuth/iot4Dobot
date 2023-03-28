terraform {
  required_providers {
    hcloud = {
      source = "hetznercloud/hcloud"
    }
  }
  required_version = ">= 0.13"
}

provider "hcloud" {
  token = var.token
}

resource "hcloud_ssh_key" "iot-dobot" {
  name = "iot-dobot"
  public_key = file("/Users/bhomuth/.ssh/id_ecdsa.pub")
}

resource "hcloud_server" "test-terraform-server" {
  name="Test-Terraform-Server"
  image="debian-11"
  server_type="cx11"
  ssh_keys = [hcloud_ssh_key.iot-dobot.name]
  user_data = var.docker_setup
  public_net {
    ipv4_enabled = true
    ipv6_enabled = true
  }
}

output "ipv4_address" {
  value = hcloud_server.test-terraform-server.ipv4_address
}