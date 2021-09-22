#!/bin/bash

cd /usr/local/nginx/sbin/

./nginx

cd /casb_poc/go-app/box_auth_server
./boxauthserver &

cd /casb_poc/go-app/gapps
./gappsauthserver &

cd /casb_poc/go-app/box_app_server
./boxappserver &

cd /casb_poc/go-app/gapps_app_server
./gappsappserver &

/bin/bash