package pkg

import (
        "testing"

        "github.com/stretchr/testify/assert"
)

var (
        site1 = SiteInfo{
                ID:                          "1",
                Url:                         "site1",
                PrimaryAdministratorEmail:   "primary_admin1",
                SecondaryAdministratorEmail: "secondary_admin1",
                State:                       Running,
                ConfirmationHash:            "confirmation_hash1",
                ResolvedHash:                "resolved_hash1",
                Frequency:                   1,
                TimeUntilReporting:          1,
                AllowedResponseTime:         1,
        }
        site2 = SiteInfo{
                ID:                          "2",
                Url:                         "site2",
                PrimaryAdministratorEmail:   "primary_admin2",
                SecondaryAdministratorEmail: "secondary_admin2",
                State:                       Unavailable,
                ConfirmationHash:            "confirmation_hash2",
                ResolvedHash:                "resolved_hash2",
                Frequency:                   2,
                TimeUntilReporting:          2,
                AllowedResponseTime:         2,
        }
)

func TestLocalDatabase_GetSiteInfo_Success(t *testing.T) {
        var database = CreateLocalDatabase([]SiteInfo{site1, site2})

        cases := []struct {
                name         string
                field        string
                value        interface{}
                expectedInfo SiteInfo
        }{
                {
                        "site1 found by url",
                        "url",
                        "site1",
                        site1,
                },
                {
                        "site2 found by primary administrator",
                        "primary_administrator_email",
                        "primary_admin2",
                        site2,
                },
                {
                        "site1 found by secondary administrator",
                        "secondary_administrator_email",
                        "secondary_admin1",
                        site1,
                },
                {
                        "site2 found by state",
                        "state",
                        Unavailable,
                        site2,
                },
                {
                        "site1 found by confirmation hash",
                        "confirmation_hash",
                        "confirmation_hash1",
                        site1,
                },
                {
                        "site2 found by resolved hash",
                        "resolved_hash",
                        "resolved_hash2",
                        site2,
                },
                {
                        "site1 found by frequency",
                        "frequency_seconds",
                        1,
                        site1,
                },
                {
                        "site2 found by alerting_window",
                        "time_until_reporting_seconds",
                        2,
                        site2,
                },
                {
                        "site1 found by allowed_response_time",
                        "allowed_response_time_seconds",
                        1,
                        site1,
                },
        }
        for _, c := range cases {
                t.Run(c.name, func(t *testing.T) {
                        info, err := database.GetSiteInfo(c.field, c.value)
                        assert.Equal(t, c.expectedInfo, info)
                        assert.Nil(t, err)
                })
        }
}

func TestLocalDatabase_GetSiteInfo_SiteNotFoundError(t *testing.T) {
        var database = CreateLocalDatabase([]SiteInfo{site1})

        cases := []struct {
                name  string
                field string
                value interface{}
        }{
                {
                        "site not found by url",
                        "url",
                        "incorrect_value",
                },
                {
                        "site not found by primary administrator",
                        "primary_administrator_email",
                        "incorrect_value",
                },
                {
                        "site not found by secondary administrator",
                        "secondary_administrator_email",
                        "incorrect_value",
                },
                {
                        "site not found by state",
                        "state",
                        "incorrect_value",
                },
                {
                        "site not found by confirmation hash",
                        "confirmation_hash",
                        "incorrect_value",
                },
                {
                        "site not found by resolved hash",
                        "resolved_hash",
                        "incorrect_value",
                },
                {
                        "site not found by frequency",
                        "frequency_seconds",
                        "incorrect_value",
                },
                {
                        "site not found by alerting window",
                        "time_until_reporting_seconds",
                        "incorrect_value",
                },
                {
                        "site not found by allowed response time",
                        "allowed_response_time_seconds",
                        "incorrect_value",
                },
        }
        for _, c := range cases {
                t.Run(c.name, func(t *testing.T) {
                        info, err := database.GetSiteInfo(c.field, c.value)
                        assert.Equal(t, SiteInfo{}, info)
                        assert.Equal(t, siteNotFound, err)
                })
        }
}

func TestLocalDatabase_UpdateSite_Success(t *testing.T) {
        var database = CreateLocalDatabase([]SiteInfo{site1})
        newSiteInfo := SiteInfo{ID: site1.ID, Url: site1.Url}

        err := database.UpdateSite(newSiteInfo)
        assert.Nil(t, err)

        site, err := database.GetSiteInfo("url", newSiteInfo.Url)
        assert.Nil(t, err)
        // LastChangeTimestamp gets overwritten in UpdateSite
        newSiteInfo.LastChangeTimestamp = site.LastChangeTimestamp
        assert.Equal(t, newSiteInfo, site)
}

func TestLocalDatabase_UpdateSite_SiteNotFoundError(t *testing.T) {
        var database = CreateLocalDatabase([]SiteInfo{site1})
        newSiteInfo := SiteInfo{ID: "some_id", Url: site1.Url}

        err := database.UpdateSite(newSiteInfo)
        assert.Equal(t, siteNotFound, err)

        site, err := database.GetSiteInfo("url", newSiteInfo.Url)
        assert.Equal(t, site1, site)
        assert.Nil(t, err)
}