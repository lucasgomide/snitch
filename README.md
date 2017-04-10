# Snitch
Keep updated about each deploy via [Tsuru](https://docs.tsuru.io/stable/)

This program will notify your team when someone has deployed any application via [Tsuru](https://docs.tsuru.io/stable/).

## Quick Start

You must add this code into your file (`tsuru.yaml` or `tsuru.yml`) located in the root of the application [at hook's section](https://docs.tsuru.io/stable/using/tsuru.yaml.html)


### From binary

Download the latest release

``` bash
$ curl -sSL https://github.com/lucasgomide/snitch/releases/download/0.0.1/snitch-0.0.1-darwin_amd64.tar.gz \
  | tar xz
```

Put into your tsuru app config

```yaml
hooks:
  build:
    - snitch -h hookName -url webHookURL
```

## Hooks availabe

[Slack](hook/slack.go)

