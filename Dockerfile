#Base image
FROM golang:1.13

# Set the Current Working Directory inside the container
WORKDIR /go/src/github.com/Dayasagara/meeting-scheduler

# Copy everything from the current directory to the PWD(Present Working Directory) inside the container
COPY . .

# Download all the dependencies
RUN go get -d -v ./...

# Install the package
RUN go install -v ./...

# This container exposes port 8080 to the outside world
EXPOSE 8080

# Run the executable
CMD ["meeting-scheduler"]
