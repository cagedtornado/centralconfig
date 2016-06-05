# Start Alpine linux latest
FROM golang

# To configure the app, set environment variables and use the command line flags

# Copy the local package files to the container's workspace and set permissions
COPY centralconfig /
RUN chmod +x /centralconfig ; sync; sleep 1

WORKDIR /

# Start with the 'serve' command
CMD ["/centralconfig", "serve"]

# Document that the app listens on port 3000.
EXPOSE 3000
