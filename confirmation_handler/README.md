## Confirmation Handler
Cloud function triggered by GET request with ?hash=aaaaa param.
- Searches in database for Site with equivalent ConfirmationHash.
- Changes its state to Confirmed.
- Returns readable response about process successfulness.

## API
Confirmation Handler can be triggered with GET requests, hash is provided in query parameter named `hash`.