# lgtm-rsvp

# how to

## deploy


1. You do need your own domain name to point to the server

2. It is recommended to deploy it on a cloud provider where DDoS protection feature is built-in, e.g. AWS, GCP...

3. Clone this repository

4. Carefully read and modify config.yaml

5. You need to create api.json file

If you set "google_comment" other than blank string, it means you're going to use Google mail service for user comments, hence Google mail api should be configured accordingly. The value is Google mail app password.

If you set "kakao_share" other than blank string, it means you're going to use Kakao share button on index page, hence Kakao developer setting should be configured accordingly. The value is Kakao developer REST API key.


```json
{
    "google_comment":"",
    "kakao_share": ""
}

```

6. If you choose to use Oauth2, you also need to create oauth.json file,

Google how to get it.

```json
{
  "web": {
    "client_id": "",
    "project_id": "",
    "auth_uri": "",
    "token_uri": "",
    "auth_provider_x509_cert_url": "",
    "client_secret": "",
    "redirect_uris": [
      "",
      ""
    ]
  }
}

```

7. Place your album folder under ./public directory

It needs to have at least three pictures for title, groom, bride in that order.

All other pictures should be sorted to come after them.

8. Set up dependencies

```shell

# this will install dependencies

./hack/setup.sh

```

9. Set up reverse proxy and get certificate for your domain


10. Run!

```shell

./docker.sh

```



## develop


