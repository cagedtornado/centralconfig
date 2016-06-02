# Start Alpine linux latest
FROM alpine:3.2

# To configure the app, set environment variables and use the command line flags

# Copy the local package files to the container's workspace.
ADD centralconfig /centralconfig
RUN chmod +x /centralconfig ; sync; sleep 1

WORKDIR /

# Run the app by default when the container starts.
ENTRYPOINT ["/centralconfig"]

# Start with the 'serve' command
CMD ["serve"]

# Document that the app listens on port 3000.
EXPOSE 3000
