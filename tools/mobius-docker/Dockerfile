FROM alpine:3.20@sha256:b3119ef930faabb6b7b976780c0c7a9c1aa24d0c75e9179ac10e6bc9ac080d0d
LABEL maintainer="Mobius Developers"

RUN apk --update add ca-certificates
RUN apk --no-cache add jq

# Create mobius group and user
RUN addgroup -S mobius && adduser -S mobius -G mobius

USER mobius

COPY mobius /usr/bin/

CMD ["mobius", "serve"]
