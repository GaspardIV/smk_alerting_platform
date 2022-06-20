## Site Checker
- checks availability of provided sites (in separates goroutines)
- it changes state of given site based on its availability
- if site has been unavailable for long enough, then it triggers notifier by inserting two tasks into TaskQueue:
  - one that triggers instantly notifier in order to notify primary administrator.
  - second one that triggers notifier after Allowed Response Time in order to notify secondary administrator, if primary didn't confirmed.

## API
Site checker can be triggered with POST requests, sites are provided in body e.g.
```
{
        "urls": ["https://google.com", "Asdfasdfas"]
}
```