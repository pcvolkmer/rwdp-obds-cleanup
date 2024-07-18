# A simple pipeline step to clean up incoming kafka messages

This can be used to alter kafka messages for [onco-analytics-on-fhir](https://github.com/bzkf/onco-analytics-on-fhir)
before any other processing of data.

## Current changes made to Kafka record values

* This images removes leading zeros from `Patient_ID` in oBDS message.

## Run docker container

See [docker-compose.yml](docker-compose/docker-compose.yml) for an example on how to use the image, available 
[here](https://github.com/pcvolkmer/rwdp-obds-cleanup/pkgs/container/rwdp-obds-cleanup).

Change `BOOTSTRAP_SERVERS`, `INPUT_TOPICS` or `OUTPUT_TOPIC` to apply any configuration changes.