# Digital Dobot Twin Project

<p>
    <a href="https://github.com/benedicthomuth/iot4dobot/actions/workflows/backend.yaml">
      <img src="https://github.com/benedicthomuth/iot4dobot/actions/workflows/backend.yaml/badge.svg?branch=main"/>
    </a>
</p>

This project aims to create a digital twin of the Dobot M1 robot arm. The digital twin can be controlled through a Python script inside M1 Studio, which sends positional data to a NATS message queue. A backend Go program listens to the queue and sends the positional data to a WebSocket connection. This connection is read by BuildWagon, an online IDE specialized in developing for the HoloLens from Microsoft. BuildWagon updates the position of the 3D model, creating the illusion of a digital twin.

Whenever there are modifications made to the backend Go code, a GitHub Action is triggered, which then compiles the code, builds a container, and publishes it to the Packages tab of the repository. 

## Prerequisites

- Dobot M1 robot arm
- M1 Studio
- Hetzner Cloud account
- BuildWagon account

## Setup

- Clone this repository to your local machine.
- Set up a VM on Hetzner Cloud using the Terraform file provided in the repository.
  - NATS will be started automatically via a cloud-init script
  - The same goes for the go backend
- Establish a connection from BuildWagon to the WebSocket connection.
- Run the Python script inside M1 Studio to send positional data to the NATS message queue.
- The 3D model should now update its position based on the positional data sent by the Dobot M1 robot arm.

## SSL Certificate Configuration & General Considerations
Following this https://nodeployfriday.com/posts/self-signed-cert/ tutorial the cert is created.

The following cert is then integrated into the docker container. This is NOT SAFE for production and is only used as a first prototype.

When the Hetzner IP changes it might be neccessary to add the new IP to the "/subscriber/requestCert.txt". This must then be published and run through the GitHub Action Workflow in order to have a new container.