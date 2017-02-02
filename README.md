# mackerel-plugin-gcp-billing

GCP Billing custom metrics plugin for mackerel.io agent.

## Requirement

The billing data is in BigQuery. And the following schema is required.

|Name|Type|
|:--|:--|
|cost|FLOAT|
|start_time|TIMESTAMP|

- Recommend to use [export billing data to BigQuery](https://support.google.com/cloud/answer/7233314?hl=en&ref_topic=7106112).

And [Application Default Credentials(JSON key)](https://developers.google.com/identity/protocols/application-default-credentials) is required.

## Example of mackerel-agent.conf

```
[plugin.metrics.gcp-billing]
command = "env GOOGLE_APPLICATION_CREDENTIALS=/path/to/JSONKEY.json /path/to/mackerel-plugin-gcp-billing -d DATASET -p PROJECTID -t TABLE"
```

## Options

```bash
Application Options:
  -p, --projectid= Project ID (require)
  -d, --dataset=   Dataset    (require)
  -t, --table=     Table      (require)

Help Options:
  -h, --help       Show this help message
```