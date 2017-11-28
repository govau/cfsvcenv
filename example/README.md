# Example

This example illustrates how a credential named `endpoint` from the `hackernews`
user-provided service is bound to an application environment variable named
`HACKERNEWS_ENDPOINT`.

Note: you can build the example without the `cloudfoundry` build tag and provide
the environment variable in the usual way. This illustrates that the example can
run perfectly without Cloud Foundry.

## Deployment

### Build

Use the `./build.sh` script. Take note of the `cloudfoundry` build tag.

### Deploy

To deploy this example application, use some commands like this:

```sh
#!/usr/bin/env sh

set -eu

# This will be the CF app name.
NAME="cfsvcenv-example"

# This will create or update the CF user-provided service.
SERVICE_NAME="hackernews"
cf create-user-provided-service $SERVICE_NAME -p "./$SERVICE_NAME.json" || true
cf update-user-provided-service $SERVICE_NAME -p "./$SERVICE_NAME.json"

cf push $NAME -f manifest.yml
```
