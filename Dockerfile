FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY static/ ./
COPY spb ./

EXPOSE 1234
USER nonroot:nonroot

ENTRYPOINT ["/spb"]
