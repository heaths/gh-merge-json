# Merge JSON from GitHub CLI

A [GitHub CLI](https://github.com/cli/cli) to merge JSON responses to return valid JSON when passing `--paginate` to `gh api`:

```powershell
$issues = gh api 'repos/{owner}/{repo}/issues' --paginate | gh merge-json | ConvertFrom-Json
```

## Install

Make sure you have version 2.0 or [newer](https://github.com/cli/cli/releases/latest) of the GitHub CLI installed.

```bash
gh extension install heaths/gh-merge-json
```

## Details

The GitHub CLI can return paginated JSON responses from both REST and GraphQL APIs;
however, due to bug [#1268](https://github.com/cli/cli/issues/1268) when `--paginate` is passed to `gh api` it does not
return value JSON:

```json
{
  "data": {
    "viewer": {
      "repositories": {
        "nodes": [
          {
            "nameWithOwner": "owner/repo1"
          },
          {
            "nameWithOwner": "owner/repo2"
          }
        ],
        "pageInfo": {
          "hasNextPage": true,
          "endCursor": "Y3Vyc29yOnYyOpHOBAaq0A=="
        }
      }
    }
  }
}
{
  "data": {
    "viewer": {
      "repositories": {
        "nodes": [
          {
            "nameWithOwner": "owner/repo3"
          },
          {
            "nameWithOwner": "owner/repo4"
          }
        ],
        "pageInfo": {
          "hasNextPage": true,
          "endCursor": "Y3Vyc29yOnYyOpHOBzdvuQ=="
        }
      }
    }
  }
}
```

Until [PR #5652](https://github.com/cli/cli/pull/5652) or some other solution is merged, you can use this extension
without having to install additional system utilities like `jq` which may not be available for your platform or
otherwise difficult to obtain e.g., requiring a bunch of prerequisites.

## License

Licensed under the [MIT](LICENSE.txt) license.
