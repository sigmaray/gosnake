FROM golang:1.23

RUN apt-get update && apt-get install -y ncurses-dev

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/reference/dockerfile/#copy
COPY *.go ./

# Build
RUN GOOS=linux go build cmd/main.go -o ./gosnake

# Run
CMD ["./gosnake"]
