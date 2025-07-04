FROM golang:alpine AS build-env
WORKDIR /root/
COPY ./ ./
RUN apk update \
  && apk add git \
  && CGO_ENABLED=0 GOOS=linux go build -a

FROM alpine:latest AS runtime-env
WORKDIR /root/
COPY --from=build-env /root/frigate-share ./
COPY ./main.gohtml ./
COPY ./icons ./icons
COPY ./manifest.json ./
RUN apk update \
  && apk add ca-certificates curl bash
CMD ["./frigate-share"]
