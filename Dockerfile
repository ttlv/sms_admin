FROM alpine
RUN apk --update upgrade && \
    apk add ca-certificates && \
    apk add tzdata && \
    rm -rf /var/cache/apk/*

ADD sms-admin /bin/
CMD /bin/sms-admin

