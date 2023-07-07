FROM cs-connect-base AS BUILDER
WORKDIR home/mattermost/cc-data
COPY . .
ARG CONFIG_FILE_NAME=config.yml
RUN make CONFIG_FILE_NAME=$CONFIG_FILE_NAME

FROM mattermost/mattermost-enterprise-edition:7.8.0
WORKDIR /mattermost
COPY --from=BUILDER /home/mattermost/cc-data/dist/climate-change-data-+.tar.gz ./prepackaged_plugins/climate-change-data-+.tar.gz
