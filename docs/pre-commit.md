
## Basics of pre-commit

### Installing pre-commit (per machine)

For macOS

```bash
brew install pre-commit
```

For linux:

```bash
curl -O https://bootstrap.pypa.io/get-pip.py
python3 get-pip.py --user
# Add the executable path, ~/.local/bin, to your PATH variable
pip install pre-commit
```
### Configuring pre-commit for the repository

In the root of the repository create file  `.pre-commit-config.yaml`.  There are [hundreds of supported hooks](https://pre-commit.com/hooks.html).  This just lists the most common ones used across all projects.  If creating a new repository you may want to look at other repositories for other hooks to add.

You can of course remove non relevant hooks (e.g. Remove SwiftFormat if your code is all Python)

```yaml
repos:
-   repo: https://github.com/pre-commit/pre-commit-hooks
    rev: v4.5.0
    hooks:
    -   id: check-json
    -   id: check-toml
    -   id: check-yaml
    -   id: check-xml
    -   id: check-added-large-files
    -   id: check-merge-conflict
    -   id: check-symlinks
    -   id: check-byte-order-marker
    -   id: check-case-conflict

-   repo: https://github.com/Lucas-C/pre-commit-hooks
    rev: v1.5.4
    hooks:
    -   id: forbid-crlf
    -   id: forbid-tabs
        types_or: [objective-c,objective-c++,swift,swiftdeps,java]

-   repo: https://github.com/KarlKraft/header_check
    rev: v1.5.3
    hooks:
    -   id: header_check

-   repo: https://github.com/dnephin/pre-commit-golang
    rev: v0.5.1
    hooks:
    -  id: go-fmt

 - repo: https://github.com/nicklockwood/SwiftFormat
   rev: 0.52.11
   hooks:
    - id: swiftformat

 - repo: https://github.com/adrienverge/yamllint.git
   rev: v1.33.0
   hooks:
    - id: yamllint
      args: [--strict]

```

You can then download the latest version of each hook at any time by running this command

```bash
pre-commit autoupdate
```

Or create a github action to keep your versions up to date automatically.  This repository has an example in ``.github/workflows/pre-commit.yml``

### Run whenever the repo is cloned to install the git hook
```bash
pre-commit install
```

### Using pre-commit
By default pre-commit will only run against modified files.  To run against all files in the repo you can use ``--all-files``

```bash
pre-commit run --all-files
```

If you need to commit without running pre-commit you can pass the `--no-verify` flag

```bash
git commit --no-verify
```

## Useful Links

The label `types_or` defines what file types the hooks rungs again.  When not present the default set is used. This uses the `identify` python library.  A list of all the known types can be found [here](https://github.com/pre-commit/identify/blob/main/identify/extensions.py).

