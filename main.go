package main

import (
        "smk_alerting_platform/confirmation_handler"
        "smk_alerting_platform/distributor"
        "smk_alerting_platform/notifier"
        "smk_alerting_platform/pkg"
        "smk_alerting_platform/resolved_handler"
        "smk_alerting_platform/scheduler"
        "smk_alerting_platform/site_checker"

        "context"
        "log"
        "os"
        "time"

        "github.com/gin-gonic/gin"
        "github.com/robfig/cron/v3"
        "net/http"
)

// TODO: Consider getting this from config or something
var google = pkg.SiteInfo{
        ID:                          "1",
        PrimaryAdministratorEmail:   "brzezinskimarcin97@gmail.com",
        SecondaryAdministratorEmail: "ma.brzezinski@student.uw.edu.pl",
        Url:                         "https://google.com",
        Frequency:                   10,
        TimeUntilReporting:          15,
        AllowedResponseTime:         60,
        State:                       pkg.Unavailable,
        LastChangeTimestamp:         time.Now(),
}
var asdf = pkg.SiteInfo{
        ID:                          "2",
        PrimaryAdministratorEmail:   "brzezinskimarcin97@gmail.com",
        SecondaryAdministratorEmail: "ma.brzezinski@student.uw.edu.pl",
        Url:                         "Asdfasdfas",
        Frequency:                   5,
        TimeUntilReporting:          15,
        AllowedResponseTime:         10,
        State:                       pkg.Running,
        LastChangeTimestamp:         time.Now(),
}
var database = pkg.CreateLocalDatabase([]pkg.SiteInfo{google, asdf})
var queue = pkg.CreateLocalTaskQueue()
var schedulerInterval = "1m"

func notifierHandler(c *gin.Context) {
        notifier.NotifierWithDatabase(c.Writer, c.Request, database)
}

func distributorHandler(c *gin.Context) {
        distributor.Distributor(c.Writer, c.Request)
}

func siteCheckerHandler(c *gin.Context) {
        site_checker.SiteCheckerWithDatabaseHTTPClientQueue(c.Writer, c.Request, database, &http.Client{}, queue)
}

func resolvedHandler(c *gin.Context) {
        resolved_handler.ResolvedHandlerWithDatabase(c.Writer, c.Request, database)
}

func confirmationHandler(c *gin.Context) {
        confirmation_handler.ConfirmationHandlerWithDatabase(c.Writer, c.Request, database)
}

func schedulerHandler() {
        err := scheduler.SchedulerWithDatabaseQueue(context.Background(), scheduler.PubSubMessage{Data: []byte(schedulerInterval)}, database, queue)
        if err != nil {
                log.Fatalf("%v ", err)
        }
}

func main() {
        os.Setenv("FUNCTION_BASE_URL", "http://localhost:8080")
        os.Setenv("URLS_PER_CHECKER", "1")
        os.Setenv("TURN_ON_EMAIL_SENDING", "false")
        os.Setenv("ROUNDING_FACTOR", "1s")

        c := cron.New()
        c.AddFunc("* * * * *", schedulerHandler)
        c.Start()

        router := gin.Default()
        router.POST("/site-checker", siteCheckerHandler)
        router.POST("/distributor", distributorHandler)
        router.POST("/notifier", notifierHandler)
        router.GET("/resolved-handler", resolvedHandler)
        router.GET("/confirmation-handler", confirmationHandler)
        router.Run("localhost:8080")
}