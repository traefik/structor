# Messor Structor: Manage multiple documentation versions with Mkdocs.

[![GitHub release](https://img.shields.io/github/release/containous/structor.svg)](https://github.com/containous/structor/releases/latest)
[![Build Status](https://travis-ci.org/containous/structor.svg?branch=master)](https://travis-ci.org/containous/structor)

Structor use git branches to create the versions of a documentation, only works with Mkdocs.

To use Structor a project must respect [semver](https://semver.org) and creates a git branch for each `MINOR` and `MAJOR` version.

Used by [Traefik](https://github.com/containous/traefik): https://docs.traefik.io

## Description

```yaml
Messor Structor: Manage multiple documentation versions with Mkdocs.

Usage: structor [--flag=flag_argument] [-f[flag_argument]] ...     set flag_argument to flag(s)
   or: structor [--flag[=true|false| ]] [-f[true|false| ]] ...     set true/false to boolean flag(s)

Available Commands:
        version                                            Display the version.
Use "structor [command] --help" for more information about a command.

Flags:
    --debug           Debug mode.                                                                    (default "false")
    --dockerfile-name Search and use this Dockerfile in the repository (in './docs/' or in './') for (default "docs.Dockerfile")
                      building documentation.                                                        
-d, --dockerfile-url  Use this Dockerfile when --dockerfile-name is not found. Can be a file path. [required]                                                                     
    --exp-branch      Build a branch as experimental.                                                
    --force-edit-url  Add a dedicated edition URL for each version.                                  (default "false")
    --image-name      Docker image name.                                                             (default "doc-site")
    --menu            Menu templates files.                                                          (default "false")
    --menu.css-file   File path of the template of the CSS file use for the multi version menu.      
    --menu.css-url    URL of the template of the CSS file use for the multi version menu.            
    --menu.js-file    File path of the template of the JS file use for the multi version menu.       
    --menu.js-url     URL of the template of the JS file use for the multi version menu.             
    --no-cache        Set to 'true' to disable the Docker build cache.                               (default "false")
-o, --owner           Repository owner. [required]                                                   
-r, --repo-name       Repository name. [required]                                                    
    --rqts-url        Use this requirements.txt to merge with the current requirements.txt. Can be a file path.
-h, --help            Print Help (this message) and exit
```

The environment variable `STRUCTOR_LATEST_TAG` allow to override the real latest tag name.

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
