FROM debian:12.9-slim

RUN apt-get update -y && apt-get upgrade -y && \
  apt-get install -y tpm2-tools && \
  rm -rf /var/lib/apt/lists/*


