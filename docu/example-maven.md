# Configuration - Example 'Maven'
## Structure
```
   |
   * master [1.4.0.0]
   |\
   | * feat/myfeature1 [1.4.0-SNAPSHOT]
   | |
   | * feat/myfeature1 [1.4.0-SNAPSHOT]
   |/
   * master [1.3.4.0]
   |
```

## semanticversion.yaml
```
branches:
  - branch_pattern: 'master'
    release_channel: 'FINAL'
    version_pattern: '{major}.{minor}.{patch}.{build}
  
  - branch_pattern: 'feat.*'
    version_pattern: '{major}.{minor}.{patch}-SNAPSHOT

```