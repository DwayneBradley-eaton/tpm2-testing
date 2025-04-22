# Testing TPM 2.0 Chipsets in ECU-150 and UNO-2271g

## Requirements

Have the following installed on your local machine:

* [Docker](https://docs.docker.com/engine/install/)
* [Go](https://go.dev/doc/install)

## Steps

### Create a new Docker container images with the TPM2 tools installed in it

Use the `Dockerfile` included in this repo to build different images for each target:

#### AMD64

`docker build --platform=linux/amd64 -t debian-tpm2-tools:12.9-slim-amd64 .`

#### ARM64

`docker build --platform=linux/arm64 -t debian-tpm2-tools:12.9-slim-arm64 .`


### Save the images to tar files

#### AMD64

`docker save debian-tpm2-tools:12.9-slim-amd64 > debian-tpm2-tools-arm64.tar`

#### ARM64

`docker save debian-tpm2-tools:12.9-slim-arm64 > debian-tpm2-tools-arm64.tar`

### SCP the tar files to the target machine

#### AMD64 (i.e. UNO-2271g)

`scp debian-tpm2-tools-amd64 <<USERNAME>>@<<HOST_IP>>:/tmp`

#### ARM64 (i.e. ECU-150)

`scp debian-tpm2-tools-arm64 <<USERNAME>>@<<HOST_IP>>:/tmp`


### SSH to the devices

`ssh <<USERNAME>>@<<HOST_IP>>`

### Load the container images from the tar files to Docker

#### AMD64

`docker load < /tmp/debian-tpm2-tools-amd64.tar`

#### ARM64

`docker load < /tmp/debian-tpm2-tools-arm64.tar`

### Run the container

#### AMD64 (i.e. UNO-2271g)

`docker run --device=/dev/tpm0 --device=/dev/tpmrm0 -v /tmp:/tmp --rm -ti debian-tpm2-tools:12.9-slim-amd64 /bin/bash`

#### ARM64 (i.e. ECU-150)

`docker run --device=/dev/tpm0 --device=/dev/tpmrm0 -v /tmp:/tmp --rm -ti debian-tpm2-tools:12.9-slim-arm64 /bin/bash`

### Extract the Certs

You can run the same command in each container to extract the certs:

`tpm2_getekcertificate -o /tmp/ek_cert.der -o /tmp/ecc_cert.der`

### Exit the container

Just run `exit` to exit and stop the container

### SCP the certs back to your machine

Do this for both device but be sure to save them to different locations.

`scp <<USERNAME>>@<<HOST_IP>>:/tmp/ek_cert.der <<YOUR_LOCAL_PATH_HERE>>`

### Parse the certs with `openssl`

Parsing the certificate with the generic `openssl`` command does not cause an error for any of the certificates as
`openssl` is just politely ignoring the extra “junk” data in the cert files:

#### AMD64 (i.e. UNO-2271g)

`openssl x509 -in ./ecu-150/ek_cert.der -text`

#### ARM64 (i.e. ECU-150)

`openssl x509 -in ./uno-2271g/ek_cert.der -text`

### Parse the certificates with Go

Run the Go program in this repo to parse the certificates:

#### ARM64 (i.e. ECU-150)

`go run main.go -f=./ecu-150/ek_cert.der`

This results in the “trailing data” error:

```
2025/04/22 09:24:34 x509: trailing data
exit status 1
```

#### AMD64 (i.e UNO-2271g)

`go run main.go -f=./uno-2271g/ek_cert.der`

This one parses just fine and outputs a few fields from the certificate:

```
Name
Not before 2022-02-14 21:05:47 +0000 UTC
Not after 2042-02-10 21:05:47 +0000 UTC
```

__NOTE__: The “common name” field was blank in the cert so you do not see anything in that line

