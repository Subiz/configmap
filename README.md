# Usage

configmap.yaml
```yaml
stripe_apikey:
  secret/stripe/dev(api_key): "222222222223333333333333"

  #- default vaule
~/workspace/x:
  "secret/stripe/file(ke)": asdlkfjkalsjdfkljasdklfj

s3_apikey: default value

```

```sh
$ configmap -addr=https://vault.subiz.com -token=12345 configmap.yaml > config
```
```
stripe_apikey="12345555555555555555555555"
s3_apikey="default value"

```
```
cat ~/workspace/x
11111111
```
