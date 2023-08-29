# climate-change
A collaboration platform based on Mattermost and Hyperlinking Object-Oriented Discussion to identify mis/disinformation in climate change discussions.

# Components
1. [cc-data]() which enables the object-oriented collaboration mechanism with support for the hyper-linking system.
1. [cc-data-provider]() a web server that provides fake data using the RESTful protocol.

# Installing

- Build the previous package using the instructions of each project.
- Execute the command: `sudo ./start.sh` to clean the compose and run mattermost and cs-connect with the faker data provider.
