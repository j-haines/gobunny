FROM ubuntu

RUN apt-get update
RUN apt-get install --yes ca-certificates curl redis-server
