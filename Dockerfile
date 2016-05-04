# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

# To configure the app, set environment variables and use the command line flags

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/cagedtornado/centralconfig

# Build and install the app inside the container.
RUN go get github.com/cagedtornado/centralconfig/...

# Run the app by default when the container starts.
ENTRYPOINT /go/bin/centralconfig

# Document that the app listens on port 3000.
EXPOSE 3000