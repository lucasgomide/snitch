# Snitch
[![Build Status](https://travis-ci.org/lucasgomide/snitch.svg?branch=master)](https://travis-ci.org/lucasgomide/snitch)
[![Coverage Status](https://coveralls.io/repos/github/lucasgomide/snitch/badge.svg?branch=master)](https://coveralls.io/github/lucasgomide/snitch?branch=master)

Keep updated about each deploy via [Tsuru](https://docs.tsuru.io/stable/).

This program will notify your team and many tools when someone has deployed any application via [Tsuru](https://docs.tsuru.io/stable/).

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

So now, You must add this code into your file (`tsuru.yaml` or `tsuru.yml`) located in the root of the application [at hook's section](https://docs.tsuru.io/stable/using/tsuru.yaml.html).

```yaml
hooks:
  build:
    - curl -sSL https://github.com/lucasgomide/snitch/releases/download/0.1.0/snitch_0.1.0_linux_amd64.tar.gz | tar xz
    - ./snitch_linux/snitch -c path/snitch_config.yml
```

> Put the hook's configuration file path as argument for the *-c* option.

## Options

**-c**

This option indicates where's the file with the hook's configurations [see more](#hooks-configurations).

**-app-name-contains**

Use it to validate if the snitch should be run. If you tsuru app name does not match it the value of `-app-name-contains`, the program will stop, and no errors will be raised.

## Hook's Configurations

Here is all avaliables hook's configurations and your descriptions. Remember that you may use environment variables to define the options's values.

- Slack
  - **webhook_url** Indicates the Webhook URL to dispatch messages to Slack.

- Sentry
  - **host** Tell to snith your sentry host (e.g http://sentry.io or http://sentry.self.hosted)
  - **organization_slug** The organization slug is a unique ID used to identify your organization. (You'll find it at your sentry's configuration, probably)
  - **project_slug** The Project Slug is a unique ID used to identify your project (You'll find it at your project config)
  - **auth_token** The Auth Token to use the Sentry Web API. You can find [here](https://docs.sentry.io/api/auth/#auth-tokens)
  - **env** The application's environment variable (e.g development, production)

## Example

[Snitch App Sample](https://github.com/lucasgomide/snitch-app-example)
