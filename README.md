# header_check
This program replaces the pre-commit SVN hook that was used in the olden days for ensuring that all common source files have copyright messages in the header.  

It is used with the [pre-commit](https://pre-commit.com) software package.

## Basics of Usage
###Run once per machine (macOS).  
```bash
brew install pre-commit
```
### Run once per machine (linux)
```bash
curl -O https://bootstrap.pypa.io/get-pip.py
python3 get-pip.py --user
#Add the executable path, ~/.local/bin, to your PATH variable
pip install pre-commit
```

###Configure once per repo

In root of repo create file  `.pre-commit-config.yaml`.  There are [hundreds of supported hooks](https://pre-commit.com/hooks.html).  This just lists the most common ones used across all projects.  If creating a new repo you may want to look at other repos for other hooks to add.

```yaml
repos:
-   repo: https://github.com/pre-commit/pre-commit-hooks
    rev: v4.4.0
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
    rev: v1.3.1
    hooks:
    -   id: forbid-crlf
    -   id: forbid-tabs
        types_or: [objective-c,objective-c++,swift,swiftdeps,java]

-   repo: https://github.com/KarlKraft/header_check
    rev: v1.3.0
    hooks:
    -   id: header_check

-   repo: https://github.com/dnephin/pre-commit-golang
    rev: v0.5.1
    hooks:
    -  id: go-fmt

```

You can then download the latest version of each hook at any time by running this command

```bash
pre-commit autoupdate
```

### Run whenever the repo is cloned to install the git hook
```bash
pre-commit install
```

### Using pre-commit
By default pre-commit will only run against modified files.  To run against all files in the repo you can use:

```bash
pre-commit run --all-files
```

If you need to commit without running pre-commit you can pass the `--no-verify` flag

```bash
git commit --no-verify
```

## Useful Links

The label `types_or` defines what file types the hooks rungs again.  When not present the default set is used. This uses the `identify` python library.  A list of all the known types can be found [here](https://github.com/pre-commit/identify/blob/main/identify/extensions.py):


# Future Plans

* If you don't run pre-commit install when you clone a repo the hooks are never run.  There is no way to force clients to run the hooks. At some point we need to apply the hooks at the pull-request level using github actions. [https://pre-commit.ci/](https://pre-commit.ci/)


* Add support for markdown checking
