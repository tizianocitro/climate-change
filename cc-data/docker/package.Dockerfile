FROM mattermost/mattermost-enterprise-edition:7.8.0
WORKDIR /mattermost
COPY --chown=mattermost:mattermost docker/package/climate-change-data-+.tar.gz ./prepackaged_plugins/climate-change-data-+.tar.gz
