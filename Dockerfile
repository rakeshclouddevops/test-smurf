# FROM golang:1.23



# WORKDIR /go/src/app
# COPY . .

# # Build the Go application
# RUN go build -o smurf main.go
# RUN chmod +x smurf


# # Set entrypoint to the build binary
# ENTRYPOINT ["./go/src/app/smurf"]



# Dockerfile
FROM golang:1.23-alpine

# Install required packages
RUN apk add --no-cache \
    git \
    bash \
    curl \
    unzip

# Install Terraform
RUN curl -fsSL https://releases.hashicorp.com/terraform/1.5.7/terraform_1.5.7_linux_amd64.zip -o terraform.zip && \
    unzip terraform.zip && \
    mv terraform /usr/local/bin/ && \
    rm terraform.zip

WORKDIR /go/src/app
COPY . .

# Build the Go application
RUN go build -o smurf main.go
RUN chmod +x smurf

RUN chmod +x entrypoint.sh

ENTRYPOINT ["./entrypoint.sh"]

