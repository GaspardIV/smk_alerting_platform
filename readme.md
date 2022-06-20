# DS Course Project: Alerting Platform 
Final project of Distributed Systems course done together with [@olearczuk](https://github.com/olearczuk) and [@brzezinskimarcin](https://github.com/brzezinskimarcin).

# Design doc
### Used technologies
 - Golang
 - Firestore for persistent and highly available storage
 - Cloud Functions for executing monitoring steps
 - Cloud Scheduler, Cloud Tasks and Cloud Tasks Queue for repeating calling of cloud functions
 - Kubernetes for deploying fake service to test
 - Cloud Build for continuous deployment

## Overview of the system
![Screenshot 2022-06-20 at 19 59 57](https://user-images.githubusercontent.com/30477366/174657679-535f5e68-5d48-4086-bf75-bb53ec4c4a90.png)

## Components overview
 - Cloud scheduler - a cron job that starts the whole checking process periodically (every minute or so)
 - Scheduler Function - gets awaken by Cloud Scheduler (via HTTP), retrieves sites eligible for monitoring from database and spawns tasks
 - Cloud Tasks - trigger Distributor Function with specified services to check
 - Distributor Function - splits specified sites into chunks (we don’t want to pass e.g. 100 services to a single cloud function) and triggers Site Checker Functions
 - Site Checker Function - performs a single check of given set of sites, updates database and triggers Notifier Function if site has not been available for long enough
 - Notifier Function - sends email to Service Administrator with links to Confirmation Handler and Resolve Handler (secondary administrator gets email only to Resolve Handler) and updates database
 - Confirmation Handler Function -  handles confirmation of primary administrator (this way he confirms that he started working on the issue and there is no need to notify secondary administrator) and updates the database
 - Resolved Handler Function - handles confirmation that problem has been resolved by an administrator, this means that the site is eligible for being monitored again

## Database scheme
![Screenshot 2022-06-20 at 19 58 27](https://user-images.githubusercontent.com/30477366/174657670-49401195-50d0-4980-a794-b90573ea17e8.png)

Possible states are:
 - Running - site is (or should be according to administrator that resolves the issue) working correctly
 - Unavailable - SiteChecker spotted that site is unavailable
 - Notified - primary administrator has been notified
 - Confirmed - primary administrator confirmed that he is working on the issue
 
## Pricing
According to Google Cloud, first 2 million invocations of Cloud Functions are free and 0.40 USD per 2 million after that. Similarly for Cloud Tasks. We won’t consider Notifier and Confirmation/Resolved handlers since they are related to handling sites being unavailable, which produces negligible traffic. Assuming every site is checked every minute, a single site generates 3 (cloud task + distributor + checker) * 60 * 24 * 30 = 129600. That means that the first ~15 sites are free and every next ~15 sites cost us 0.4 USD monthly. This number in reality will be bigger than 15 because the distributor reduces the number of Checker calls.

According to Google Cloud first 3 jobs of Cloud Scheduler are free and we only need one.

According to Firestore pricing we will be charged for each document read, write, and delete that we perform with Firestore. Assuming each cloud function run results in execution up to 1 documents reads and 1 document write on average, then each site generates costs up to 0.324 USD + 0.18 per GiB of stored data.



## Monitoring and alerting
Any needed monitoring/metrics for all used Cloud Services (Firestore, Scheduler, Tasks, Cloud Functions) are provided by Google Cloud in gcloud console

## Auto healing
All services used (Cloud Functions, Firestore, Scheduler, Tasks) are serverless and entirely managed by Google, so auto healing is supported out of the box by all of them.

## High availability (multi-node)
 - Cloud Firestore - automatically scales up and down based on demand. It requires no maintenance, and provides high availability of >= 99.999% (monthly uptime percentage) achieved through strongly consistent data replication.
 - Cloud Scheduler 
 - Cloud Tasks - fully managed by Google, and have high availability of >= 99.95%, guaranteed by Google.
 - Cloud Functions - on zone-level are entirely managed by Google, so the high availability >= 99.95% (monthly uptime percentage) is guaranteed across all the zones.

## High availability (multi-region)
 - Firestore - we will use multi-regional configuration, so thanks to the automatic data replication across regions Firestore database will be available in different regions.
 - Cloud Scheduler 
 - Cloud Tasks are regional, which means the infrastructure managing the queue is located in a specific region, so using it from different region could result in some latency or availability problems. We could solve this issue by having multiple queues across regions. 
 - Similarly, Cloud Functions are regional, and we could solve it in a similar way, by deploying all functions in multiple regions.

## Various concerns
 - Ideally we would like to keep as much traffic as possible within a single region, however this might be complicated in Cloud Functions approach
## Testing scenarios
 1. End-to-end tests
We are going to deploy a test service and switch it on and off. This is going to allow us to perform manual tests of the entire project (testers would mimic administrators).
Integration tests
 2. We can verify integration between SiteChecker and Notifier. It could be done by triggering SiteChecker to the moment when it should realise that a site is not available long enough. At this point it should trigger Notifier and we should receive an email (end state of this site should be set to Notified).
 3. Unit tests
For each cloud function we will implement dedicated unit tests.
 4. Stress and load tests
First way to test the system against a number of observable services. We could do it by creating a lot of sites in the database and test system behavior on a big scale.
Secondly, in order to test throughput of the system (including services provided by Google Cloud) we could adjust the frequency of sites pinging.


## Code build & Deployment
Thanks to the use of cloud functions, and firestore the deployment process is significantly simplified. In addition we would like to take advantage of Continuous Integration and Deployment (CI/CD) pipelines. As Cloud Functions are not updated automatically, we will  configure CI/CD pipelines to automatically test and redeploy our functions from Cloud Source Repositories. To achieve that, we will use Cloud Build, triggered on each push to main branch.

Deployment steps:
 - run unit tests
 - pause cloud scheduler
 - deploy cloud functions
 - propagate config to database
 - resume cloud scheduler

## Encryption of notification messages
Communication channels are encrypted out of the box, since we are using Sendgrid for sending messages.
Links are obfuscated, because we are generating salted hashes which are stored in the database so they are not predictable and can’t be easily silenced.


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
