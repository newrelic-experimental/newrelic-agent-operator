# To build one auto-instrumentation image for Node.js, please:
# - Ensure the packages are installed in the `/instrumentation` directory. This is required as when instrumenting the pod,
#   one init container will be created to copy all the content in `/instrumentation` directory to your app's container. Then
#   update the `NODE_OPTIONS` environment variable accordingly.
#   In the first stage, install all the required packages in one custom directory.
#   Then in the second stage, copy the directory to `/instrumentation`.
# - Ensure you have `newrelic`, which also contains instrumentation for `aws-sdk`, `koa`, `superagent`, and `native-metrics`.
# - Grant the necessary access to `/instrumentation` directory. `chmod -R go+r /instrumentation`
ARG NODE_RUNTIME_VERSION=20.12
ARG NEWRELIC_NODE_AGENT_VERSION=latest

FROM node:${NODE_RUNTIME_VERSION} AS build
WORKDIR /operator-build
COPY . .
RUN npm install newrelic@${NEWRELIC_NODE_AGENT_VERSION}

FROM busybox
COPY --from=build /operator-build/build/workspace /instrumentation
RUN chmod -R go+r /instrumentation
