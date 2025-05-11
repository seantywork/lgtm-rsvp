#!/bin/bash

sudo apt update

sudo apt install -y make curl nginx ca-certificates apache2-utils sqlite3 unzip podman

sudo snap install --classic certbot 

sudo ln -s /snap/bin/certbot /usr/bin/certbot 


curl -OL https://golang.org/dl/go1.23.2.linux-amd64.tar.gz

sudo tar -C /usr/local -xvf go1.23.2.linux-amd64.tar.gz

rm -rf *.tar.gz

echo "export PATH=\$PATH:/usr/local/go/bin" >> $HOME/.profile


echo "do: source ~/.profile"
