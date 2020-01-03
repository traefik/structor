# Messor Structor: Manage multiple documentation versions with Mkdocs.

[![GitHub release](https://img.shields.io/github/release/containous/structor.svg)](https://github.com/containous/structor/releases/latest)
[![Build Status](https://travis-ci.org/containous/structor.svg?branch=master)](https://travis-ci.org/containous/structor)

Structor use git branches to create the versions of a documentation, only works with Mkdocs.

To use Structor a project must respect [semver](https://semver.org) and creates a git branch for each `MINOR` and `MAJOR` version.

Created for [Traefik](https://github.com/containous/traefik) and used by:

* [Traefik](https://github.com/containous/traefik) for https://docs.traefik.io
* [JanusGraph](https://github.com/JanusGraph/janusgraph) for https://docs.janusgraph.org
* [ONOS Project](https://github.com/onosproject) for https://docs.onosproject.org

## Prerequisites

* [git](https://git-scm.com/)
* [Docker](https://www.docker.com/)
* `requirements.txt`, `mkdocs.yml`, and a Dockerfile.

## Description

For the following git graph:

```
* gaaaaaaag - (branch master) commit
| * faaaaaaaf - (branch v1.2) commit
| * eaaaaaaae - commit
|/
* haaaaaaah - commit
| * kaaaaaaak - (branch v1.1) commit
| * jaaaaaaaj - commit
|/
* iaaaaaaai - commit
| * daaaaaaad - (branch v1.0) commit
| * caaaaaaac - commit
|/
* baaaaaaab - commit
* aaaaaaaaa - initial commit
```

Structor generates the following files structure:

```
. (latest, branch v1.2)
├── index.html
├── ...
├── v1.0 (branch v1.0)
│   ├── index.html
│   └── ...
├── v1.1 (branch v1.1)
│   ├── index.html
│   └── ...
└── v1.2 (branch v1.2)
    ├── index.html
    └── ...
```

If the content from `.` is served on `mydoc.com`, documentation will be available at the following URLs:

- `http://mydoc.com` (latest, branch v1.2)
- `http://mydoc.com/v1.0` (branch v1.0)
- `http://mydoc.com/v1.1` (branch v1.1)
- `http://mydoc.com/v1.2` (branch v1.2)

The multi version menu is created from templates provided by the following options:

- `--menu.js-url` (or `--menu.js-file`)
- `--menu.css-url` (or `--menu.css-file`)

## Configuration

```yaml
Messor Structor: Manage multiple documentation versions with Mkdocs.

Usage:
  structor [flags]
  structor [command]

Available Commands:
  help        Help about any command
  version     Display version

Flags:
      --debug                    Debug mode.
      --dockerfile-name string   Search and use this Dockerfile in the repository (in './docs/' or in './') for building documentation. (default "docs.Dockerfile")
  -d, --dockerfile-url string    Use this Dockerfile when --dockerfile-name is not found. Can be a file path. [required]
      --exclude strings          Exclude branches from the documentation generation.
      --exp-branch string        Build a branch as experimental.
      --force-edit-url           Add a dedicated edition URL for each version.
  -h, --help                     help for structor
      --image-name string        Docker image name. (default "doc-site")
      --menu.css-file string     File path of the template of the CSS file use for the multi version menu.
      --menu.css-url string      URL of the template of the CSS file use for the multi version menu.
      --menu.js-file string      File path of the template of the JS file use for the multi version menu.
      --menu.js-url string       URL of the template of the JS file use for the multi version menu.
      --no-cache                 Set to 'true' to disable the Docker build cache.
  -o, --owner string             Repository owner. [required]
  -r, --repo-name string         Repository name. [required]
      --rqts-url string          Use this requirements.txt to merge with the current requirements.txt. Can be a file path.
      --version                  version for structor
```

The environment variable `STRUCTOR_LATEST_TAG` allow to override the latest tag name obtains from GitHub.

The [sprig](http://masterminds.github.io/sprig/) functions for Go templates can be used inside the JS template file.

## Download / CI Integration

```bash
curl -sfL https://raw.githubusercontent.com/containous/structor/master/godownloader.sh | bash -s -- -b $GOPATH/bin v1.7.0
```

<!--
To generate the script:

```bash
godownloader --repo=containous/structor -o godownloader.sh

# or

godownloader --repo=containous/structor > godownloader.sh
```
-->

## Examples

A simple example is available on the repository https://github.com/mmatur/structor-sample.

With menu template URL:

```shell
sudo ./structor -o containous -r traefik \
--dockerfile-url="https://raw.githubusercontent.com/containous/traefik/master/docs.Dockerfile" \
--menu.js-url="https://raw.githubusercontent.com/containous/structor/master/traefik-menu.js.gotmpl" \
--exp-branch=master --debug
```

With local menu template file:

```shell
sudo ./structor -o containous -r traefik \
--dockerfile-url="https://raw.githubusercontent.com/containous/traefik/master/docs.Dockerfile" \
--menu.js-file="~/go/src/github.com/containous/structor/traefik-menu.js.gotmpl" \
--exp-branch=master --debug
```

## What does Messor Structor mean? 

![Messor Structor](http://www.antwiki.org/wiki/images/8/8d/Messor_structor_antweb1008070_h_1_high.jpg)
