# GitHub Actions Setup Guide

This guide explains how to set up GitHub Actions for CI/CD in this repository.

## Overview

We have two workflows:

1. **CI Workflow** (`.github/workflows/ci.yml`)
   - Runs on every push and pull request
   - Tests code on multiple Go versions
   - Runs linters
   - Builds and tests Docker image
   - Validates builds on multiple platforms

2. **Release Workflow** (`.github/workflows/release.yml`)
   - Triggers on version tags (e.g., `v1.0.0`)
   - Validates semantic versioning
   - Builds binaries for all platforms
   - Creates Docker images
   - Publishes to GitHub Container Registry and optionally Docker Hub
   - Creates GitHub releases with artifacts

## Required Setup

### 1. GitHub Container Registry (GHCR)

**No setup required!** GitHub Actions automatically has access to `GITHUB_TOKEN` which can push to GHCR.

Images will be published to: `ghcr.io/your-username/mcp-email:tag`

### 2. Docker Hub (Optional)

If you want to also publish to Docker Hub:

1. **Create a Docker Hub account** (if you don't have one)
2. **Create an access token**:
   - Go to Docker Hub → Account Settings → Security
   - Click "New Access Token"
   - Give it a name (e.g., "github-actions")
   - Copy the token (you won't see it again!)

3. **Configure in GitHub**:
   - Go to your repository → Settings → Secrets and variables → Actions
   - Add a **Repository Variable**:
     - Name: `DOCKER_HUB_USERNAME`
     - Value: Your Docker Hub username
   - Add a **Repository Secret**:
     - Name: `DOCKER_HUB_TOKEN`
     - Value: Your Docker Hub access token

4. **Verify**: The release workflow will automatically use these if they exist.

### 3. Codecov (Optional)

If you want code coverage reports:

1. Sign up at [codecov.io](https://codecov.io)
2. Add your repository
3. Copy the repository upload token
4. Add it as a GitHub Secret:
   - Name: `CODECOV_TOKEN`
   - Value: Your Codecov token

**Note**: The workflow will work without Codecov, it just won't upload coverage reports.

## Workflow Details

### CI Workflow

**Triggers:**
- Push to `main` or `develop` branches
- Pull requests to `main` or `develop` branches

**Jobs:**
1. **Test**: Runs unit tests on Go 1.21, 1.22, and 1.23
2. **Lint**: Runs golangci-lint
3. **Build**: Builds on Linux, macOS, and Windows
4. **Docker Build**: Builds and tests Docker image

### Release Workflow

**Triggers:**
- Push a tag matching `v*.*.*` (e.g., `v1.0.0`)
- Manual workflow dispatch

**Version Validation:**
- Must be in SemVer format: `vMAJOR.MINOR.PATCH`
- Must be greater than the previous tag
- Fails if version is not incremented

**Jobs:**
1. **Validate Version**: Checks version format and increment
2. **Build Binaries**: Creates binaries for:
   - Linux (amd64, arm64)
   - macOS (amd64, arm64)
   - Windows (amd64, arm64)
3. **Build Docker**: Creates and pushes Docker images
4. **Create Release**: Creates GitHub release with all artifacts

## Creating a Release

### Automatic Release (Recommended)

1. **Create and push a tag**:
   ```bash
   git tag -a v1.0.0 -m "Release v1.0.0"
   git push origin v1.0.0
   ```

2. **Workflow runs automatically** and creates the release.

### Manual Release

1. Go to Actions → Release workflow
2. Click "Run workflow"
3. Enter version tag (e.g., `v1.0.1`)
4. Click "Run workflow"

## Troubleshooting

### Version Validation Fails

**Error**: "Version v1.0.0 is not greater than previous tag v1.0.1"

**Solution**: Use a version that's higher than the previous tag. For example:
- If last tag was `v1.0.0`, next can be `v1.0.1`, `v1.1.0`, or `v2.0.0`
- You cannot go backwards (e.g., `v0.9.0` after `v1.0.0`)

### Docker Hub Push Fails

**Error**: "unauthorized: authentication required"

**Solutions**:
1. Check that `DOCKER_HUB_USERNAME` variable is set correctly
2. Check that `DOCKER_HUB_TOKEN` secret is set correctly
3. Verify the token has the correct permissions
4. If you don't need Docker Hub, the workflow will skip it automatically

### Build Fails

**Error**: "go: cannot find module providing package..."

**Solution**: Run `go mod tidy` locally and commit `go.sum`

### Release Not Created

**Error**: Release appears in Actions but not in Releases

**Solutions**:
1. Check that the workflow completed successfully
2. Verify `GITHUB_TOKEN` has write permissions (should be automatic)
3. Check the workflow logs for errors
4. Ensure the version validation passed

## Testing Locally

You can test the workflows locally using [act](https://github.com/nektos/act):

```bash
# Install act
brew install act  # macOS
# or
curl https://raw.githubusercontent.com/nektos/act/master/install.sh | sudo bash

# Run CI workflow
act push

# Run release workflow (dry run)
act workflow_dispatch -e .github/workflows/release.yml
```

## Next Steps

1. ✅ Set up Docker Hub credentials (optional)
2. ✅ Set up Codecov (optional)
3. ✅ Create your first tag and release
4. ✅ Verify artifacts are created correctly
5. ✅ Test Docker images from GHCR

## Questions?

- Check the workflow files: `.github/workflows/ci.yml` and `.github/workflows/release.yml`
- See `RELEASE.md` for release process details
- Check GitHub Actions logs for detailed error messages

