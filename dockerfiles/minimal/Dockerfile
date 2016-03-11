FROM ubuntu:14.04
MAINTAINER https://github.com/cloudfoundry/mega-ci

RUN \
      apt-get update && \
      apt-get -qqy install --fix-missing \
            build-essential \
            curl \
            git \
            libreadline6 \
            libreadline6-dev\
            wget \
      && \
      apt-get clean
