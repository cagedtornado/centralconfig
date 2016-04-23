# centralconfig [![Circle CI](https://circleci.com/gh/danesparza/centralconfig.svg?style=svg)](https://circleci.com/gh/danesparza/centralconfig)
A simple REST based service for managing application configuration using a SQL back-end.  Runs on Linux/Windows/OSX/FreeBSD/Raspberry Pi.

Back-ends supported:
- [BoltDB](https://github.com/boltdb/bolt) (default)
- [MySQL](https://www.mysql.com/)

### Quick start
To get up and running, [grab the latest release](https://github.com/danesparza/centralconfig/releases/latest) for your platform

Start the server:
```
centralconfig serve
```
Then visit the url http://localhost:3000 and you can add/edit configuration through the built-in web interface.

To customize the config, first generate a default config file (with the name centralconfig.yaml):
```
centralconfig defaults > centralconfig.yaml
```
