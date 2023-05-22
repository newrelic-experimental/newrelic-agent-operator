# To build one auto-instrumentation image for Python, please:
# - Ensure the packages are installed in the `/instrumentation` directory. This is required as when instrumenting the pod,
#   one init container will be created to copy all the content in `/instrumentation` directory to your app's container. Then
#   update the `PYTHONPATH` environment variable accordingly.
#   In the first stage, install all the required packages in one custom directory with `pip install --target`.
#   Then in the second stage, copy the directory to `/instrumentation`.
# - Ensure you have `newrelic-distro` and `newrelic` installed.
#   Those two packages are essential to Python auto-instrumentation.
# - Grant the necessary access to `/instrumentation` directory. `chmod -R go+r /instrumentation`

FROM python:3.10-alpine AS build

WORKDIR /operator-build

# WARNING: Disabling optional C extension components of the Python agent 
# will result in some non core features of the Python agent, such as
# capacity analysis instance busy metrics, that will not be available.
# Pure Python versions of code supporting some features, rather than the
# optimised C versions, will also be used resulting in additional overheads.

ENV NEW_RELIC_EXTENSIONS = False

ADD requirements.txt .

RUN mkdir workspace && pip install --target workspace -r requirements.txt

FROM busybox

COPY --from=build /operator-build/workspace /instrumentation

# Bypassing the python prefix match to the python version since we are not mounting to the interpreters per user site-packages location.
RUN sed -i 's/if python_prefix_matches and python_version_matches:/if True:/g' /instrumentation/newrelic/bootstrap/sitecustomize.py

RUN chmod -R go+r /instrumentation
