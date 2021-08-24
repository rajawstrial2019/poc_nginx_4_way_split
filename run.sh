#!/bin/bash

cd /usr/local/nginx/sbin/

./nginx

cd /casb_poc/go-app/box_auth_server
./boxauthserver &

cd /casb_poc/go-app/office_auth_server
./officeauthserver &

cd /casb_poc/go-app/box_app_server
./boxappserver &

cd /casb_poc/go-app/office_app_server
./officeappserver &

/bin/bash