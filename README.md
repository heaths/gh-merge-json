# Merge JSON from GitHub CLI

A [GitHub CLI](https://github.com/cli/cli) to merge JSON responses to return valid JSON when passing `--paginate` to `gh api`:

```powershell
$issues = gh api 'repos/{owner}/{repo}/issues' --paginate | gh merge-json | ConvertFrom-Json
```

## Install

Make sure you have version 2.0 or [newer](https://github.com/cli/cli/releases/latest) of the GitHub CLI installed.

```bash
gh extension install --force heaths/gh-merge-json
```

The `--force` flag I [introduced](https://github.com/cli/cli/pull/7173) in `gh` [version 2.25](https://github.com/cli/cli/releases/tag/v2.25.0)
to better support automation. It will install or upgrade an extension, or do nothing if the latest extension version is
already installed.

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

### Formatting

Like several `gh` commands, you can pass `--template` (or `-t`) and a [Go template](https://pkg.go.dev/text/template)
using most functions described in `gh help formatting` - specifically, those defined [in source](https://github.com/cli/go-gh/blob/ee75bbcac35dc7d72f0a39d52e74ded343d17810/pkg/template/template.go#L58).
Output is always written to standard output.

### Teeing

If you would like to write formatted output e.g., colorful, indented JSON, or formatted using a `--template` but save
compressed JSON to a file for further processing by other tools, you can pass `--tee` and a file path.

```bash
gh api graphql --paginate -f query='
    query($endCursor: String) {
      viewer {
        repositories(first: 100, after: $endCursor) {
          nodes { nameWithOwner }
          pageInfo {
            hasNextPage
            endCursor
          }
        }
      }
    }
  ' | gh merge-json --tee repos.json --template '{{range .data.viewer.repositories.nodes}}{{.nameWithOwner}}{{"\n"}}{{end}}'
```

## License

Licensed under the [MIT](LICENSE.txt) license.
