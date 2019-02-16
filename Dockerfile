FROM golang:1.11.5-alpine3.9 AS build
LABEL maintainer=<tim.curless@thinkahead.com>

COPY ./ /go/src/github.com/timcurless/eta
WORKDIR /go/src/github.com/timcurless/eta

RUN apk update && apk add --no-cache git
RUN go get github.com/golang/dep/cmd/dep
RUN dep ensure -v
RUN CGO_ENABLED=0 go build -o /go/bin/eta

FROM alpine:3.9

COPY --from=build /go/bin/eta /usr/local/bin/eta
RUN chmod +x /usr/local/bin/eta
EXPOSE 3000
ENTRYPOINT ["eta"]
