#NOTE - Apache is using Python 3.7, so use the same version of Image
FROM ubuntu:16.04

#Update apt-get
RUN apt-get update
RUN apt install -y apt-utils

#Install Vi
RUN apt-get install -y vim

#install lsb-release
RUN apt-get install -y lsb-release

#Install Curl
RUN apt-get install -y curl wget


############################################
# Download nginx Code and compile it
############################################
WORKDIR /localbuild
RUN curl -O https://nginx.org/download/nginx-1.20.1.tar.gz
RUN tar xf nginx-1.20.1.tar.gz
WORKDIR /localbuild/nginx-1.20.1

RUN apt-get install -y build-essential

#Install ZLIB Packages
RUN apt-get install zlib1g-dev

#Install PCRE Packages
RUN apt-get install -y libpcre3
RUN apt-get install -y libpcre3-dev

#Install SSL Packages
RUN apt-get install -y libssl-dev

#RUN ./configure --sbin-path=/usr/bin/nginx --conf-path=/etc/nginx/nginx.conf --error-log-path=/var/log/nginx/error.log --http-log-path=/var/log/nginx/access.log --with-pcre --pid-path=/var/run/nginx.pid --with-http_ssl_module
RUN ./configure --prefix=/usr/local/nginx \
            --sbin-path=/usr/sbin/nginx \
            --modules-path=/usr/lib/nginx/modules \
            --conf-path=/etc/nginx/nginx.conf \
            --error-log-path=/var/log/nginx/error.log \
            --http-log-path=/var/log/nginx/access.log \
            --pid-path=/run/nginx.pid \
            --lock-path=/var/lock/nginx.lock \
            --with-ipv6 \
            --with-http_auth_request_module \
            --with-http_realip_module \
            --with-http_ssl_module \
            --with-http_stub_status_module
RUN make
RUN make install

#Remove the DONLOADED Tar
RUN rm /localbuild/nginx-1.20.1.tar.gz

############################################
# Install and Configure Go
############################################
WORKDIR /localbuild/go_install
#Install Go Dependencies
RUN apt install -y build-essential software-properties-common curl gdebi net-tools wget curl sqlite3 dirmngr apt-transport-https leafpad git sudo unzip socat bash-completion checkinstall imagemagick openssl

#Clean any tmp files
RUN apt-get clean \
    && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

#Download Go tar
RUN curl -O https://dl.google.com/go/go1.16.7.linux-amd64.tar.gz

#Install Go
RUN tar -xvf go1.16.7.linux-amd64.tar.gz -C /usr/local

#Remove the DONLOADED Tar
RUN rm go1.16.7.linux-amd64.tar.gz

#Setup Go Env
RUN chown -R root:root /usr/local/go \
    && mkdir -p $HOME/go/bin \
    && mkdir -p $HOME/go/src

RUN echo "export GOPATH=\$HOME/go" >> ~/.profile \
    && echo "export PATH=\$PATH:\$GOPATH/bin" >> ~/.profile \
    && echo "export PATH=\$PATH:\$GOPATH/bin:/usr/local/go/bin" >> ~/.profile

CMD /casb_poc/run.sh
