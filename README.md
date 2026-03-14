# header_check
This program replaces a pre-commit SVN hook that I used long ago for ensuring that all common source files have copyright messages in the header.  

It is used with the [pre-commit](https://pre-commit.com) software package.  If you aren't familiar with pre-commit you can find my quick notes [here](./docs/pre-commit.md).

The headers need to start on the first line of the file and be in this form:
```
//
//  filename.suffix
//  project_name
//
//  Created by FIRST LAST on MM/DD/YYYY
//  Copyright 2026 OWNER. All rights reserved
//
```
The project_name is typically the name of the repository, but can also be the name of any enclosing folder.  For instance if the file "foo.go" is located in "https://github.com/KarlKraft/MySuperProject/toolkits/network/foo.go", the project_name can be toolkits, network, or MySuperProject.

You can look at the .go files in this project for examples of how these headers look.

## A special note on Package.swift

Package.swift requires that the first two lines be "comments" that identify the minimum tool version.  In this case these two lines can appear before the copyright block.  This is allowed for any file named Package.swift

Example:
```
// swift-tools-version: 6.0
// The swift-tools-version declares the minimum version of Swift required to build this package.
//
//  Package.swift
//  ColorPopoverWell
//
//  Created by Karl Kraft on 12/25/22
//  Copyright 2022-2026 Karl Kraft. Licensed under Apache License, Version 2.0
//
```


## Configuration
When configuring there are flags you can pass to header_check:

```
 - repo: https://github.com/KarlKraft/header_check
   rev: v1.5.3
   hooks:
    - id: header_check
      args: [--owner, "Karl Kraft", --license, arr,--autodate]
```

| Flag                   | Description                                                                                                                                  |
|------------------------|----------------------------------------------------------------------------------------------------------------------------------------------|
| --license [arr,apache] | Sets the expected license in the copyright block. The default is arr (All Rights Reserved.)  apache is for the Apache 2.0 license            |
| --autodate             | Automatically update the date in copyright headers.                                                                                          |
| --owner                | This can appear multiple times.  When looking at the copyright line it will ensure that the file is copyrighted by one of the listed owners. |

There is an additional flag "--infoplist" that should scan for copyrights in Info.plist files and keep them up to date, but this is not currently complete.

### Future Plans

* Scan for and fix dates in Xcode project Info.plist files

