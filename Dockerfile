FROM alpine:3.8
RUN apk update && apk add ca-certificates gettext
COPY configmap /usr/local/bin/
ENTRYPOINT configmap
