variable "podman_setup" {
  type = string
  default = <<-EOF
    #!/bin/bash
    #  -ex -> exits, when error accours and shows log
    set -ex
    sudo apt-get update
    sudo apt-get -y install \
        ca-certificates \
        curl \
        gnupg
    sudo apt-get -y install podman
    podman run -p 4222:4222 --name=nats --network=host -d docker.io/library/nats:2.9.15
    podman run -p 8080:8080 --name=backend --network=host -d -e NATS_CONN_STRING="localhost:4222" ghcr.io/benedicthomuth/iot4dobot/server:latest
  EOF
  description = "Sets up podman and start NATS + Backend"
}