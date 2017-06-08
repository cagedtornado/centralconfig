# centralconfig [![Circle CI](https://circleci.com/gh/cagedtornado/centralconfig.svg?style=shield)](https://circleci.com/gh/cagedtornado/centralconfig) [![stable](http://badges.github.io/stability-badges/dist/stable.svg)](http://github.com/badges/stability-badges)

A simple REST based service for managing application configuration across a cluster.  
Runs natively on: Linux, [Windows](https://www.microsoft.com/en-us/windows), [OSX](http://www.apple.com/osx/), [FreeBSD](https://www.freebsd.org/), [NetBSD](https://www.netbsd.org/), [OpenBSD](http://www.openbsd.org/), and even [Raspberry Pi](https://www.raspberrypi.org/).

Storage back-ends supported:
- [BoltDB](https://github.com/boltdb/bolt) (default)
- [MySQL](https://www.mysql.com/)
- [Microsoft SQL server (MSSQL)](https://www.microsoft.com/en-us/server-cloud/products/sql-server/)

### Quick start
To get up and running, [grab the latest release](https://github.com/danesparza/centralconfig/releases/latest) for your platform

Start the server:
```
centralconfig serve
```
Then visit the url [http://localhost:3000](http://localhost:3000) and you can add/edit your configuration through the built-in web interface.  

If no other configuration is specified, BoltDB will be used to store your config items in a file called 'config.db' in the working directory.

### Docker quick start
To use the centralconfig docker image: 

[Install Docker](https://docs.docker.com/mac/started/)

Start the server:
```
docker run --restart=on-failure -d -p 3000:3000 cagedtornado/centralconfig:latest
```

### Configuration
To customize the config, first generate a default config file (with the name centralconfig.yaml):
```
centralconfig defaults > centralconfig.yaml
```
