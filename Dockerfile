FROM alpine:latest

RUN apk --no-cache add curl
COPY demo-app /usr/bin/demo-app

EXPOSE 8080

CMD ["/usr/bin/demo-app"]
