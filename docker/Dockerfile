FROM ubuntu:14.04

ADD https://github.com/ashwanthkumar/marathon-alerts/releases/download/v0.3.1/marathon-alerts-linux-amd64 /marathon-alerts 
RUN chmod 755  /marathon-alerts

RUN apt-get update && \
    apt-get install -y ca-certificates

WORKDIR /
ENTRYPOINT ["./marathon-alerts"]
