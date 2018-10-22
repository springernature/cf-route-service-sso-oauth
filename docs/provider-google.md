# Google OIDC Provider

## Manual CURL commands

Step 1: Get Authorization Code

```bash
https://accounts.google.com/o/oauth2/v2/auth?\
client_id=${CLIENT_ID}.apps.googleusercontent.com&\
response_type=code&\
scope=openid%20email&\
redirect_uri=https://sso-test-gerard.snpaas.eu/cfsso/callback
```

Step 1: Get Authorization Code (RESPONSE)

```bash
response:
https://sso-test-gerard.snpaas.eu/cfsso/callback?code=4%2Fdg_funny_code_in_long_string_kVl8&scope=https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fuserinfo.email+https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fplus.me&authuser=0&hd=springernature.com&session_state=54d_another_funny_code_36s&prompt=none#
```

Step 2: Redeem! Get an access token and oidc tokenID

```bash
curl -s "https://www.googleapis.com/oauth2/v3/token" \
-X POST \
-H "Content-Type: application/x-www-form-urlencoded" \
-H "Accept: application/json" -d \
"code=4%2Fdg_funny_code_in_long_string_kVl8&\
client_id=${CLIENT_ID}.apps.googleusercontent.com&\
client_secret=$(cat credentials.json | jq -r .web.client_secret)&\
redirect_uri=https://sso-test-gerard.snpaas.eu/cfsso/callback&\
grant_type=authorization_code"
```
