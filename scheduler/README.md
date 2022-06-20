## Scheduler
Cloud function triggered by PubSub message from Cloud scheduler periodically.
- Divides the given interval by frequency and create task for each site for each scheduled site check time.
- Groups tasks by trigger time.
- Schedule groups to distributor.

## PubSub message data
`interval` - what time period he has to cover with tasks
