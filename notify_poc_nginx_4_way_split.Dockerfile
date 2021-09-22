#NOTE - Apache is using Python 3.7, so use the same version of Image
FROM notify_poc_nginx_go_base


RUN mv /usr/local/nginx/conf/nginx.conf /usr/local/nginx/conf/nginx.conf.original
COPY simple_nginx.conf /usr/local/nginx/conf/nginx.conf

RUN . ~/.profile; mkdir $GOPATH/go-web

############################################
# Build GO - App Web
############################################
WORKDIR /casb_poc/go-app/box_app_server

COPY appserver/boxappserver.go boxappserver.go
RUN . ~/.profile; go build boxappserver.go;

WORKDIR /casb_poc/go-app/gapps_app_server

COPY appserver/gappsappserver.go gappsappserver.go
RUN . ~/.profile; go build gappsappserver.go;

############################################
# Build GO - Auth Web
############################################
WORKDIR /casb_poc/go-app/box_auth_server

COPY authserver/boxauthserver.go boxauthserver.go
RUN . ~/.profile; go build boxauthserver.go;

WORKDIR /casb_poc/go-app/gapps

COPY gapps .
RUN . ~/.profile; go build httpd/gappsauthserver.go;


RUN mkdir -p /var/log/casb/

#Add GoWeb as a service - NOT WORKING
#COPY goweb.service /lib/systemd/system/goweb.service

#Configure NGINX Proxy to Goweb
#COPY appserver.conf /usr/local/nginx/conf.d/0000_appserver.conf

#RUN service nginx start

#CMD ["nginx", "-g", "daemon off;"]

COPY run.sh /casb_poc/run.sh
RUN chmod 777 /casb_poc/run.sh

CMD /casb_poc/run.sh
