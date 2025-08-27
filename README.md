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
      args: [--owner, "Karl Kraft", --license, arr,--autodate]
```

| Flag                    | Description |
| -----------             | ----------- |
| --license [arr,apache]  | Sets the expected license in the copyright block. The default is arr (All Rights Reserved.)  apache is for the Apache 2.0 license      |
| --autodate              | Automatically update the date in copyright headers. |
| --owner              | This can appear multiple times.  When looking at the copyright line it will ensure that the file is copywritten by one of the listed owners. |


### Future Plans

* Scan for and fix dates in Xcode project Info.plist files

