# lgtm-rsvp

![sample](sample.gif)

# how to

## deploy - the easy and cheap way (recommended)

- coming soon

[Video Tutorial](#)


## deploy - the harder and expensive way


1. You do need your own domain name to point to the server

2. It is recommended to deploy it on a cloud provider where DDoS protection feature is built-in, e.g. AWS, GCP... 

```shell
        TCP on Port 80, 443, 22 should be allowed
```

3. Clone this repository

4. Carefully read and modify config.yaml.tmpl and modify name as config.yaml

5. Place "album" under public/images/album

```shell
        The folder name should be exactly "album"

        The folder should contain at least three images

        - title image
        - groom
        - bride

        and needs to be sorted in that exact order.

        All other images will be displayed under "Gallery" sector.

```

6. Set up dependencies

```shell

# this will install dependencies

./hack/setup.sh

```

7. Set up reverse proxy and get certificate for your domain

```shell
# /etc/nginx/nginx.conf
# modify the line below to have only TLS1.2 <=
        ssl_protocols  TLSv1.2 TLSv1.3; # Dropping SSLv3, ref: POODLE

# use file template at net/rsvp-chal.conf
# place it at /etc/nginx/conf.d/
# with domain name changed to yours

sudo systemctl restart nginx

# check if you could reach nginx by your domain name at port 80
# using `curl http://yourdomain.com`

# if all good
# then get the certificates on that domain

sudo certbot --nginx --rsa-key-size 4096 --no-redirect 

# after success
# change the file at /etc/nginx/conf.d/ to net/rsvp.conf with the domain name changed to yours

sudo systemctl restart nginx

```

8. Configure podman

```shell
# /etc/containers/registries.conf
# add this to the last line
unqualified-search-registries = ["docker.io"]

# log in
# with docker ip pw
podman login


```

9. Run!

```shell

# use `help` arg to find out what's available

./run_container.sh

```



## Admin Features

### sign in as admin

- $url/signin

### write your own story

- $url/story/w

### delete story

- $url/r/$storyhash/delete


### sign out

- $url/signout

### allow or block a list of comments

- $url/comment/sudo

```shell
        on this page, you can use command `allow` or `block` to specify
        what action you want to take on a comment-list.json file, whose
        structure is seen below.

        this is also generated under `data/` directory
        when there is any failed attempt to send mail
        when someone posts a comment on the website.
```

```json
[
        {
                "commentid":"",
                "title":"",
                "content":""
        },
        {
                "commentid":"",
                "title":"",
                "content":""
        }
]

```


## develop

- go 1.23+
- sqlite3
- build-essential
- make
- podman

## thanks to

[AndersonChoi](https://github.com/AndersonChoi/wedding-card)
