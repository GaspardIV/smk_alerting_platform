## Notifier
- if no-one was notified(State == Unavailable), then notify primary administrator: 
  - generate confirmation hash 
  - generate resolved hash
  - send mail with links via sendgrid api
  - set state to notified

- if primary administrator was already notified (state == Notified), then notify secondary administrator:
  - get resolved hash from db
  - send mail with resolved link via sendgrid api


## API
Notifier can be triggered with POST requestt