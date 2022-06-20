module smk_alerting_platform/resolved_handler

go 1.13

replace smk_alerting_platform/pkg => ../pkg

require (
        github.com/fatih/structs v1.1.0 // indirect
        github.com/stretchr/testify v1.6.1 // indirect
        smk_alerting_platform/pkg v0.0.0-00010101000000-000000000000
)