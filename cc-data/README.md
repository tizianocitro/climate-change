# Climate Change Data
A plugin to enable HOOD within Mattermost.

## Installation
Build the Docker image for the plugin build environment using the following command.

It is possible to change the base container name, but it will require updating the start script.

```sh
docker build -t cs-connect-base -f docker/dev.Dockerfile .
```

Launch the `build.sh` script to build the custom Mattermost Docker image with the plugin installed.
