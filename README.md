# Structured Log

Package provides to structured log both for displaying it sanely in the stderr
logs as well as writing key-valued logs into ElasticSearch:

Using this package it's possible to use it with Goa as well and logs will be
displayed in the stderr like this:

```
2017-09-21 18:58:19 [DEBUG] {goa} completed
                            ├─ req_id: aIlGSFNnk8-1
                            ├─ status: 200
                            ├─ bytes: 134030
                            ├─ time: 1.288885343s
                            ├─ ctrl: MetricsController
                            └─ action: get
```
