# To build one auto-instrumentation image for dotnet, please:
#  - Download the newrelic dotnet artefacts to `/instrumentation` directory. This is required as when instrumenting the pod,
#    one init container will be created to copy the files to your app's container.
#  - Grant the necessary access to the files in the `/instrumentation` directory.
#  - Following environment variables are injected to the application container to enable the auto-instrumentation.
#    CORECLR_ENABLE_PROFILING=1
#    CORECLR_PROFILER={36032161-FFC0-4B61-B559-F6C5D41BAE5A}
#    CORECLR_PROFILER_PATH=%InstallationLocation%/libNewRelicProfiler.so
#    CORECLR_NEWRELIC_HOME=%InstallationLocation% 

FROM alpine:latest as build

RUN apk update && apk add ca-certificates

ARG version

WORKDIR /instrumentation

RUN wget -c "https://download.newrelic.com/dot_net_agent/latest_release/newrelic-dotnet-agent_${version}_arm64.tar.gz" -O - | tar -xz --strip-components 1

FROM busybox

COPY --from=build /instrumentation /instrumentation

RUN chmod -R go+r /instrumentation