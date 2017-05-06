# Snitch
[![Build Status](https://travis-ci.org/lucasgomide/snitch.svg?branch=master)](https://travis-ci.org/lucasgomide/snitch)
[![Coverage Status](https://coveralls.io/repos/github/lucasgomide/snitch/badge.svg?branch=master)](https://coveralls.io/github/lucasgomide/snitch?branch=master)

Keep updated about each deploy via [Tsuru](https://docs.tsuru.io/stable/)

This program will notify your team when someone has deployed any application via [Tsuru](https://docs.tsuru.io/stable/).

## Quick Start

First one, you have to create a hook's configuration file. This file describe wich hook will be dispatched and the your configurations (e.g webhook_url).

You can add this code into your file, hardcode mode:
```yaml
slack:
  webhook_url: http://your.webhook.here
```

or using environment variable:
```yaml
slack:
  webhook_url: $SLACK_WEBHOOK_URL
```

So now, You must add this code into your file (`tsuru.yaml` or `tsuru.yml`) located in the root of the application [at hook's section](https://docs.tsuru.io/stable/using/tsuru.yaml.html)

```yaml
hooks:
  build:
    - curl -sSL https://github.com/lucasgomide/snitch/releases/download/0.1.0/snitch_0.1.0_linux_amd64.tar.gz | tar xz
    - ./snitch_linux/snitch -c path/snitch_config.yml
```

> Put the hook's configuration file path as argument for the *-c* option.

## Options

**-c**

This option indicates where's the file with the hook's configurations [see more](#hooks-configurations)

**-app-name-contains**

Use it to validate if the snitch should be run. If you tsuru app name does not match it the value of `-app-name-contains`, the program will stop, and no errors will be raised.

## Hook's Configurations

Here is all avaliables hook's configurations and your descriptions.

- Slack
  - webhook_url: Indicates the Webhook URL to dispatch messages to Slack.

## Example

[Snitch App Sample](https://github.com/lucasgomide/snitch-app-example)
