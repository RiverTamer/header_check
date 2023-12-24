# header_check
This program replaces a pre-commit SVN hook that I used long ago for ensuring that all common source files have copyright messages in the header.  

It is used with the [pre-commit](https://pre-commit.com) software package.  If you aren't familiar with pre-commit you can find my quick notes [here](./docs/pre-commit.md).

## Configuration
When configuring there are flags you can pass to header_check:

```
 - repo: https://github.com/KarlKraft/header_check
   rev: v1.5.3
   hooks:
    - id: header_check
      args: [--license, arr]
```

| Flag                    | Description |
| -----------             | ----------- |
| --license [arr,apache]  | Sets the expected license in the copyright block. The default is arr (All Rights Reserved.)  apache is for the Apache 2.0 license      |
| --autodate              | Automatically update the date in copyright headers. |
| --infoplist              | If autodate is enabled, then scan for Info.plist files and update the copyright dates.  If autodate is not enabled then scan for the file and validate, but do not auto update.|


### Future Plans

* If you don't run pre-commit install when you clone a repo the hooks are never run.  There is no way to force clients to run the hooks. At some point we need to apply the hooks at the pull-request level using github actions. [https://pre-commit.ci/](https://pre-commit.ci/)