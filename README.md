# cataloggen — Backstage Catalog Generator

`cataloggen` is a local CLI tool that reads a `services.yaml` file and generates
[Backstage](https://backstage.io) `catalog-info.yaml` files for Components, Systems,
APIs, and Resources.

## Why

Writing `catalog-info.yaml` by hand for every service is repetitive and produces
inconsistent metadata. `cataloggen` gives platform and SRE teams a single source of
truth and generates deterministic, reviewable Backstage catalog entities from it.

## Install

```sh
git clone https://github.com/example/backstage-catalog-generator
cd backstage-catalog-generator
go build -o cataloggen .
```

Or install directly:

```sh
go install github.com/example/backstage-catalog-generator@latest
```

## Commands

### `cataloggen version`

```sh
cataloggen version
# cataloggen version 0.1.0
```

### `cataloggen init`

Creates a demo project directory with a sample `services.yaml` and pre-generated
Backstage catalog files.

```sh
cataloggen init --output ./cataloggen-demo
```

### `cataloggen validate`

Validates a `services.yaml` file and prints warnings and errors.

```sh
cataloggen validate --file ./services.yaml
cataloggen validate --file ./cataloggen-demo/services.yaml
```

### `cataloggen generate`

Reads `services.yaml` and generates Backstage catalog YAML files.

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

Fields missing in `services.yaml` are filled from the `global` section, then from
CLI flags.

| Field | Default source |
|---|---|
| `service.owner` | `global.owner` → `--owner` |
| `service.lifecycle` | `global.lifecycle` → `--lifecycle` |
| `service.system` | `global.system` → `--system` |
| `service.type` | `service` |
| `api.type` | `openapi` |
| `resource.type` | `other` |
| `global.namespace` | `default` |

## Validation warnings

The `validate` command emits warnings (non-blocking) for:

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

- API definition files are referenced by path but not read or embedded.
- No Backstage API integration — registration is manual.
- No repository scanning or auto-discovery.
- Review all generated files before registering them in Backstage.

## Roadmap (not implemented)

- GitHub Action for CI-based catalog generation
- Repository auto-discovery
- Backstage API integration
- OpenAPI file ingestion
- TechDocs generation support
- Ownership validation against CODEOWNERS
