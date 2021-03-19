FROM alpine:latest
RUN apk add --no-cache ca-certificates
RUN apk --no-cache add tzdata
ENV TZ Asia/Jakarta
EXPOSE 8080
COPY ./bin/ /
COPY ./files/etc/testing /
ENTRYPOINT ["/testing"]
