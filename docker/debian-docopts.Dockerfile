FROM debian:buster
ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update -y \
    && apt-get install -y --no-install-recommends \
    unzip \
    wget \
    python3-setuptools \
    python3-pip \
    vim \
		git \
    less

RUN apt-get install -y --no-install-recommends \
    sudo \
    tar \
    curl \
    make

WORKDIR /app
COPY *.sh \
    /app/

# install precompiled binary published ==> docopts0
ARG VERSION
RUN wget https://github.com/docopt/docopts/releases/download/$VERSION/docopts_linux_amd64
RUN install -o root -g root -m a+x docopts_linux_amd64 /usr/local/bin/docopts0

# install a golang build env
# predownload the tgz so it get docker cached
RUN wget --quiet https://dl.google.com/go/go1.17.1.linux-amd64.tar.gz
RUN ./update_go.sh
ENV PATH=$PATH:/app:/usr/local/go/bin:/root/go/bin

RUN go get github.com/docopt/docopt-go && go get github.com/docopt/docopts

# intall python version 0.6.1 ==> /usr/local/bin/docopts
RUN pip3 install docopts

# return to basic /app dir
WORKDIR /app
