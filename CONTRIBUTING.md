# Contributing to QRAP

Thank you for your interest in contributing to QRAP (Quantum Risk Assessment Platform). This document provides guidelines and instructions for contributing to the project.

## Table of Contents

- [Welcome](#welcome)
- [Code of Conduct](#code-of-conduct)
- [How to Report Bugs](#how-to-report-bugs)
- [How to Suggest Features](#how-to-suggest-features)
- [Development Workflow](#development-workflow)
- [Branch Naming](#branch-naming)
- [Commit Message Format](#commit-message-format)
- [Pull Request Guidelines](#pull-request-guidelines)
- [Code Review Process](#code-review-process)
- [Release Process](#release-process)
- [Getting Help](#getting-help)

---

## Welcome

QRAP is an open-source project and we welcome contributions from the community. Whether you are fixing a bug, adding a feature, improving documentation, or writing tests, your contribution is valued.

Before contributing, please take a moment to read through this guide to understand our processes and expectations. Following these guidelines helps maintain a high-quality codebase and makes the review process smoother for everyone.

---

## Code of Conduct

This project follows the [Contributor Covenant Code of Conduct](https://www.contributor-covenant.org/version/2/1/code_of_conduct/). By participating, you are expected to uphold this code. Please report unacceptable behavior to the project maintainers.

---

## How to Report Bugs

If you find a bug, please open a GitHub Issue with the following information:

1. **Title**: A clear, concise summary of the bug.
2. **Environment**: Operating system, browser (if applicable), Go/Python/Node.js versions, and QRAP version or commit hash.
3. **Steps to reproduce**: A minimal set of steps that reliably reproduce the issue.
4. **Expected behavior**: What you expected to happen.
5. **Actual behavior**: What actually happened, including error messages, log output, or screenshots.
6. **Additional context**: Any other information that might help diagnose the issue (configuration, database state, network conditions).

**Security vulnerabilities** should NOT be reported as public issues. See [SECURITY.md](SECURITY.md) for responsible disclosure instructions.

---

## How to Suggest Features

We welcome feature suggestions. To propose a new feature:

1. **Search existing issues** to check if the feature has already been requested.
2. **Open a new issue** with the `enhancement` label.
3. **Describe the feature**: Explain what the feature does and why it is useful.
4. **Provide context**: Describe the use case, who would benefit, and how it fits into the existing architecture.
5. **Suggest an approach** (optional): If you have ideas about implementation, include them. This is not required.

For significant features that change the architecture or public API, please open a discussion first before writing code.

---

## Development Workflow

### 1. Fork the Repository

Fork the QRAP repository to your GitHub account and clone your fork locally:

```bash
git clone https://github.com/your-username/qrap.git
cd qrap
git remote add upstream https://github.com/quantun/qrap.git
```

### 2. Create a Branch

Create a feature branch from the latest `main`:

```bash
git fetch upstream
git checkout -b feature/your-feature-name upstream/main
```

### 3. Set Up the Development Environment

Follow the instructions in [docs/DEVELOPMENT.md](docs/DEVELOPMENT.md) to set up your local environment.

### 4. Write Code

Make your changes, following the code style guidelines described in the development guide. Key points:

- Write tests for new functionality.
- Update documentation if your change affects public APIs or user-facing behavior.
- Keep changes focused. If you find an unrelated issue, open a separate PR for it.

### 5. Run Tests and Linters

Before committing, ensure all tests pass and linters are clean:

```bash
make test
make lint
```

### 6. Commit Your Changes

Write clear, descriptive commit messages following the [Conventional Commits](#commit-message-format) format.

### 7. Push and Open a Pull Request

```bash
git push origin feature/your-feature-name
```

Open a pull request against the `main` branch of the upstream repository.

---

## Branch Naming

Use the following prefixes for branch names:

| Prefix | Purpose | Example |
|---|---|---|
| `feature/` | New features | `feature/hndl-batch-analysis` |
| `fix/` | Bug fixes | `fix/jwt-expiry-validation` |
| `docs/` | Documentation changes | `docs/deployment-guide` |
| `test/` | Test additions or fixes | `test/ml-scoring-edge-cases` |
| `chore/` | Build, CI, or maintenance tasks | `chore/upgrade-go-1.24` |

Use lowercase with hyphens for multi-word names. Keep branch names concise but descriptive.

---

## Commit Message Format

This project uses [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/). Every commit message must follow this format:

```
<type>(<scope>): <description>

[optional body]

[optional footer(s)]
```

### Types

| Type | Description |
|---|---|
| `feat` | A new feature |
| `fix` | A bug fix |
| `docs` | Documentation-only changes |
| `test` | Adding or correcting tests |
| `chore` | Build process, CI, or auxiliary tool changes |
| `refactor` | Code change that neither fixes a bug nor adds a feature |
| `perf` | Performance improvement |
| `style` | Formatting, whitespace, or other non-functional changes |
| `ci` | Changes to CI/CD configuration |

### Scopes

Use the component name as the scope: `api`, `ml`, `web`, `db`, `infra`, `docs`.

### Examples

```
feat(api): add batch assessment endpoint

Allows creating multiple assessments in a single request.
Includes transaction support for atomic batch operations.

fix(ml): correct HNDL score normalization for short shelf lives

The Mosca inequality calculation produced negative values when
shelf_life_years was less than migration_time_years. Clamp the
output to the [0, 100] range.

docs: update deployment guide with Kubernetes manifests

test(api): add integration tests for rate limiting middleware

chore(ci): upgrade GitHub Actions to use Node.js 22
```

### Breaking Changes

If your commit introduces a breaking change, add `BREAKING CHANGE:` in the footer:

```
feat(api)!: rename /assessments endpoint to /risk-assessments

BREAKING CHANGE: The /api/v1/assessments endpoint has been renamed
to /api/v1/risk-assessments. Update all client integrations.
```

---

## Pull Request Guidelines

When opening a pull request:

1. **Title**: Use the same format as commit messages (e.g., `feat(api): add batch assessment endpoint`).
2. **Description**: Clearly explain what the PR does and why. Include:
   - Summary of changes
   - Motivation and context
   - How it was tested
   - Screenshots (for UI changes)
3. **Size**: Keep PRs focused. A PR should address one concern. Large PRs are harder to review and more likely to introduce issues.
4. **Tests**: All new functionality must include tests. All existing tests must pass.
5. **No breaking changes** without prior discussion. If a breaking change is necessary, open an issue to discuss the approach first.
6. **Documentation**: Update relevant documentation if your change affects APIs, configuration, or user-facing behavior.
7. **Clean history**: Squash work-in-progress commits before requesting review. Each commit in the final PR should be meaningful.
8. **CI must pass**: All GitHub Actions checks must be green before a PR can be merged.

### PR Template

When you open a PR, please include the following sections:

```markdown
## What does this PR do?

<!-- A brief description of the changes -->

## Why is this change needed?

<!-- Motivation, context, or link to an issue -->

## How was this tested?

<!-- Steps to test, test output, or screenshots -->

## Checklist

- [ ] Tests added/updated
- [ ] Documentation updated (if applicable)
- [ ] Linters pass (`make lint`)
- [ ] All tests pass (`make test`)
- [ ] No breaking changes (or discussed in an issue first)
```

---

## Code Review Process

1. **All PRs require at least one approving review** from a project maintainer before merging.
2. Reviewers will check for:
   - Correctness and completeness
   - Test coverage
   - Code style and consistency
   - Security considerations (especially for auth, input validation, and database queries)
   - Performance implications
   - Documentation accuracy
3. **Address all review comments**. If you disagree with a suggestion, explain your reasoning in the review thread.
4. Once approved and CI passes, a maintainer will merge the PR using squash-and-merge.
5. The merge commit message will follow the Conventional Commits format.

---

## Release Process

QRAP follows [Semantic Versioning](https://semver.org/spec/v2.0.0.html):

- **MAJOR** (1.0.0): Incompatible API changes
- **MINOR** (0.2.0): New functionality, backwards compatible
- **PATCH** (0.1.1): Bug fixes, backwards compatible

The release process is managed by the project maintainers:

1. A release branch is created from `main`.
2. The `CHANGELOG.md` is updated with all changes since the last release.
3. Version numbers are updated across the codebase.
4. The release is tagged and published on GitHub with release notes.
5. Container images are built and pushed to the container registry.

Contributors do not need to manage version numbers or changelog entries; the maintainers handle this during the release cycle.

---

## Getting Help

If you need help with your contribution or have questions:

- **GitHub Discussions**: Open a discussion for questions, ideas, or general conversation about the project.
- **GitHub Issues**: Search existing issues or open a new one for specific bugs or feature requests.
- **Development Guide**: See [docs/DEVELOPMENT.md](docs/DEVELOPMENT.md) for detailed setup and development instructions.

We aim to respond to issues and discussions within a few business days. Thank you for contributing to QRAP.
