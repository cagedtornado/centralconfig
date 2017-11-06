# centralconfig [![Circle CI](https://circleci.com/gh/cagedtornado/centralconfig.svg?style=shield)](https://circleci.com/gh/cagedtornado/centralconfig) [![Go Report Card](https://goreportcard.com/badge/github.com/cagedtornado/centralconfig)](https://goreportcard.com/report/github.com/cagedtornado/centralconfig) [![stable](http://badges.github.io/stability-badges/dist/stable.svg)](http://github.com/badges/stability-badges)

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

### Supported environment variables
If you're using centralconfig in as part of a [12 factors app](https://12factor.net/config) environment or just want to set centralconfig service settings through environment variables, you have the following settings available:

| Command | Description |
| --- | --- |
| `SERVER.SSLCERT` | Path to the SSL certificate file |
| `SERVER.SSLKEY` | Path to the SSL certificate key |
| `DATASTORE.TYPE` | The type of backing storage for configuration.  One of: mysql, mssql, boltdb |
| `DATASTORE.ADDRESS` | Location of the backing store |
| `DATASTORE.DATABASE` | Database name to use in the backing store |
| `DATASTORE.USER` | Databse user to use |
| `DATASTORE.PASSWORD` | Database password to use |

#### Example (with docker)
```
docker run --restart=unless-stopped -d -p 3800:3000 -v /private/etc/ssl:/certs -e "SERVER.SSLCERT=/certs/sslcert.pem" -e "SERVER.SSLKEY=/certs/sslcert.key" -e "DATASTORE.TYPE=mysql" -e "DATASTORE.ADDRESS=mysqldatabaseserver:3306" -e "DATASTORE.DATABASE=centralconfig" -e "DATASTORE.USER=myusername" -e "DATASTORE.PASSWORD=thepasswordhere" cagedtornado/centralconfig:154
```


