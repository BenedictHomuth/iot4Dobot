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
    sudo mkdir -m 0755 -p /etc/apt/keyrings
    curl -fsSL https://download.docker.com/linux/debian/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
    # echo \
    # "deb [arch="$(dpkg --print-architecture)" signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/debian \
    # "$(. /etc/os-release && echo "$VERSION_CODENAME")" stable" | \
    # sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

    # sudo apt-get update
    # sudo apt-get install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
    # sudo groupadd docker
    sudo apt-get -y install podman
    podman run -p 4222:4222 docker.io/library/nats:2.9.15
    sudo usermod -aG docker $USER
    apt install apparmor
    systemctl restart docker
  EOF
  description = "Sets up docker, starts it and adds user to docker group"
}