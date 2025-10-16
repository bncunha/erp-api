# ERP API

This repository contains a Go-based ERP API with automated CI/CD and a comprehensive unit test suite that focuses on the application and domain layers.

## Requirements

- Go 1.21 or newer (GitHub Actions uses 1.24).
- A Render deploy hook URL stored as the `RENDER_DEPLOY_HOOK_URL` secret in the GitHub repository.

## Local development

Install dependencies and run the test suite with coverage:

```bash
go mod download
go test ./src/application/... ./src/domain/... -cover
```

To examine a detailed coverage report run:

```bash
go test ./src/application/... ./src/domain/... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Continuous integration and deployment

The workflow defined in [`.github/workflows/deploy.yml`](.github/workflows/deploy.yml) runs on pushes and pull requests to `main`.

1. Checkout the repository and set up Go.
2. Build the project.
3. Run the targeted unit tests with a coverage gate of 90%.
4. Trigger a Render deployment via the `RENDER_DEPLOY_HOOK_URL` secret when the build occurs on `main`.

Ensure that the Render deploy hook URL is configured in the repository secrets before merging to `main` so successful builds automatically deploy.

## Test coverage expectations

The existing suite exercises all code under `src/application` and `src/domain`, and the coverage gate enforces a minimum of 90%. You can review package level coverage locally with:

```bash
go tool cover -func=coverage.out
```

If coverage falls below the threshold, the CI pipeline will fail and block deployment, ensuring that new contributions maintain the required level of automated verification.
