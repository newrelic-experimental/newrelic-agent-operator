# To build one auto-instrumentation image for PHP, please:
#  - Download the newrelic php artefacts to `/instrumentation` directory. This is required as when instrumenting the pod,
#    one init container will be created to copy the files to your app's container.
#  - Grant the necessary access to the files in the `/instrumentation` directory.
FROM alpine:latest as build

RUN apk update && apk add ca-certificates

ARG version

WORKDIR /instrumentation

RUN wget -c "https://download.newrelic.com/php_agent/release/newrelic-php5-${version}-linux-musl.tar.gz" -O - | tar -xz --strip-components 1

FROM busybox

COPY --from=build /instrumentation /instrumentation

RUN chmod -R go+r /instrumentation