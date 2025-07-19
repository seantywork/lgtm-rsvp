#!/bin/bash

sudo apt update

sudo apt install -y curl nginx ca-certificates sqlite3 unzip podman

sudo snap install --classic certbot 

sudo ln -s /snap/bin/certbot /usr/bin/certbot 
