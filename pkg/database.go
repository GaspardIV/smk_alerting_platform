package pkg

import (
        "testing"
        "time"

        "github.com/stretchr/testify/assert"
)

const (
        // possible site states

        Running     string = "Running"
        Unavailable string = "Unavailable"
        Notified    string = "Notified"
        Confirmed   string = "Confirmed"
)

type SiteInfo struct {
        ID                          string    `structs:"-"`
        Url                         string    `firestore:"url" structs:"url"`
        PrimaryAdministratorEmail   string    `firestore:"primary_administrator_email" structs:"primary_administrator_email"`
        SecondaryAdministratorEmail string    `firestore:"secondary_administrator_email" structs:"secondary_administrator_email"`
        LastChangeTimestamp         time.Time `firestore:"last_change_timestamp" structs:"last_change_timestamp"`
        StateChangeTimestamp        time.Time `firestore:"state_change_timestamp" structs:"state_change_timestamp"`
        State                       string    `firestore:"state" structs:"state"`
        ConfirmationHash            string    `firestore:"confirmation_hash" structs:"confirmation_hash"`
        ResolvedHash                string    `firestore:"resolved_hash" structs:"resolved_hash"`
        Frequency                   int       `firestore:"frequency_seconds" structs:"frequency_seconds"`
        TimeUntilReporting          int       `firestore:"time_until_reporting_seconds" structs:"time_until_reporting_seconds"`
        AllowedResponseTime         int       `firestore:"allowed_response_time_seconds" structs:"allowed_response_time_seconds"`
}

func AssertAlmostEqual(t *testing.T, expected, actual SiteInfo) {
        assert.NotEqual(t, expected.LastChangeTimestamp, actual.LastChangeTimestamp)
        actual.LastChangeTimestamp = expected.LastChangeTimestamp
        assert.Equal(t, expected, actual)
}

func (s SiteInfo) GetFrequency() time.Duration {
        return time.Duration(s.Frequency) * time.Second
}

func (s SiteInfo) GetTimeUntilReporting() time.Duration {
        return time.Duration(s.TimeUntilReporting) * time.Second
}

func (s SiteInfo) GetAllowedResponseTime() time.Duration {
        return time.Duration(s.AllowedResponseTime) * time.Second
}

type Database interface {
        GetSiteInfo(field string, value interface{}) (SiteInfo, error)
        GetAllSites() ([]SiteInfo, error)
        UpdateSite(siteInfo SiteInfo) error
        AddSite(info SiteInfo) error
        Clear() error
}