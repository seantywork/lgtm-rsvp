#!/bin/bash

sudo apt update

sudo apt install -y make curl nginx ca-certificates

sudo snap install --classic certbot 

sudo ln -s /snap/bin/certbot /usr/bin/certbot 

sudo install -m 0755 -d /etc/apt/keyrings
sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
sudo chmod a+r /etc/apt/keyrings/docker.asc

# Add the repository to Apt sources:
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu \
  $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
  sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
sudo apt-get update

sudo apt-get install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin


curl -OL https://golang.org/dl/go1.23.2.linux-amd64.tar.gz

sudo tar -C /usr/local -xvf go1.23.2.linux-amd64.tar.gz

rm -rf *.tar.gz

echo "export PATH=\$PATH:/usr/local/go/bin" >> $HOME/.profile


echo "do: source ~/.profile"
