# protoget

A very simple tool for fetching and managing proto dependencies. **This doesn't do transitive dependencies (yet?)** I'm not sure if I even ned them, honestly. Here's the problem this solves:

The Github model now dominates the way we manage our projects. If I create a little service, I put the code for that service in its own repo. If I want to run the service, that's great, I just clone it, build it, and run it. But if I want to connect to that service, things get complicated. I just need the client code. But the client code is often for a completely different language than the service is written in. We end up solving that problem by creating `.proto` files to describe the interface to the service in a way where the clients can be generated by `protoc`. Now should my service repo generate the code for every imaginable language? It's nicer to just have the dependency grab the `.proto` files and generate their own client. That's what `protoget` does. For the service repos, it provides a convention for declaring the list of `.proto` files that make up the interface to the service (`protoget.yaml`). The the client repo uses the `protoget` tool to fetch the necessary `.proto` files to generate the client in the language of your choice.

## Installation

```console
go install github.com/kellegous/protoget@latest
```

## Service Usage

The service repos needs to have a `protoget.yaml` file in the root of the repo that looks kind of like this,

```yaml
# protoget.yaml
name: my-service
sources:
  - proto/v1/service.proto
  - proto/v1/types.proto
  - proto/v2/service.proto
```

## Client Usage

**Basic**

This will fetch the `protoget.yaml` file from github.com/kellegous/my-service at the tag `v1.2.3` and write the files listed in the `sources` section to the directory `./external`.

```console
protoget github.com/kellegous/my-service@v1.2.3
```

**Getting a branch**

```console
protoget github.com/kellegous/my-service@main
```

**Getting a specific commit**

```console
protoget github.com/kellegous/my-service@0e7d1e8b6db501aef104f82ebd7f5b21cfa963a1
```

When specifying a dependency, you are required to specify a reference after the `@` symbol. This can be a tag, branch, or commit hash but it isn't optional.

**Changing the destination directory**

```console
protoget github.com/kellegous/my-service@main --destination-directory=./proto
```

**Changing the cache directory**

By default, `protoget` keeps a cache in `$HOME/.cache/protoget`. This can be overridden with the `--cache-directory` flag.

```console
protoget github.com/kellegous/my-service@main --cache-directory=./cache
```

## Authors

- Kelly Norton [kellegous.com](https://kellegous.com/about)
