FROM cfinfrastructure/golang
MAINTAINER https://github.com/cloudfoundry/mega-ci
USER root

RUN \
      apt-get update && \
      apt-get -y install --fix-missing \
        runit \
      && \
      apt-get clean

# Create bbluser
RUN mkdir -p /home/bbluser && \
	groupadd -r bbluser -g 433 && \
	useradd -u 431 -r -g bbluser -d /home/bbluser -s /usr/sbin/nologin -c "Docker image user" bbluser && \
	chown -R bbluser:bbluser /go
