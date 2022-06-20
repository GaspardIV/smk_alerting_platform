module smk_alerting_platform

go 1.13

require (
        cloud.google.com/go/secretmanager v1.0.0 // indirect
        github.com/gin-gonic/gin v1.7.7
        github.com/robfig/cron/v3 v3.0.0
        github.com/stretchr/testify v1.6.1 // indirect
        google.golang.org/api v0.63.0 // indirect
        smk_alerting_platform/confirmation_handler v0.0.0-00010101000000-000000000000
        smk_alerting_platform/distributor v0.0.0-00010101000000-000000000000
        smk_alerting_platform/notifier v0.0.0-00010101000000-000000000000
        smk_alerting_platform/pkg v0.0.0-00010101000000-000000000000
        smk_alerting_platform/resolved_handler v0.0.0-00010101000000-000000000000
        smk_alerting_platform/scheduler v0.0.0-00010101000000-000000000000
        smk_alerting_platform/site_checker v0.0.0-00010101000000-000000000000
)

replace (
        smk_alerting_platform/confirmation_handler => ./confirmation_handler
        smk_alerting_platform/distributor => ./distributor
        smk_alerting_platform/notifier => ./notifier
        smk_alerting_platform/pkg => ./pkg
        smk_alerting_platform/resolved_handler => ./resolved_handler
        smk_alerting_platform/scheduler => ./scheduler
        smk_alerting_platform/site_checker => ./site_checker
)