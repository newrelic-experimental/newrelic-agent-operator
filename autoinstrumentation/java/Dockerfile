# To build one auto-instrumentation image for Java, please:
#  - Download the newrelic `newrelic-agent-$version.jar` to `/newrelic-agent.jar`. This is required as when instrumenting the pod,
#    one init container will be created to copy the jar to your app's container.
#  - Grant the necessary access to the jar. `chmod -R go+r //newrelic-agent.jar`
FROM busybox

ARG version

ADD https://download.newrelic.com/newrelic/java-agent/newrelic-agent/$version/newrelic-agent-$version.jar /newrelic-agent.jar

RUN chmod -R go+r /newrelic-agent.jar
