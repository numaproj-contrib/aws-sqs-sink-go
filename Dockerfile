####################################################################################################
# base
####################################################################################################
FROM alpine:3.12.3 as base
RUN apk update && apk upgrade && \
    apk add ca-certificates && \
    apk --no-cache add tzdata

COPY dist/aws-sqs-sink /bin/aws-sqs-sink
RUN chmod +x /bin/aws-sqs-sink

####################################################################################################
# sqs-sink
####################################################################################################
FROM scratch as sqs-sink
ARG ARCH
COPY --from=base /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=base /bin/aws-sqs-sink /bin/aws-sqs-sink
ENTRYPOINT [ "/bin/aws-sqs-sink" ]
