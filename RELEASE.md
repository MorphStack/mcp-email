# Release Process

This document describes the release process for MCP Email Server.

## Versioning

We use [Semantic Versioning](https://semver.org/) (SemVer) for releases:
- Format: `vMAJOR.MINOR.PATCH` (e.g., `v1.0.0`)
- **MAJOR**: Breaking changes
- **MINOR**: New features (backward compatible)
- **PATCH**: Bug fixes (backward compatible)

## Creating a Release

### 1. Prepare the Release

1. Update `CHANGELOG.md` (if it exists) with the changes in this release
2. Ensure all tests pass: `go test ./...`
3. Ensure the build works: `go build ./cmd/server`

### 2. Create and Push a Tag

Create a new tag with the version number:

```bash
# For a patch release (bug fixes)
git tag -a v1.0.1 -m "Release v1.0.1"

# For a minor release (new features)
git tag -a v1.1.0 -m "Release v1.1.0"

# For a major release (breaking changes)
git tag -a v2.0.0 -m "Release v2.0.0"

# Push the tag
git push origin v1.0.1
```

**Important**: The tag must be in the format `vX.Y.Z` and must be greater than the previous tag, or the release workflow will fail.

### 3. Automatic Release Process

Once you push the tag, the GitHub Actions workflow will automatically:

1. **Validate the version**: Ensure it's in SemVer format and greater than the previous tag
2. **Build binaries**: Create binaries for all platforms:
   - Linux (amd64, arm64)
   - macOS (amd64, arm64)
   - Windows (amd64, arm64)
3. **Build Docker image**: Create and push Docker image to:
   - GitHub Container Registry (ghcr.io)
   - Docker Hub (if configured)
4. **Create GitHub Release**: 
   - Create a new release on GitHub
   - Upload all binary artifacts
   - Generate release notes

### 4. Manual Release (Alternative)

You can also trigger a release manually:

1. Go to the "Actions" tab in GitHub
2. Select the "Release" workflow
3. Click "Run workflow"
4. Enter the version tag (e.g., `v1.0.1`)
5. Click "Run workflow"

## Required Secrets

The following secrets need to be configured in your GitHub repository:

### For GitHub Container Registry (GHCR)
- **GITHUB_TOKEN**: Automatically provided by GitHub Actions (no setup needed)

### For Docker Hub (Optional)
- **DOCKER_HUB_USERNAME**: Your Docker Hub username (configure as a repository variable)
- **DOCKER_HUB_TOKEN**: Your Docker Hub access token (configure as a repository secret)

To configure secrets:
1. Go to your repository settings
2. Navigate to "Secrets and variables" → "Actions"
3. Add the required secrets/variables

## Docker Images

After a successful release, Docker images will be available at:

- **GHCR**: `ghcr.io/your-username/mcp-email:latest` and `ghcr.io/your-username/mcp-email:v1.0.0`
- **Docker Hub** (if configured): `your-username/mcp-email-server:latest` and `your-username/mcp-email-server:v1.0.0`

## Release Artifacts

Each release includes:
- Binary archives for all platforms (tar.gz for Linux/macOS, zip for Windows)
- Docker images for all architectures
- Release notes (auto-generated)

## Troubleshooting

### Version validation fails

If the version validation fails, check:
- The tag format is correct: `vX.Y.Z`
- The version is greater than the previous tag
- The tag exists in the remote repository

### Docker build fails

If the Docker build fails:
- Check that Docker Hub credentials are correct (if using Docker Hub)
- Ensure the Dockerfile is valid
- Check GitHub Actions logs for detailed error messages

### Release not created

If the release is not created:
- Check that the workflow completed successfully
- Verify that the `GITHUB_TOKEN` has write permissions
- Check the GitHub Actions logs for errors

