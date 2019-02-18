FROM node:10-alpine as webbuilder

RUN git clone --depth=1 https://github.com/vicanso/location.git /location \
  && cd \location \
  && npm i \
  && npm run build

FROM golang:1.11-alpine as builder

RUN apk update \
  && apk add git make gcc \
  && git clone --depth=1 https://github.com/vicanso/location.git /location 

COPY --from=webbuilder /location/web/build /location/web/build

RUN go get -u github.com/gobuffalo/packr/v2/packr2 \
  && cd /location \
  && packr2 \
  && make build

FROM alpine

EXPOSE 7001

COPY --from=builder /location/location /usr/local/bin/location

CMD ["location"]

HEALTHCHECK --interval=10s --timeout=3s \
  CMD location --mode=check || exit 1