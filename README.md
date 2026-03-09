[![REUSE status](https://api.reuse.software/badge/github.com/openmcp-project/service-provider-aso)](https://api.reuse.software/info/github.com/openmcp-project/service-provider-aso)

# service-provider-aso

## 📖 About this project

An [OpenMCP](https://openmcp-project.github.io/) Service Provider that enables platform owners to offer [Azure Service Operator (ASO)](https://azure.github.io/azure-service-operator/) as a service to end users.

**⚠️ Experimental Status**: This service provider is in experimental state and might not be feature complete. Use with caution in production environments.

### ✨ What it provides

This provider brings two core APIs to your OpenMCP platform:

1. **`AzureServiceOperator`** (end-user facing) - Allows end users to request managed ASO instances, enabling them to provision and manage Azure resources through Kubernetes manifests
2. **`ProviderConfig`** (platform owner facing) - Enables platform owners to configure how ASO is offered, including image locations, version constraints, and environment-specific settings

The provider's reconciler manages the full lifecycle of ASO installations, deploying them into managed control planes or workload clusters based on your platform configuration. For more details on service providers, see the [OpenMCP documentation](https://openmcp-project.github.io/docs/developers/service-providers).

## 🏗️ Requirements and Setup

### Building the Image

To build the Docker image, use the following command:

```shell
task build:img:build
```

### Running End-to-End Tests

```shell
task test-e2e
```

## 🧑‍💻 Development and Debugging

The repository includes a [.vscode](.vscode) folder with a [launch.json](.vscode/launch.json) containing two launch configurations for debugging:

1. **Debug Provider**: Debug the service provider itself during development
2. **Debug E2E Tests**: Debug the end-to-end tests - this is particularly useful as you can stop at any point in time and play around with the OpenMCP setup that was created

The E2E debugging configuration allows you to pause test execution and inspect or interact with the running OpenMCP environment, making it invaluable for troubleshooting and understanding the system behavior.

## ⚙️ CLI Flags

### Service Provider Runtime Flags

The generated service provider supports the following runtime flags:

- `--verbosity`: Logging verbosity level (see [controller-runtime logging](https://github.com/kubernetes-sigs/controller-runtime/blob/main/TMP-LOGGING.md))
- `--environment`: Name of the environment (required for operation)
- `--provider-name`: Name of the provider resource (required for operation)
- `--metrics-bind-address`: Address for the metrics endpoint (default: `0`, use `:8443` for HTTPS or `:8080` for HTTP)
- `--health-probe-bind-address`: Address for health probe endpoint (default: `:8081`)
- `--leader-elect`: Enable leader election for controller manager (default: `false`)
- `--metrics-secure`: Serve metrics endpoint securely via HTTPS (default: `true`)
- `--enable-http2`: Enable HTTP/2 for metrics and webhook servers (default: `false`)

For a complete list of available flags, run the generated binary with `-h` or `--help`.

## ❤️ Support, Feedback, Contributing

This project is open to feature requests/suggestions, bug reports etc. via [GitHub issues](https://github.com/openmcp-project/service-provider-aso/issues). Contribution and feedback are encouraged and always welcome. For more information about how to contribute, the project structure, as well as additional contribution information, see our [Contribution Guidelines](CONTRIBUTING.md).

## 🔐 Security / Disclosure

If you find any bug that may be a security problem, please follow our instructions at [in our security policy](https://github.com/openmcp-project/service-provider-aso/security/policy) on how to report it. Please do not create GitHub issues for security-related doubts or problems.

## 🤝 Code of Conduct

We as members, contributors, and leaders pledge to make participation in our community a harassment-free experience for everyone. By participating in this project, you agree to abide by its [Code of Conduct](https://github.com/SAP/.github/blob/main/CODE_OF_CONDUCT.md) at all times.

## 📋 Licensing

Copyright 2025 SAP SE or an SAP affiliate company and service-provider-aso contributors. Please see our [LICENSE](LICENSE) for copyright and license information. Detailed information including third-party components and their licensing/copyright information is available [via the REUSE tool](https://api.reuse.software/info/github.com/openmcp-project/service-provider-aso).
