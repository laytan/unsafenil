# unsafenil

Checks that there is no return of a nil error or false, and a nil/default value before it.

This is a fork of [Antonboom/nilnil](https://github.com/Antonboom/nilnil) which in addition to the
original check of returning a `nil` result and a `nil` error, adds checking for a `nil` return and `false`.

This is a handy addition that catches the common go pattern of returning an 'ok' boolean.

It also extends the checking to not only functions with 2 return values,
but any amount >= 2 that returns an `error` or `bool` as the last return type.

## Installation & usage

### Standalone

```sh
go install github.com/laytan/unsafenil@latest
unsafenil ./...
```

### golangci-lint

The linter is not available from the official golangci-lint package.
These steps add it as a custom linter in your project.

#### Install the package as a submodule

The below command installs a git submodule in the third_party folder recommended
in the go community, [learn more](https://github.com/golang-standards/project-layout).

```sh
git submodule add https://github.com/laytan/unsafenil third_party/unsafenil
```

#### Build the plugin

Builds the .so file for golangci-lint, you probably want to gitignore this binary.

```sh
cd third_party/unsafenil
go build -o ../../unsafenil.so -buildmode=plugin plugin/unsafenil.go
```

#### Add to your config

Add the following snippet to your `.golangci.yml`

```yaml
linters-settings:
    custom:
        unsafenil:
            path: ./unsafenil.so
            description: Checks that there is no return of a nil error or false, and a nil/default value before it
            original-url: github.com/laytan/unsafenil
```

