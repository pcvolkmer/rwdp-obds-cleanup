# A simple pipeline step to clean up incoming kafka messages

This can be used to alter kafka messages for [onco-analytics-on-fhir](https://github.com/bzkf/onco-analytics-on-fhir)
before any other processing of data.

## Current modifications and actions

* Removes leading zeros from `Patient_ID` in oBDS message.
* Drops Kafka records not containing an oBDS 2.x message (e.g. new oBDS 3.x)

All features are enabled by default.

## Run docker container

See [docker-compose.yml](docker-compose/docker-compose.yml) for an example on how to use the image, available 
[here](https://github.com/pcvolkmer/rwdp-obds-cleanup/pkgs/container/rwdp-obds-cleanup).

Change `BOOTSTRAP_SERVERS`, `INPUT_TOPICS` or `OUTPUT_TOPIC` to apply any configuration changes.

### Enable or disable features

To enable or disable a feature, set the given env var to `true` or `false`. All features are enabled by default.

* `ENABLE_REMOVE_PATIENT_ID_ZERO`: Enable removing of leading patient id zeros
* `ENABLE_DROP_NON_OBDS_2`: Drop records not containing oBDS 2.x messages 