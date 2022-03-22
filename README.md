# Pi App Depoyer Action

This action deploys applications using [pi-app-deployer](https://github.com/andrewmarklloyd/pi-app-deployer). The action will trigger a deploy then poll the API checking the deploy status until a success or a timeout.

## Inputs

### `repoName`

**Required** The Github repo name including the org or user name.

### `manifestName`

**Required** The name of the pi-app-deployer/v1 Manifest. Must be defined in the root of the repo in a file named `.pi-app-deployer.yaml`.

## Environment Variables

### `PI_APP_DEPLOYER_API_KEY`

**Required** Authenticates with the pi-app-deployer API server.

### `GITHUB_SHA`

**Required** This is inherited from all Github Actions.

## Example usage

```        
uses: andrewmarklloyd/pi-app-deployer-action@v1
env:
  PI_APP_DEPLOYER_API_KEY: ${{ secrets.PI_APP_DEPLOYER_API_KEY }}
with:
  repoName: andrewmarklloyd/pi-test
  manifestName: pi-test-arm
```
