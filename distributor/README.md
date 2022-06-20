## Distributor
Splits specified sites into chunks (we don’t want to pass e.g. 100 services to a single cloud function).
- Reads env variable "URLS_PER_CHECKER"
- get urls from 
- splits urls into Ceil(urlsCount / urlsPerChecker) parts
- and for each part triggers Site Checker Functions

## API
Distributor can be triggered with POST requests, sites are provided in request body e.g.
```
{
"urls": ["https://google.com", "Asdfasdfas"]
}
```