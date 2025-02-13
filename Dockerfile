FROM alpine:3.19.6
RUN apk add git>=2.38
COPY devx /usr/bin/devx
RUN mkdir /app
WORKDIR /app
ENTRYPOINT ["devx"]
