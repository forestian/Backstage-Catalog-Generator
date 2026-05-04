# cataloggen — Backstage Catalog Generator

Turn a single `services.yaml` into a full set of Backstage `catalog-info.yaml` files — no repetitive YAML, no inconsistent metadata.

## Demo

GIF demo coming soon.

## Quick Demo

```sh
$ cataloggen generate --file services.yaml --output ./catalog
Generated 7 file(s) in ./catalog
  payment-api/catalog-info.yaml
  worker-service/catalog-info.yaml
  systems/payment-system.yaml
  apis/payment-api-openapi.yaml
  resources/payment-db.yaml
  resources/redis-cache.yaml
  locations.yaml
```

## Quick Start

Download a prebuilt binary from [GitHub Releases](https://github.com/forestian/Backstage-Catalog-Generator/releases):

```sh
# Linux/macOS
tar -xzf cataloggen_<version>_linux_amd64.tar.gz
chmod +x cataloggen

# Or build with Go
go install github.com/example/backstage-catalog-generator@latest
```

Run your first generation:

```sh
cataloggen init                                    # scaffold a demo project
cataloggen validate --file ./services.yaml         # check for issues first
cataloggen generate --file ./services.yaml --output ./catalog
```

## Use Cases

- Onboard all services into Backstage without writing YAML by hand
- Enforce consistent metadata (owner, lifecycle, system) across every service
- Validate `services.yaml` in CI before merging catalog changes
- Generate a `locations.yaml` to register everything in Backstage in one step
- Bootstrap a new platform team's service catalog from scratch

## Commands

### `cataloggen init`

Creates a demo project directory with a sample `services.yaml` and pre-generated Backstage catalog files.

```sh
cataloggen init                          # creates ./cataloggen-demo
cataloggen init --output ./my-demo
```

### `cataloggen validate`

Validates a `services.yaml` and prints warnings and errors. Exits non-zero on errors.

```sh
cataloggen validate --file ./services.yaml
cataloggen validate --file ./cataloggen-demo/services.yaml
```

### `cataloggen generate`

Reads `services.yaml` and generates Backstage catalog YAML files. Runs validation first; aborts on errors.

```sh
# generate one file per service (default)
cataloggen generate --file ./services.yaml --output ./catalog

# overwrite existing output
cataloggen generate --file ./services.yaml --output ./catalog --force

# generate all entities in a single file
cataloggen generate --file ./services.yaml --output ./catalog-single --format single
```

#### Flags

| Flag | Default | Description |
|---|---|---|
| `--file` | (required) | Path to services.yaml |
| `--output` | `./catalog` | Output directory |
| `--format` | `files` | `files` or `single` |
| `--owner` | `unknown` | Default entity owner |
| `--system` | `default-system` | Default system name |
| `--lifecycle` | `experimental` | Default lifecycle |
| `--include-location` | `true` | Generate locations.yaml |
| `--force` | `false` | Overwrite existing files |

### `cataloggen version`

```sh
cataloggen version
# cataloggen version 0.1.0
```

## services.yaml format

```yaml
global:
  owner: platform
  lifecycle: production
  system: payment
  namespace: default

systems:
  - name: payment
    title: Payment System
    description: Services related to payment processing
    owner: platform
    domain: commerce

services:
  - name: payment-api
    title: Payment API
    description: Handles payment requests from external clients
    type: service
    lifecycle: production
    owner: payment-team
    system: payment
    repo: https://github.com/example/payment-api
    docs: https://docs.example.com/payment-api
    tags:
      - go
      - api
    annotations:
      github.com/project-slug: example/payment-api
    provides_apis:
      - payment-api
    depends_on:
      - resource:payment-db

apis:
  - name: payment-api
    type: openapi
    lifecycle: production
    owner: payment-team
    system: payment
    definition_path: ./openapi/payment.yaml

resources:
  - name: payment-db
    title: Payment Database
    type: database
    owner: payment-team
    system: payment
```

## Example Output

Generated `payment-api/catalog-info.yaml`:

```yaml
apiVersion: backstage.io/v1alpha1
kind: Component
metadata:
  name: payment-api
  title: Payment API
  description: Handles payment requests from external clients
  annotations:
    github.com/project-slug: example/payment-api
    backstage.io/techdocs-ref: dir:.
    backstage.io/source-location: url:https://github.com/example/payment-api
  tags:
    - go
    - api
  links:
    - url: https://docs.example.com/payment-api
      title: Documentation
spec:
  type: service
  lifecycle: production
  owner: payment-team
  system: payment
  providesApis:
    - payment-api
  dependsOn:
    - resource:payment-db
```

Generated `locations.yaml`:

```yaml
apiVersion: backstage.io/v1alpha1
kind: Location
metadata:
  name: generated-catalog-location
spec:
  type: file
  targets:
    - ./payment-api/catalog-info.yaml
    - ./systems/payment-system.yaml
    - ./apis/payment-api-openapi.yaml
    - ./resources/payment-db.yaml
```

## Generated output (--format files)

```
catalog/
  payment-api/
    catalog-info.yaml      # Component entity
  worker-service/
    catalog-info.yaml      # Component entity
  systems/
    payment-system.yaml    # System entity
  resources/
    payment-db.yaml        # Resource entity
    redis-cache.yaml       # Resource entity
  apis/
    payment-api-openapi.yaml  # API entity
  locations.yaml           # Location entity pointing to all targets
```

## Generated output (--format single)

```
catalog/
  catalog-info.yaml    # all entities, separated by ---
  locations.yaml
```

## Supported Backstage entity kinds

| Kind | Source | Spec fields |
|---|---|---|
| `Component` | `services[]` | type, lifecycle, owner, system, providesApis, consumesApis, dependsOn |
| `System` | `systems[]` | owner, domain |
| `API` | `apis[]` | type, lifecycle, owner, system, definition |
| `Resource` | `resources[]` | type, owner, system |
| `Location` | auto | targets list |

## Defaulting behavior

Fields missing in `services.yaml` are filled from the `global` section, then from CLI flags.

| Field | Default source |
|---|---|
| `service.owner` | `global.owner` → `--owner` |
| `service.lifecycle` | `global.lifecycle` → `--lifecycle` |
| `service.system` | `global.system` → `--system` |
| `service.type` | `service` |
| `api.type` | `openapi` |
| `resource.type` | `other` |
| `global.namespace` | `default` |

## Validation

The `validate` command (and the automatic pre-flight check in `generate`) emits warnings for:

- Service missing `repo`
- Service `repo` using plain HTTP instead of HTTPS
- Service missing `owner` after defaulting
- Service missing `system` after defaulting
- Service missing `tags`
- Service `lifecycle=experimental` in a production-like system
- No systems defined
- No services defined
- API missing `definition_path`
- Resource missing `type`

## Limitations

- Runs locally — does not connect to any Kubernetes cluster or Backstage API.
- Does not read or validate API definition files (referenced by path only).
- Registration in Backstage is manual — register `locations.yaml` via the UI or catalog-import API.
- Review all generated files before committing them to your repository.

## Roadmap

- GitHub Actions workflow for automated catalog generation in CI
- Repository auto-discovery from a GitHub org
- Backstage API integration for automated registration
- OpenAPI file ingestion and inline definition embedding
- Ownership validation against CODEOWNERS

## Install from GitHub Releases

Download a prebuilt binary from the [GitHub Releases page](https://github.com/forestian/Backstage-Catalog-Generator/releases).

**Linux/macOS:**

```sh
tar -xzf cataloggen_<version>_<os>_<arch>.tar.gz
chmod +x cataloggen
./cataloggen version
```

**Windows:**

Download the Windows archive, extract it, and run:

```
cataloggen.exe version
```

---

Part of the [Forestian Cloud Native Toolkit](https://github.com/forestian) — small CLI tools for Kubernetes, observability, GitOps, and platform engineering.
