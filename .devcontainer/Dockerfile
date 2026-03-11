FROM mcr.microsoft.com/devcontainers/go:1.25-bookworm

ARG KIND_VERSION=0.31.0
ARG CONTROLLER_GEN_VERSION=0.20.1
ARG GOFUMPT_VERSION=0.9.2
ARG SETUP_ENVTEST_VERSION=latest

RUN go install sigs.k8s.io/controller-tools/cmd/controller-gen@v${CONTROLLER_GEN_VERSION} \
    && go install mvdan.cc/gofumpt@v${GOFUMPT_VERSION} \
    && go install sigs.k8s.io/controller-runtime/tools/setup-envtest@${SETUP_ENVTEST_VERSION} \
    && go clean -cache -modcache

RUN curl -Lo /usr/local/bin/kind "https://kind.sigs.k8s.io/dl/v${KIND_VERSION}/kind-linux-amd64" \
    && chmod +x /usr/local/bin/kind \
    && curl -s "https://raw.githubusercontent.com/kubernetes-sigs/kustomize/master/hack/install_kustomize.sh" | bash \
    && mv kustomize /usr/local/bin/ \
    && curl -sSfL "https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh" | sh -s -- -b "$(go env GOPATH)/bin" \
    && curl -sSfL "https://github.com/kubernetes-sigs/kubebuilder/releases/latest/download/kubebuilder_linux_$(go env GOARCH)" -o /usr/local/bin/kubebuilder \
    && chmod +x /usr/local/bin/kubebuilder \
    && curl -LsSf https://astral.sh/uv/install.sh | sh \
    && UV_TOOL_BIN_DIR=/usr/local/bin UV_TOOL_DIR=/usr/local/share/uv/tools \
       /root/.local/bin/uv tool install reuse \
    && rm -rf /root/.local/bin/uv /root/.local/bin/uvx /root/.cache/uv
