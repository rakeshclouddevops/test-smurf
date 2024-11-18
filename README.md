
# Smurf

Smurf is a command-line interface built with Cobra, designed to streamline DevOps workflows by providing unified commands for essential tools like Terraform, Helm, and Docker. With Smurf, you can execute Terraform, Helm, and Docker commands seamlessly from a single interface. This CLI simplifies tasks such as environment provisioning, container management, and infrastructure-as-code deployment, improving productivity and minimizing context-switching.

## Features

- **Terraform Command Wrapper (stf):** Run `init`, `plan`, `apply`, `output`, `drift`, `validate`, `destroy`, `format` commands, and `provision`, a combined operation of `init`, `validate`, and `apply`.
- **Helm Command Wrapper (selm):** Run `create`, `install`, `lint`, `list`, `status`, `template`, `upgrade`, `uninstall` commands, and `provision`, a combination of `install`, `upgrade`, `lint`, and `template`.
- **Docker Command Wrapper (sdkr):** Run `build`, `scan`, `tag`, `publish`, `push` commands, and `provision`, a combination of `build`, `scan`, and `publish`.
- **Multicloud Container registry :** Push images from multiple cloud registries like AWS ECR, GCP GCR, Azure ACR, and Docker Hub.Run `smurf sdkr push --help` to push images from the specified registry.
- **Git Integration:** *(Yet to come)*
- **Unified CLI Interface:** Manage multi-tool operations from one interface, reducing the need for multiple command sets.

## Installation

### Prerequisites

- Go 1.20 or higher
- Git
- Terraform, Helm, and Docker Daemon installed and accessible via your PATH

### Installation Steps

1. **Clone the repository:**

   ```bash
   git clone https://github.com/clouddrove/smurf.git
   ```

2. **Change to the project directory:**

   ```bash
   cd smurf
   ```

3. **Build the tool:**

   ```bash
   go build -o smurf .
   ```

   This will build `smurf` in your project directory.

## Usage

### Terraform Commands

Use `smurf stf <command>` to run Terraform commands. Supported commands include:

- **Help:** `smurf stf --help`
- **Initialize Terraform:** `smurf stf init`
- **Generate and Show Execution Plan:** `smurf stf plan`
- **Apply Terraform Changes:** `smurf stf apply`
- **Detect Drift in Terraform State:** `smurf stf drift`
- **Provision Terraform Environment:** `smurf stf provision`

The `provision` command for Terraform performs `init`, `validate`, and `apply`.

### Helm Commands

Use `smurf selm <command>` to run Helm commands. Supported commands include:

- **Help:** `smurf selm --help`
- **Create a Helm Chart:** `smurf selm create`
- **Install a Chart:** `smurf selm install`
- **Upgrade a Release:** `smurf selm upgrade`
- **Provision Helm Environment:** `smurf selm provision --help`

The `provision` command for Helm combines `install`, `upgrade`, `lint`, and `template`.

### Docker Commands

Use `smurf sdkr <command> <flags>` to run Docker commands. Supported commands include:

- **Help:** `smurf sdkr --help`
- **Build an Image:** `smurf sdkr build`
- **Scan an Image:** `smurf sdkr scan`
- **Push an Image:** `smurf sdkr push --help`
- **Provision Registry Environment:** `smurf sdkr provision-hub [flags] `(for Docker Hub)

The `provision-hub` command for Docker combines `build`, `scan`, and `publish`.
The `provision-ecr` command for Docker combines `build`, `scan`, and `publish` for AWS ECR.
THE `provision-gcr` command for Docker combines `build`, `scan`, and `publish` for GCP GCR.
THE `provision-acr` command for Docker combines `build`, `scan`, and `publish` for Azure ACR.


## Contributing

Contributions are welcome! Open issues or pull requests on the [GitHub repository](https://github.com/clouddrove/smurf).

## License

This project is licensed under the Apache License - see the [LICENSE](LICENSE) file for details.
