FROM alpine:latest
MAINTAINER Alexey Pegov <iam@alexeypegov.com>

RUN mkdir /blog /blog/import
ADD b4v_docker /blog/
ADD blog.toml /blog/
ADD blog.tpl /blog/
ADD public/ /blog/public/
ADD import/blog-23-12-2016.json /blog/import/

RUN apk update && apk add tzdata

EXPOSE 8080
WORKDIR /blog

CMD ./b4v_docker --logtostderr --import ./import/blog-23-12-2016.json