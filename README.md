# Resistor

Resistor is a complement to the InfluxData Kapactor tool https://github.com/influxdata/kapacitor  and has 3 functional components.

* Alert filtering system: it acts as alert filter for diferent WebHooks , it can filter by ALERTID's, time and tags, without need to change tasks variables or template definition. It can exclude alerts only on some devices or a group of them based on tags.

* Easy alert management: it can  deploy alerts based on basic templates.

* It has and resitor_udf with habilty to inject some tags / fields over datapoints depending on the value for another tag ( by example the deviceid)


If you wish to compile from source code you can follow the next steps

## Run from master
If you want to build a package yourself, or contribute. Here is a guide for how to do that.

### Dependencies

- Go 1.5
- NodeJS >=6.2.1

### Get Code

```bash
go get github.com/toni-moreno/resistor
```

### Building the backend


```bash
cd $GOPATH/src/github.com/toni-moreno/resistor
go run build.go setup            (only needed once to install godep)
godep restore                    (will pull down all golang lib dependencies in your current GOPATH)
```

### Building frontend and backend in production mode

```bash
npm install
PATH=$(npm bin):$PATH
npm run build:pro #will build fronted and backend
```
### Creating minimal package tar.gz

```bash
npm run postbuild #will build fronted and backend
```

### Creating rpm and deb packages
You  will need previously installed the fpm/rpm and deb packaging tools

```bash
go run build.go latest
```

### Recompile backend on source change (only for developers)

To rebuild on source change (requires that you executed godep restore)
```bash
go get github.com/Unknwon/bra
npm start
```
will init a change autodetect webserver with angular-cli (ng serve) and also a autodetect and recompile process with bra for the backend

## Running first time
To execute without any configuration you need a minimal config.toml file on the conf directory.

### main resistor tool

```bash
cp conf/sample.resistor.toml conf/resistor.toml
./bin/resistor
```

### resistor Injector
```bash
cp conf/sample.resinjector.toml conf/resinjector.toml
./bin/resinjector
```

## First Setup Guide 

Now you wil be able to configure metrics/measuremnets and devices from the builting web server at  http://localhost:8090 or http://localhost:4200 if working in development mode (npm start)


Resistor depends on InfluxDB and Kapacitor Tools. Be sure you have InfluxDB and Kapacitor Installed.
When first executed , you will want.

1. Add all your InfluxDB instances.
2. Add all your Kapacitor instances.
3. Build a product/measurment/fields/tags catalog.(with these 2 steps)
* import all your influx catalog first
* Organize measurements by product then.
4. Import The base Templates.

When done you will be ready to create basic alerts for all your products.

### Template naming standars

(pending)