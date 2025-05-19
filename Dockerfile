FROM golang:1.24-alpine 
# Enable CGO and set the target OS/architecture
# ENV CGO_ENABLED=1 GOOS=linux GOARCH=amd64
# Install required dependencies for CGO and SQLite3
RUN apk add --no-cache gcc musl-dev sqlite-dev make

# Enable CGO and set the target OS/architecture
ENV CGO_ENABLED=1

WORKDIR /app

# Copy Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Install required Go tools
RUN go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest
RUN go install github.com/mattn/go-sqlite3

# Set configuration path
ENV CONFIG_PATH=/app/.config

# Copy application files
COPY . /app
COPY ./config ${CONFIG_PATH}

# Initialize the application
RUN make init

# Build the application with SQLite3 support
RUN go build -tags "sqlite_omit_load_extension" -o app ./cmd/main.go

# Expose the application port
EXPOSE ${PORT}

# Run the application
CMD ["./app"]
