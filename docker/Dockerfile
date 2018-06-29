FROM alpine:latest
MAINTAINER Toni Moreno <toni.moreno@gmail.com>

ADD ./resistor-last.tar.gz /

VOLUME ["/opt/resistor/conf", "/opt/resistor/log"]

EXPOSE 8090

WORKDIR /opt/resistor
COPY ./resistor.toml ./conf/
COPY ./start.sh /

ENTRYPOINT ["/start.sh"]
