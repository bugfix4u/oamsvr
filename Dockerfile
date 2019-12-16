FROM ubuntu:18.04

RUN apt-get update && \
    apt-get -y install \
   vim net-tools iputils-ping openssl

# Load Go binary executable
COPY oamsvr /
COPY oamsvr_cert.conf /

# Build self signed cert
RUN openssl req -x509 -nodes -days 3650 -newkey rsa:2048 -keyout /oamsvr.key -out /oamsvr.crt -config /oamsvr_cert.conf

# Luanch it at startup
ENTRYPOINT /oamsvr