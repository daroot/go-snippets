---
title: Go Snippets
description: Small Go sub-packages.
author: me
tags: [ go, snippets, code, reuseable ]
date: 2023-06-01T15:54:12-0700
lastmod: 2023-10-08T09:22:55+0000
draft: false
license: "[Grug 1-Clause License](./LICENSE)"
---

This repo serves as a central place from which to copy
Go "modules" that I frequently use
that are either small enough
to not be worth making into separate 3rd party dependencies,
or which need minor adapation into a given project.

Since these snippets are meant to be copied and adapted to fit a new project,
they occasionally have things hard coded that in a "proper" library package
would be configured via option patterns or structs or the like.
Check that any hard coded defaults make sense for the destination project.
Also update all instances of `myapp` and `importfromprojectlocally`,
as those are used in instances where some of these snippets depend on each other,
such as `testgoldenproto` using `testgolden`.

- **buildinfo**: populate a struct containing git commit hash and date of build.
- **consterr**: string-based errors, instead of errors.New(), so you can make them `const`
- **envflag**: set flag variables via ENV without any extra third party dependencies like viper.
- **slogext**: various slog helpers for contexts, errors, and time formats.
- **httptools**: http.Client constructor and http.Handler serving with graceful shutdowns.
- **skeleton**: new project templates
- **testbuffer**: a sync.Mutex locked buffer for use in tests with goroutines.
- **testgolden**: test helpers for comparing results to a golden file and updating said files.
- **testgoldenproto**: as above, but includes protobuf comparison support

## TODOS

- slogext needs a ReplaceAttr chainer/composer.
