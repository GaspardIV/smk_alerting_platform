
# Project structure:
    main.go - used for test local env initialization.
    config_propagator/ - script used to init database from config.
    confirmation_handler/ - confirmation handler cloud function implementation.
    distributor/ - distributor cloud function implementation.
    fake_service/ - service with ability of switching it state off/on used for functional tests.
    notifier/ - notifier  cloud function implementation.
    resolved_handler/ -  cloud function implementation.
    scheduler/ -  cloud function implementation.
    site_checker/ - site checker cloud function implementation.
    pkg/ - common files.
        cloud_tasks_queue.go - Local tasks queue interface implementation 
        consts.go - Constants.
        database.go - Database interface.
        firestore_database.go - Firestore database interface implementation.
        http_client.go - HttpClient interface and FakeHttpClient implementation.
        local_database.go - Local database interface implementation.
        local_tasks_queue.go - Local tasks queue interface implementation.
        tasks_queue.go - Tasks queue interface.

# Database:

Main collection sites consist of documents that represents single site.

sites/siteID/

each site is a document and consist of fields:
```
{
  time_until_reporting_seconds: 10 // time of inactivity after which the primary administrator is notified
  allowed_response_time_seconds: 10 // time of after which the secondary administrator is notified
  confirmation_hash: "" // hash used in url confirmation link sent to administrator 
  frequency_seconds: 2 // frequency of checking
  primary_administrator_email: "ma.brzezinski@student.uw.edu.pl" // primary administrator mail
  resolved_hash: "" // hash used in url resolved link sent to administrator
  secondary_administrator_email: "s.olearczuk@student.uw.edu.pl" // secondary administrator
  state: Running // int that represent state. State is one of four:     Running, Unavailable, Notified, Confirmed
  state_change_timestamp: 10 stycznia 2022 01:07:30 UTC+1 // timestamp of changing state
  last_change_timestamp: 10 stycznia 2022 01:07:31 UTC+1  // server timestamp of last document update - scheduler does use it. 
  url: "google.com" // page address
}
```

# Running locally
```
# install dependencies
go get .

# run locally
./run_local.sh
```

# Running tests
In order to run tests simply run `run_tests.sh` script. <br/>



```
