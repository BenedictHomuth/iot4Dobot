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

resource "hcloud_server" "serveriot" {
  name = "Backend"
  image = "debian-11"
  datacenter = "fsn1-dc14"
  server_type="cx11"
  ssh_keys = [hcloud_ssh_key.iot-dobot.name]
  user_data = var.podman_setup
  public_net {
    ipv4_enabled = true
    ipv6_enabled = true
  }
}

output "ipv4_address" {
  value = hcloud_server.serveriot.ipv4_address
}