FROM golang:1.24.2

# Set destination for COPY
WORKDIR /app

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/reference/dockerfile/#copy
COPY . .

# Download Go modules
RUN go mod download

# Build
RUN GIT_DESCRIBE=$(git describe --tags --always) && \
    TAG=$(echo "$GIT_DESCRIBE" | cut -d'-' -f1) && \
    COMMIT=$(echo "$GIT_DESCRIBE" | cut -d'-' -f3) && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags "-X main.Tag=${TAG} -X main.Commit=${COMMIT}" -o /server cmd/main.go

# Optional:
# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/reference/dockerfile/#expose
EXPOSE 8081

# Run
CMD ["/server"]
