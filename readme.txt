
YOUTUBE VIDEO : https://www.youtube.com/watch?v=eD88NZN6N6Q


Download nginx binary from https://nginx.org/en/download.html


https://medium.com/@yildirimabdrhm/nginx-is-an-open-source-web-server-software-designed-to-use-as-a-web-server-reverse-proxy-http-7e0cd0bab12


Check compile options
https://nginx.org/en/docs/configure.html


RUN ./configure --prefix=/usr/local/nginx --with-http_auth_request_module --with-http_ssl_module --with-http_stub_status_module
RUN make install

Configuration summary
  + using system PCRE library
  + using system OpenSSL library
  + using system zlib library

  nginx path prefix: "/usr/local/nginx"
  nginx binary file: "/usr/local/nginx/sbin/nginx"
  nginx modules path: "/usr/local/nginx/modules"
  nginx configuration prefix: "/usr/local/nginx/conf"
  nginx configuration file: "/usr/local/nginx/conf/nginx.conf"
  nginx pid file: "/usr/local/nginx/logs/nginx.pid"
  nginx error log file: "/usr/local/nginx/logs/error.log"
  nginx http access log file: "/usr/local/nginx/logs/access.log"
  nginx http client request body temporary files: "client_body_temp"
  nginx http proxy temporary files: "proxy_temp"
  nginx http fastcgi temporary files: "fastcgi_temp"
  nginx http uwsgi temporary files: "uwsgi_temp"
  nginx http scgi temporary files: "scgi_temp"

root@c651c3c87bcd:/nginx-1.20.1# C02W3473Hbuild . -f notify_poc_nginx_go_auth.Dockerfile -ter container run -dit -p 5003:80  notify_poc_nginx_go_auth
[+] Building 46.7s (20/20) FINISHED

Add these to /etc/hosts file on your dev machine
127.0.0.1        box-notify.casb.protect.broadcom.com
127.0.0.1        office-notify.casb.protect.broadcom.com
127.0.0.1        gapps-notify.casb.protect.broadcom.com

docker build . -f notify_poc_nginx_go_base.Dockerfile -t notify_poc_nginx_go_base

docker build . -f notify_poc_nginx_4_way_split.Dockerfile -t notify_poc_nginx_4_way_split

docker container run -dit -p 5000:80 -p 5091:9991 -p 5092:9992 -p 5081:9981 -p 5082:9982 notify_poc_nginx_4_way_split

docker container run -dit -p 80:80 -p 9991:9991 -p 9992:9992 -p 9981:9981 -p 9982:9982 notify_poc_nginx_4_way_split

docker container list

docker exec -it 99931aedd1f3 bash


Demo
- Direct hit the Box Auth
  - Without Authorization Error - App Error
  - With Wrong Origination URL - App Misconfigured - Error
  - ALL Right (Tenant Header)

- Nginx Box Endpoint
 - Without Authorization Error - Nginx
 - Box Valid Request 
 - Gapps Valid Request



