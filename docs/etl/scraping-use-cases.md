## Wiz scraping use cases


```yaml
data:
    url: http://cnn.com/home or article # (required)
    # handle file writing (e.g. what should the file be named, just the document name?)
    method: GET # | POST | other HTTP verbs
    body: |
        http body content for POST requests, etc
    headers: # todo: add default headers
        Accept: text/html
```

Outline: initial homepage or article, HTTP requester -> web parser, find links -> back to http requester

In this scenario, the task graph is not a DAG because each child node links back to the parent with more data.

This could be done with queue or stream. research how distributed scrapers work

For now, the scraper could all be in one package. include requests and parsing together

## data splitting

```yaml
initial_data:
    base_url: https://download.bls.gov/pub/time.series/la
    files: 
    url:
        - https://download.bls.gov/pub/time.series/la/la.period
        - https://download.bls.gov/pub/time.series/la/la.series
        - https://download.bls.gov/pub/time.series/la/la.area_type
```

could be generic as in

```yaml
pipeline: 
    name: splitter
    data: 
        url:
            - https://download.bls.gov/pub/time.series/la/la.period
            - https://download.bls.gov/pub/time.series/la/la.series
            - https://download.bls.gov/pub/time.series/la/la.area_type
    template:
        url: $url
        method: GET
    loop: data.url
```