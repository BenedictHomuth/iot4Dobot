variable "docker_setup" {
  type = string
  default = <<-EOF
    #!/bin/bash
    #  -ex -> exits, when error accours and shows log
    set -ex
    sudo apt-get update
    sudo apt-get install \
        ca-certificates \
        curl \
        gnupg
    sudo apt-get -y install podman
    podman run -p 4222:4222 docker.io/library/nats:2.9.15
  EOF
  description = "Sets up docker, starts it and adds user to docker group"
}