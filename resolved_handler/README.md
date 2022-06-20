## Resolved Handler
Cloud function triggered by GET request with ?hash=aaaaa param. 
- Searches in database for Site with equivalent ResolvedHash.
- Changes its state back to Running.
- Returns readable response about process successfulness.

## API
Resolved Handler can be triggered with GET requests, hash is provided in query parameter named `hash`