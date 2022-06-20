module config_propagator

go 1.13

require (
        smk_alerting_platform/pkg v0.0.0-00010101000000-000000000000
)

replace (
        smk_alerting_platform/pkg => ../pkg
)