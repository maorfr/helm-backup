# Helm Backup Plugin

This is a Helm plugin which performs backup/restore of releases in a namespace to/from a file

## Usage

backup releases from namespace to file

```
$ helm backup [flags] NAMESPACE
```

restore releases from file to namespace

```
$ helm backup [flags] NAMESPACE --restore
```

### Flags:

```
      --file string        file name to use (.tgz file). If not provided - will use <namespace>.tgz
  -h, --help               help for backup
  -l, --label string       label to select tiller resources by (default "OWNER=TILLER")
  -r, --restore            restore instead of backup
  -t, --tiller-ns string   namespace of Tiller (default "kube-system")

```

## Install

```
$ helm plugin install https://github.com/maorfr/helm-backup
```

The above will fetch the latest binary release of `helm backup` and install it.

### Developer (From Source) Install

If you would like to handle the build yourself, instead of fetching a binary,
this is how recommend doing it.

First, set up your environment:

- You need to have [Go](http://golang.org) installed. Make sure to set `$GOPATH`
- If you don't have [Dep](https://github.com/golang/dep) installed, this will install it into
  `$GOPATH/bin` for you.

Clone this repo into your `$GOPATH`. You can use `go get -d github.com/maorfr/helm-backup`
for that.

```
$ cd $GOPATH/src/github.com/maorfr/helm-backup
$ make bootstrap build
$ HELM_PUSH_PLUGIN_NO_INSTALL_HOOK=1 helm plugin install $GOPATH/src/github.com/maorfr/helm-backup
```

That last command will skip fetching the binary install and use the one you
built.
