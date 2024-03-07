FROM golang:1.21 AS build-stage

WORKDIR /app

COPY * ./
RUN go mod download
RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -o /spb

# Run the tests in the container
# FROM build-stage AS run-test-stage
# RUN go test -v ./...

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /spb /spb

EXPOSE 1234
USER nonroot:nonroot

ENTRYPOINT ["/spb"]
