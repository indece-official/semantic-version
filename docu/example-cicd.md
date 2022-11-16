# Configuration - Example 'CICD'
## Structure
```
  |
  * main [v1.2.0.0]
  |\
  | * PR-2 [v1.2.0-PR-2.0]
  |  \
  |   * test [v1.2.0-test.0]
  |   |\
  |   | * PR-1 [v1.2.0-PR-1.0]
  |   |  \
  |   |   * dev [v1.2.0-dev.0]
  |   |   |
```

## semanticversion.yaml
```
strategy: LATEST
branches:
  - branch_pattern: master
    release_channel: FINAL
    version_pattern: v{major}.{minor}.{patch}.{build}

  - branch_pattern: test
    version_pattern: v{major}.{minor}.{patch}-test.{build}

  - branch_pattern: dev
    version_pattern: v{major}.{minor}.{patch}-dev.{build}

  - branch_pattern: PR-*
    version_pattern: v{major}.{minor}.{patch}-{branch}.{build}

```