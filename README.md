# Dep Ensure Cloud Native Buildpack

The Dep Ensure CNB makes use of the [`dep`](https://golang.github.io/dep) tool
to execute the `dep ensure` command in the working directory of the app. For
more info about the `dep ensure` command, see the
[documentation](https://golang.github.io/dep/docs/daily-dep.html#using-dep-ensure)

## Integration

The Dep Ensure CNB does not provide any dependencies. It requires the `dep`
dependency that can be provided by a buildpack like the [Dep
CNB](https://github.com/paketo-buildpacks/dep).

## Usage

To package this buildpack for consumption:

```
$ ./scripts/package.sh
```

## `buildpack.yml` Configurations

This buildpack does not support configurations using `buildpack.yml`.

DO NOT MERGE. Triggering PR actions.
