# Messor Structor: Manage multiple documentation versions with Mkdocs.

[![GitHub release](https://img.shields.io/github/release/containous/structor.svg)](https://github.com/containous/structor/releases/latest)
[![Build Status](https://travis-ci.org/containous/structor.svg?branch=master)](https://travis-ci.org/containous/structor)

Structor use git branches to create the versions of a documentation, only works with Mkdocs.

To use Structor a project must respect [semver](https://semver.org) and creates a branch for each `MINOR` and `MAJOR` version.

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
    --debug          Debug mode.                                                               (default "false")
-d, --dockerfile-url Dockerfile URL. [required]                                                
    --exp-branch     Build a branch as experimental.                                           
    --image-name     Docker image name.                                                        (default "doc-site")
    --menu           Menu templates files.                                                     (default "false")
    --menu.css-file  File path of the template of the CSS file use for the multi version menu. 
    --menu.css-url   URL of the template of the CSS file use for the multi version menu.       
    --menu.js-file   File path of the template of the JS file use for the multi version menu.  
    --menu.js-url    URL of the template of the JS file use for the multi version menu.        
-o, --owner          Repository owner. [required]                                              
-r, --repo-name      Repository name. [required]                                               
    --rqts-url       URL of the requirements.txt file.                                     
-h, --help           Print Help (this message) and exit                                     
```

The environment variable `STRUCTOR_LATEST_TAG` allow to override the real latest tag name.

## Download / CI Integration

```bash
curl -sfL https://raw.githubusercontent.com/containous/structor/master/godownloader.sh | bash -s -- -b $GOPATH/bin v1.4.3
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

## What's Messor Structor mean? 

![Messor Structor](http://www.antwiki.org/wiki/images/8/8d/Messor_structor_antweb1008070_h_1_high.jpg)
