FROM alpine:3.18.8
RUN apk add git>=2.38
COPY devx /usr/bin/devx
RUN mkdir /app
WORKDIR /app
ENTRYPOINT ["devx"]
