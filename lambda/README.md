[![Actions Status](https://github.com/tommzn/hdb-datasource-weather/actions/workflows/go.image.build.yml/badge.svg)](https://github.com/tommzn/hdb-datasource-weather/actions)
[![CircleCI](https://circleci.com/gh/tommzn/hdb-datasource-weather/tree/main.svg?style=svg)](https://circleci.com/gh/tommzn/hdb-datasource-weather/tree/main)

# Weather Data Collector
This package composes a [data collector](https://github.com/tommzn/hdb-datasource-core/collector.go) and [weather data source](https://github.com/tommzn/hdb-datasource-weather) to fetch weather data and publish it to a SQS queue.

## Config
This collector requires a config to get settings for OpenWeatherMap api, AWS SQS and maybe some logging settings.

### Example 
```yaml
log:
  loglevel: error
  shipper: logzio  

hdb:
  queue: MyWeatherDataQueue
  archive: MyEventArchiveQueue

weather:
  owm:
    url: https://api.openweathermap.org/data/2.5
    latitude: 123.12
    longitude: 123.12
    units: metric

aws:
  sqs:
    region: eu-west-1
```

# Links
- [HomeDashboard Documentation](https://github.com/tommzn/hdb-docs/wiki)

