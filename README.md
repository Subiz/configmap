# Configmap (1.0.12)
Read configs from vault based on configmap.yaml

# Dev
```
go test ./...
go build
./configmap help

```
# Usage

config.yaml
```yaml
secret/stripe/dev:
  api_key: 12345555555555555555555555
  file: 11111111
s3:
  key: default value

```
configmap.yaml
```yaml
stripe_apikey: secret/stripe/dev.api_key
s3_apikey: s3.key
~/workspace/x: secret/stripe/dev.file

```
```sh
$ configmap -config config configmap.yaml > config
```
```
stripe_apikey="12345555555555555555555555"
s3_apikey=""

```
```
cat ~/workspace/x
11111111
```
