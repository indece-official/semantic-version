# Configuration

## Example
`semanticversion.yaml`:

```
strategy: LATEST
branches:
  - branch_pattern: 'release.*'
    release_channel: FINAL
    version_pattern: 'v{major}.{minor}.{patch}'

  - branch_pattern: 'testing.*'
    release_channel: BETA
    version_pattern: 'v{major}.{minor}.{patch}-beta.{build}'

  - branch_pattern: 'integration.*'
    release_channel: ALPHA
    version_pattern: 'v{major}.{minor}.{patch}-alpha.{build}'
```

## Documentation
| Field | Required | Values | Description |
| --- | --- | --- | --- |
| strategy | yes | `LATEST`, `CLOSEST`, `OVERALL_LATEST` | (see [Strategies](#strategies)) |
| branches | yes | | |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;branch_pattern | yes | | |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;release_channel | no | `ALPHA`, `BETA`, `GAMMA`, `FINAL` | |
| &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;version_pattern | yes | | |


### Strategies
#### Strategy `LATEST` (**default**)
A child branch will always increment from the *latest version* of the closest 'FINAL'-release branch from the underlying git tree.

```
           * fix [v1.1.1-beta.0]
           |      ^^^^^^
  v1.1.0 * |
  ^^^^^^ | |
v2.0.0 * | |
       | | |
       | | * fix [v1.0.1-beta.0]
        \|/
         * v1.0.0
         |
```

#### Strategy `CLOSEST`
A child branch will always increment from the *closest version* of the closest 'FINAL'-release branch from the underlying git tree.

```
           * fix [v1.0.2-beta.0]
           |      ^^^^^^
  v1.1.0 * |
         | |
v2.0.0 * | |
       | | |
       | | * fix [v1.0.1-beta.0]
        \|/
         * v1.0.0
         | ^^^^^^
```

#### Strategy `OVERALL_LATEST`
A child branch will always increment from the *latest version* of all 'FINAL'-release branches,

```
           * fix [v2.0.1-beta.0]
           |      ^^^^^^
  v1.1.0 * |
         | |
v2.0.0 * | |
^^^^^^ | | |
       | | * fix [v1.0.1-beta.0]
        \|/
         * v1.0.0
         |
```
