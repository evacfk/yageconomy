Economy, gambling, waifu and so on features for YAGDPDB

To use:

1. Add 	"github.com/jonas747/yageconomy" to the imports in yags main.go
2. Add "yageconomy.RegisterPlugin()" at the bottom of all the other "RegisterFunctions()" calls in yags main.go

If using docker, you need some more large modifications:

1. clone a version of yageconomy one directory above yagpdb
example layout:

somefolder:
	yageconomy
	yagpdb

2. Use the following Dockerfile

```

FROM golang:stretch as builder

WORKDIR $GOPATH/src

RUN git clone -b yagpdb https://github.com/jonas747/discordgo github.com/jonas747/discordgo \
  && git clone -b dgofork https://github.com/jonas747/dutil github.com/jonas747/dutil \
  && git clone -b dgofork https://github.com/jonas747/dshardmanager github.com/jonas747/dshardmanager \
  && git clone -b dgofork https://github.com/jonas747/dcmd github.com/jonas747/dcmd

COPY yageconomy github.com/jonas747/yageconomy

RUN go get -d -v \
  github.com/jonas747/yagpdb/cmd/yagpdb

# Uncomment during development
COPY yagpdb github.com/jonas747/yagpdb

RUN go get -d github.com/jonas747/yageconomy

# Disable CGO_ENABLED to force a totally static compile.
RUN CGO_ENABLED=0 GOOS=linux go install -v \
  github.com/jonas747/yagpdb/cmd/yagpdb

FROM alpine:latest

WORKDIR /app
VOLUME /app/soundboard \
  /app/cert
EXPOSE 80 443

# We need the X.509 certificates for client TLS to work.
RUN apk --no-cache add ca-certificates

# Add ffmpeg for soundboard support
RUN apk --no-cache add ffmpeg

# Handle templates for plugins automatically
COPY --from=builder /go/src/github.com/jonas747/yagpdb/*/assets/*.html templates/plugins/
COPY --from=builder /go/src/github.com/jonas747/yageconomy/assets/*.html templates/plugins/

COPY --from=builder /go/src/github.com/jonas747/yagpdb/cmd/yagpdb/templates templates/
COPY --from=builder /go/src/github.com/jonas747/yagpdb/cmd/yagpdb/posts posts/
COPY --from=builder /go/src/github.com/jonas747/yagpdb/cmd/yagpdb/static static/

COPY --from=builder /go/bin/yagpdb .

# add extra flags here when running YAGPDB
# Set "-exthttps=true" if using a TLS-enabled proxy such as
# jrcs/letsencrypt-nginx-proxy-companion
# Set "-https=false" do disable https
ENV extra_flags ""

# `exec` allows us to receive shutdown signals.
CMD exec /app/yagpdb -all -pa $extra_flags

```

3. in docker-compose.yml update the services/app/build entries to the following:

```
context: ../../
dockerfile: yagpdb/yagpdb_docker/Dockerfile
```

on one of the lines under "# Handle templates for plugins automatically"

