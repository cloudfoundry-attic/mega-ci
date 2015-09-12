FROM ubuntu:14.04
#dpkg-query --show
RUN apt-get update && apt-get install -y \
  autoconf=2.69-6 \
  awscli=1.2.9-2 \
  bison=2:3.0.2.dfsg-2 \
  build-essential=11.6ubuntu6 \
  curl=7.35.0-1ubuntu2.5 \
  gcc=4:4.8.2-1ubuntu6 \
  git=1:1.9.1-1ubuntu0.1 \
  jq=1.3-1.1ubuntu1 \
  libffi-dev:amd64=3.1~rc1+r3.0.13-12 \
  libgdbm-dev=1.8.3-12build1 \
  libgdbm3:amd64=1.8.3-12build1 \
  libncurses5-dev:amd64=5.9+20140118-1ubuntu1 \
  libreadline6-dev:amd64=6.3-4ubuntu2 \
  libsqlite3-dev:amd64=3.8.2-1ubuntu2.1 \
  libssl-dev:amd64=1.0.1f-1ubuntu2.15 \
  libxml2-dev:amd64=2.9.1+dfsg1-3ubuntu4.4 \
  libxslt1-dev=1.1.28-2build1 \
  libyaml-dev:amd64=0.1.4-3ubuntu3.1 \
  make=3.81-8.2ubuntu3 \
  openssl=1.0.1f-1ubuntu2.15 \
  sqlite3=3.8.2-1ubuntu2.1 \
  unzip=6.0-9ubuntu1.3 \
  zlib1g-dev=1:1.2.8.dfsg-1ubuntu1 \
  zlibc=0.9k-4.1

RUN cd /root && \
  curl -L -J -O https://cache.ruby-lang.org/pub/ruby/2.2/ruby-2.2.3.tar.gz && \
  tar xf ruby-2.2.3.tar.gz && \
  cd /root/ruby-2.2.3 && \
  ./configure --prefix=/usr/local && \
  make && \
  make install

RUN gem install bosh_cli -v 1.3071.0

RUN cd /root && \
 curl -L -J -O https://storage.googleapis.com/golang/go1.5.linux-amd64.tar.gz && \
 tar xf go1.5.linux-amd64.tar.gz && \
 mv go /usr/local/go
ENV PATH /usr/local/go/bin:$PATH

RUN cd /root && \
  git clone https://github.com/square/certstrap && \
  cd certstrap && \
  git checkout 9741482cae85edb8f2c3f147428ec2621d849809 && \
  ./build && \
  mv bin/certstrap /usr/local/bin/certstrap

RUN cd /root && \
  curl -L -J -O https://s3.amazonaws.com/bosh-init-artifacts/bosh-init-0.0.72-linux-amd64 && \
  mv bosh-init-* /usr/local/bin/bosh-init && \
  chmod +x /usr/local/bin/bosh-init

RUN cd /root && \
  curl -L -J -O https://github.com/cloudfoundry-incubator/spiff/releases/download/v1.0.7/spiff_linux_amd64.zip && \
  unzip spiff_linux_amd64.zip && \
  mv spiff /usr/local/bin/spiff

RUN cd /root && \
  rm -Rf *

RUN echo "Checking that all dependencies are installed..." && \
  aws --version && \
  bosh --version && \
  printf "bosh-init " && \
  bosh-init --version && \
  certstrap --version && \
  git --version && \
  jq --version && \
  openssl version && \
  spiff --version
