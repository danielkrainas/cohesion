FROM alpine:3.4
MAINTAINER Danny Krainas <me@danielkrainas.com>

ENV COHESION_CONFIG_PATH /etc/cohesion.default.yml

COPY ./dist /bin/cohesion
COPY ./config.default.yml /etc/cohesion.default.yml

ENTRYPOINT ["/bin/cohesion"]
