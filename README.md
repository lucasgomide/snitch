# Snitch
[![Documentation](https://godoc.org/github.com/lucasgomide/snitch?status.svg)](http://godoc.org/github.com/lucasgomide/snitch)
[![Coverage Status](https://coveralls.io/repos/github/lucasgomide/snitch/badge.svg?branch=master)](https://coveralls.io/github/lucasgomide/snitch?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/lucasgomide/snitch)](https://goreportcard.com/report/github.com/lucasgomide/snitch)
[![Build Status](https://travis-ci.org/lucasgomide/snitch.svg?branch=master)](https://travis-ci.org/lucasgomide/snitch)

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
  - **host** Tell to Snitch your sentry host (e.g http://sentry.io or http://sentry.self.hosted)
  - **organization_slug** The organization slug is a unique ID used to identify your organization. (You'll find it at your sentry's configuration, probably)
  - **project_slug** The Project Slug is a unique ID used to identify your project (You'll find it at your project config)
  - **auth_token** The token used to authenticate on Sentry API. To generate a new token, you have to access [manager auth tokens](https://sentry.io/api) then create a token. If you are using Sentry self hosted, you need change the domain _sentry.io_ to your own domain, example: _sentry.snitch.com/api_. Find more information [on Sentry documentation](https://docs.sentry.io/api/auth/#auth-tokens)
  - **env** The application's environment variable (e.g development, production)

- Rollbar
  - **access_token** The access token with `post_server_item` scope. You can find more [here](https://rollbar.com/docs/api/#authentication)
  - **env** The application's environment variable (e.g development, production)

- NewRelic
  - **host** Tell to Snitch your NewRelic API host (e.g https://api.newrelic.com)
  - **application_id** The application ID is a unique ID used to identify your application in APM. (You'll find it at the end of the application's page URL)
  - **api_key** The API Key to use the NewRelic REST API. You can find more [here](https://docs.newrelic.com/docs/apis/rest-api-v2/getting-started/api-keys)
  - **revision** The application's current revision (e.g 0.0.1r42)

- HangoutsChat
  - **webhook_url** Indicates the Webhook URL to dispatch messages to HangoutsChat Room.

## Example

[Snitch App Sample](https://github.com/lucasgomide/snitch-app-example)
