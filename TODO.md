https://raw.githubusercontent.com/containous/traefik/master/docs.Dockerfile


sudo ./structor -o containous -r traefik \
--dockerfile-url="https://raw.githubusercontent.com/containous/traefik/master/docs.Dockerfile" \
--menu.js-file="/home/ldez/sources/go/src/github.com/containous/structor/traefik-menu.js.gotmpl" \
--exp-branch=master --debug
