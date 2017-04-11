# Snitch
[![Build Status](https://travis-ci.org/lucasgomide/snitch.svg?branch=master)](https://travis-ci.org/lucasgomide/snitch)
[![Coverage Status](https://coveralls.io/repos/github/lucasgomide/snitch/badge.svg?branch=master)](https://coveralls.io/github/lucasgomide/snitch?branch=master)

Keep updated about each deploy via [Tsuru](https://docs.tsuru.io/stable/)

This program will notify your team when someone has deployed any application via [Tsuru](https://docs.tsuru.io/stable/).

## Quick Start

You must add this code into your file (`tsuru.yaml` or `tsuru.yml`) located in the root of the application [at hook's section](https://docs.tsuru.io/stable/using/tsuru.yaml.html)


### From binary

```yaml
hooks:
  build:
    - curl -sSL https://github.com/lucasgomide/snitch/releases/download/0.1.0/snitch_0.1.0_linux_amd64.tar.gz | tar xz
    - ./snitch_linux/snitch -h hookName -url webHookURL
```

## Hooks availabe

[Slack](hook/slack.go)

## Example

[Snitch App Sample](https://github.com/lucasgomide/snitch-app-example)
