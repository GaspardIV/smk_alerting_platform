module smk_alerting_platform/notifier

go 1.13

replace smk_alerting_platform/pkg => ../pkg

require (
        cloud.google.com/go/secretmanager v1.0.0
        github.com/fatih/structs v1.1.0 // indirect
        github.com/sendgrid/rest v2.6.6+incompatible // indirect
        github.com/sendgrid/sendgrid-go v3.10.4+incompatible
        google.golang.org/genproto v0.0.0-20211208223120-3a66f561d7aa
        smk_alerting_platform/pkg v0.0.0-00010101000000-000000000000
)