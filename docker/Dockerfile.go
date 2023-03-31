FROM golang:1.20.2 as builder
WORKDIR /iot-app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY /subscriber . 


# CGO has to be disabled for scratch/alpine
ENV CGO_ENABLED=0 
RUN go build -o backend-server main.go && chmod +x ./backend-server

FROM scratch
LABEL org.opencontainers.image.source https://github.com/benedicthomuth/iot4dobot
COPY /early-dev-frontend /frontend
COPY --from=builder /iot-app/backend-server backend-server
COPY /subscriber/domain.crt .
COPY /subscriber/domain.key .
ENTRYPOINT [ "./backend-server" ]
