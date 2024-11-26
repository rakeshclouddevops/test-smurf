# FROM golang:1.23-alpine


# WORKDIR /go/src/app
# COPY . .

# # Build the Go application
# RUN go build -o smurf main.go && \
#     mv smurf /usr/local/bin/ && \
#     chmod +x /usr/local/bin/smurf

# ENTRYPOINT ["/usr/local/bin/smurf"]

# FROM golang:1.23-alpine

# # Install required packages
# RUN apk add --no-cache \
#     git \
#     bash \
#     curl \
#     unzip

# # Install Terraform
# RUN curl -fsSL https://releases.hashicorp.com/terraform/1.5.7/terraform_1.5.7_linux_amd64.zip -o terraform.zip && \
#     unzip terraform.zip && \
#     mv terraform /usr/local/bin/ && \
#     rm terraform.zip

# # Install AWS CLI v2
# RUN curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip" && \
# unzip awscliv2.zip && \
# ./aws/install && \
# rm -rf aws awscliv2.zip

# WORKDIR /go/src/app
# COPY . .

# # Build the Go application
# RUN go build -o smurf main.go && \
#     mv smurf /usr/local/bin/ && \
#     chmod +x /usr/local/bin/smurf

# # Ensure the entrypoint.sh script is executable and starts with a shebang
# RUN sed -i '1s|^|#!/bin/sh\n|' entrypoint.sh && \
#     chmod +x entrypoint.sh
# # RUN chmod +x entrypoint.sh

# ENTRYPOINT ["/go/src/app/entrypoint.sh"]


FROM golang:1.23-alpine

# Install required packages for Go, curl, unzip, and Helm
RUN apk add --no-cache \
    git \
    bash \
    curl \
    unzip

# Install Terraform (version 1.5.7) if required
RUN curl -fsSL https://releases.hashicorp.com/terraform/1.5.7/terraform_1.5.7_linux_amd64.zip -o terraform.zip && \
    unzip terraform.zip && \
    mv terraform /usr/local/bin/ && \
    rm terraform.zip

# Install Helm (v3.11.0)
RUN curl https://get.helm.sh/helm-v3.11.0-linux-amd64.tar.gz -o helm.tar.gz && \
    tar -zxvf helm.tar.gz && \
    mv linux-amd64/helm /usr/local/bin/ && \
    rm -rf linux-amd64 helm.tar.gz

# Set the working directory for Go application code
WORKDIR /go/src/app

# Copy Go application and entrypoint script into the container
COPY . .

# Build the Go application (assuming there's a main.go file in your repository)
RUN go build -o smurf main.go && \
    mv smurf /usr/local/bin/ && \
    chmod +x /usr/local/bin/smurf

# Ensure the entrypoint.sh script is executable and starts with a shebang (#!/bin/sh)
RUN sed -i '1s|^|#!/bin/sh\n|' entrypoint.sh && \
    chmod +x entrypoint.sh

# Set the entrypoint for the container to run the shell script
ENTRYPOINT ["/go/src/app/entrypoint.sh"]
