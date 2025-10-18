# Build System Documentation

This document describes the build system setup and usage for the go-chat project.

## Overview

The project uses:
- **Buf** - For Protocol Buffers management, linting, and code generation
- **Make** - For build automation and task management
- **golangci-lint** - For Go code linting

## Prerequisites

- Go 1.21 or higher
- Make
- curl (for installing golangci-lint)

## Quick Start

### 1. Install Development Tools

```bash
make install-tools
```

This will install:
- `buf` CLI (v1.28.1)
- `protoc-gen-go` (v1.31.0) - Generate Go structs from proto messages
- `protoc-gen-go-grpc` (v1.3.0) - Generate gRPC service stubs
- `protoc-gen-grpc-gateway` (v2.19.0) - Generate REST gateway for gRPC services
- `protoc-gen-openapiv2` (v2.19.0) - Generate OpenAPI/Swagger documentation
- `protoc-gen-doc` (v1.5.1) - Generate markdown documentation from protos
- `protoc-gen-validate` (v1.0.4) - Generate validation code for proto messages
- `golangci-lint` (v1.55.2)

### 2. Verify Installation

```bash
make check-tools
```

### 3. Generate Proto Code

```bash
# Generate code for all services
make proto-gen

# Or generate for a specific service
make proto-gen-auth
make proto-gen-users
make proto-gen-chat
```

## Makefile Structure

The Makefile is modular and split into logical components:

```
Makefile                    # Main entry point
├── scripts/make/tools.mk   # Tool installation and management
├── scripts/make/proto.mk   # Protocol Buffers operations
└── scripts/make/dev.mk     # Development commands
```

## Available Commands

### Tool Management

| Command | Description |
|---------|-------------|
| `make install-tools` | Install all required development tools |
| `make install-buf` | Install only buf CLI |
| `make install-protoc-plugins` | Install only protoc Go plugins |
| `make install-golangci-lint` | Install only golangci-lint |
| `make check-tools` | Verify all tools are installed |

### Proto Management

| Command | Description |
|---------|-------------|
| `make proto-gen` | Generate Go code for all services |
| `make proto-gen-<service>` | Generate Go code for specific service (auth, users, social, chat, notifications, gateway) |
| `make proto-lint` | Lint all proto files |
| `make proto-format` | Format all proto files |
| `make proto-breaking` | Check for breaking changes against main branch |
| `make proto-clean` | Clean generated proto files |

### Development

| Command | Description |
|---------|-------------|
| `make lint` | Run golangci-lint on all services |
| `make lint-fix` | Run golangci-lint with auto-fix |
| `make test` | Run tests for all services |
| `make test-coverage` | Run tests with coverage report |
| `make build` | Build all services |
| `make clean` | Clean all build artifacts |
| `make mod-tidy` | Run go mod tidy for all services |
| `make mod-download` | Download dependencies for all services |

### Help

| Command | Description |
|---------|-------------|
| `make help` | Show all available commands |

## Buf Configuration

Each service has its own buf configuration:

### Service-Level Configuration

**`<service>/buf.gen.yaml`** - Code generation configuration
- Defines output directories:
  - `<service>/pkg/` - Generated Go code
  - `<service>/pkg/docs/` - Proto documentation
  - `swagger/<service>/` - OpenAPI/Swagger specs (root level)
- Configures Go package prefix
- Uses local installed plugins
- Plugins:
  - `protoc-gen-go` - Generates Go structs from proto messages
  - `protoc-gen-go-grpc` - Generates gRPC service stubs
  - `protoc-gen-grpc-gateway` - Generates REST gateway reverse-proxy
  - `protoc-gen-openapiv2` - Generates OpenAPI/Swagger spec for REST API
  - `protoc-gen-doc` - Generates markdown documentation
  - `protoc-gen-validate` - Generates validation code with rules enforcement

### Root Configuration

**`buf.yaml`** - Workspace configuration
- Defines the root workspace containing all service modules
- Configures linting rules: STANDARD, COMMENTS, FILE_LOWER_SNAKE_CASE
- Configures breaking change detection

## Directory Structure

Generated code structure per service:

```
go-chat/
├── swagger/          # Centralized OpenAPI/Swagger specs (gitignored)
│   ├── auth/
│   │   └── auth.swagger.json
│   ├── users/
│   │   └── users.swagger.json
│   ├── social/
│   │   └── social.swagger.json
│   ├── chat/
│   │   └── chat.swagger.json
│   ├── notifications/
│   │   └── notifications.swagger.json
│   └── gateway/
│       └── gateway.swagger.json
└── <service>/
    ├── buf.gen.yaml       # Generation configuration
    ├── proto/             # Proto definitions
    │   └── *.proto
    └── pkg/              # Generated Go code (gitignored)
        ├── *.pb.go       # Protocol buffer messages
        ├── *_grpc.pb.go  # gRPC service stubs
        ├── *.pb.gw.go    # gRPC-Gateway reverse proxy
        ├── *.pb.validate.go  # Validation code
        └── docs/         # Generated documentation
            └── <service>.md
```

## Workflow

### Adding a New Proto File

1. Create proto file in `<service>/proto/`
2. Run `make proto-lint` to check syntax
3. Run `make proto-gen-<service>` to generate Go code
4. Import generated code in your service

### Making Proto Changes

1. Modify proto files
2. Run `make proto-breaking` to check for breaking changes
3. Run `make proto-gen` to regenerate code
4. Update service implementations

### Before Committing

```bash
# Format proto files
make proto-format

# Lint proto files
make proto-lint

# Regenerate code if needed
make proto-gen

# Run linter
make lint

# Run tests
make test
```

## CI/CD Integration

Recommended CI pipeline steps:

```yaml
- make install-tools
- make check-tools
- make proto-lint
- make proto-gen
- make lint
- make test
```

## Tool Versions

Current tool versions are defined in `scripts/make/tools.mk`:

- Buf: 1.28.1
- protoc-gen-go: 1.31.0
- protoc-gen-go-grpc: 1.3.0
- protoc-gen-grpc-gateway: 2.19.0
- protoc-gen-openapiv2: 2.19.0
- protoc-gen-doc: 1.5.1
- protoc-gen-validate: 1.0.4
- golangci-lint: 1.55.2

To update versions, modify the version variables in `tools.mk`.

## Plugin Outputs

Each plugin generates different artifacts:

- **protoc-gen-go**: `<service>/pkg/*.pb.go` - Go message types
- **protoc-gen-go-grpc**: `<service>/pkg/*_grpc.pb.go` - gRPC client/server interfaces
- **protoc-gen-grpc-gateway**: `<service>/pkg/*.pb.gw.go` - HTTP/JSON to gRPC reverse proxy
- **protoc-gen-openapiv2**: `swagger/<service>/*.swagger.json` - OpenAPI v2 spec for REST API (root level)
- **protoc-gen-doc**: `<service>/pkg/docs/*.md` - Human-readable proto documentation
- **protoc-gen-validate**: `<service>/pkg/*.pb.validate.go` - Validation rules enforcement

## Linting Configuration

### golangci-lint

The project uses a comprehensive `.golangci.yml` configuration with:

**Code Quality:**
- Line length: 120 characters
- Function length: max 60 lines, 40 statements
- Cyclomatic complexity: max 15
- Cognitive complexity: max 20

**Enabled Linters (40+):**
- **Core:** errcheck, gosimple, govet, staticcheck, unused
- **Quality:** gocritic, revive, funlen, cyclop, gocognit
- **Security:** gosec, sqlclosecheck
- **Style:** gofmt, gofumpt, goimports, stylecheck, misspell
- **Performance:** prealloc
- **Microservices:** bodyclose, contextcheck, noctx
- And many more...

**Exclusions:**
- Generated proto files (`*.pb.go`, `*.pb.gw.go`, `*.pb.validate.go`)
- Test files have relaxed rules (length, complexity)
- Init functions allowed in `cmd/` directories

**Usage:**
```bash
# Run linter on all services
make lint

# Auto-fix issues where possible
make lint-fix

# Run on specific service
cd <service> && golangci-lint run ./...
```

## Troubleshooting

### "buf: command not found"

Run `make install-buf` or `make install-tools`

### "protoc-gen-go: command not found"

Run `make install-protoc-plugins` or `make install-tools`

### Proto generation fails

1. Check if tools are installed: `make check-tools`
2. Verify proto syntax: `make proto-lint`
3. Ensure `buf.gen.yaml` exists in service directory

### Import paths are wrong

Check `go_package_prefix` in `<service>/buf.gen.yaml` matches your module path.

