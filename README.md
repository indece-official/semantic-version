# semantic-version
Lightweight and configurable semantic-version generator without runtime dependencies (generate version & changelog from git history)

## Features
* Configurable to match many different versioning strategies (alpha-beta-game, maven-like, ...)
* Single binary without dependencies (does *not* need Maven, Gradle, NPM etc.)

## Usage
```
Usage: semantic-version [args] <command>

Args:
  -build int
         (default -1)
  -config string
         (default "./semanticversion.yaml")
  -debug
    
  -git-branch string
    
  -v    Print the version info and exit

Commands:
  generate-config  Generate config file 'semanticversion.yaml'
  get-version      Get the new release version
  get-changelog    Get a changelog with all changes since the last release
```

### Setup
Download the binary and you are ready to go

### Generate config
You can customize the versioning rules by creating the config file `semanticrelease.yaml` in your work directory:

```
> semantic-release generate-config
```

### Git commit format
| Version increment | Prefix | Example |
| --- | --- | --- |
| Major | `break: ` | `break: Changed API model to v2` |
| Minor | `feat: ` | `feat: Added new delete() function` |
| Patch | `fix: ` | `fix: Fixed add() function` |
| Build | - | `Updated README.md` |


Multiple changes can be commited in the same commit message, separated by `;`, e.g.:

```
break: Changed API model to v2; feat: Added new delete() function;
```

### Get version from git-history
```
> semantic-release get-version
```

Output:
```
v1.0.3-feat_apimodel.0
```

### Get changelog from git-history
```
> semantic-release get-version
```

Output:
```
# BREAKING CHANGES
* Changed API model to v2

# Features
* Added new delete() function

# Fixes
* Fixed add() function

```

## Configuration
* [Documentation](./docu/config.md)
* [Example 'Maven'](./docu/example-maven.md)


## Development
### Snapshot build

```
$> make --always-make
```

### Release build

```
$> BUILD_VERSION=1.0.0 make --always-make
```
