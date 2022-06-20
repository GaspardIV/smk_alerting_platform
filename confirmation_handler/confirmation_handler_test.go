package confirmation_handler

import (
        "fmt"
        "github.com/stretchr/testify/assert"
        "log"
        "net/http"
        "net/http/httptest"
        "os"
        "smk_alerting_platform/pkg"
        "testing"
)

var (
        siteRunning = pkg.SiteInfo{
                ID:                          "id0",
                Url:                         "site",
                PrimaryAdministratorEmail:   "primary_admin",
                SecondaryAdministratorEmail: "secondary_admin",
                State:                       pkg.Running,
                ConfirmationHash:            "confirmation_link",
                ResolvedHash:                "resolved_link",
                Frequency:                   10,
                TimeUntilReporting:          10,
                AllowedResponseTime:         10,
        }
        siteUnavailable = pkg.SiteInfo{
                ID:                          "id1",
                Url:                         "site",
                PrimaryAdministratorEmail:   "primary_admin",
                SecondaryAdministratorEmail: "secondary_admin",
                State:                       pkg.Unavailable,
                ConfirmationHash:            "confirmation_link1",
                ResolvedHash:                "resolved_link1",
                Frequency:                   10,
                TimeUntilReporting:          10,
                AllowedResponseTime:         10,
        }
        siteNotified = pkg.SiteInfo{
                ID:                          "id2",
                Url:                         "site",
                PrimaryAdministratorEmail:   "primary_admin",
                SecondaryAdministratorEmail: "secondary_admin",
                State:                       pkg.Notified,
                ConfirmationHash:            "confirmation_link2",
                ResolvedHash:                "resolved_link2",
                Frequency:                   10,
                TimeUntilReporting:          10,
                AllowedResponseTime:         10,
        }
        siteConfirmed = pkg.SiteInfo{
                ID:                          "id3",
                Url:                         "site",
                PrimaryAdministratorEmail:   "primary_admin",
                SecondaryAdministratorEmail: "secondary_admin",
                State:                       pkg.Confirmed,
                ConfirmationHash:            "confirmation_link3",
                ResolvedHash:                "resolved_link3",
                Frequency:                   10,
                TimeUntilReporting:          10,
                AllowedResponseTime:         10,
        }
        siteUnexpected = pkg.SiteInfo{
                Url:   "site",
                State: "some_state",
        }
)

func TestConfirmationHandlerWithDatabase_ParamHashIsMissing(t *testing.T) {
        rr := httptest.NewRecorder()
        os.Setenv("FUNCTION_BASE_URL", "http://localhost:8080")
        hash := ""
        cloudFunctionHandlerUrl := fmt.Sprintf("%v/%v-handler?hash=%v", os.Getenv("FUNCTION_BASE_URL"), "resolved", hash)
        request, _ := http.NewRequest("GET", cloudFunctionHandlerUrl, nil)
        var database = pkg.CreateLocalDatabase([]pkg.SiteInfo{siteUnexpected, siteRunning, siteUnavailable})

        ConfirmationHandlerWithDatabase(rr, request, database)
        assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestConfirmationHandlerWithDatabase_HashOutdatedSiteRunning(t *testing.T) {
        rr := httptest.NewRecorder()
        os.Setenv("FUNCTION_BASE_URL", "http://localhost:8080")
        hash := siteRunning.ConfirmationHash
        cloudFunctionHandlerUrl := fmt.Sprintf("%v/%v-handler?hash=%v", os.Getenv("FUNCTION_BASE_URL"), "resolved", hash)
        request, _ := http.NewRequest("GET", cloudFunctionHandlerUrl, nil)
        var database = pkg.CreateLocalDatabase([]pkg.SiteInfo{siteUnexpected, siteRunning, siteUnavailable})
        ConfirmationHandlerWithDatabase(rr, request, database)
        assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestConfirmationHandlerWithDatabase_HashOutdatedSiteUnavailable(t *testing.T) {
        rr := httptest.NewRecorder()
        os.Setenv("FUNCTION_BASE_URL", "http://localhost:8080")
        hash := siteUnavailable.ConfirmationHash
        cloudFunctionHandlerUrl := fmt.Sprintf("%v/%v-handler?hash=%v", os.Getenv("FUNCTION_BASE_URL"), "resolved", hash)
        request, _ := http.NewRequest("GET", cloudFunctionHandlerUrl, nil)
        log.Println(request)
        var database = pkg.CreateLocalDatabase([]pkg.SiteInfo{siteUnexpected, siteRunning, siteUnavailable, siteConfirmed, siteNotified})
        ConfirmationHandlerWithDatabase(rr, request, database)
        assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestConfirmationHandlerWithDatabase_SiteConfirmed(t *testing.T) {
        rr := httptest.NewRecorder()
        os.Setenv("FUNCTION_BASE_URL", "http://localhost:8080")
        hash := siteConfirmed.ConfirmationHash
        cloudFunctionHandlerUrl := fmt.Sprintf("%v/%v-handler?hash=%v", os.Getenv("FUNCTION_BASE_URL"), "resolved", hash)
        request, _ := http.NewRequest("GET", cloudFunctionHandlerUrl, nil)
        log.Println(request)
        var database = pkg.CreateLocalDatabase([]pkg.SiteInfo{siteUnexpected, siteRunning, siteUnavailable, siteConfirmed, siteNotified})
        ConfirmationHandlerWithDatabase(rr, request, database)
        assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestConfirmationHandlerWithDatabase_SiteNotified(t *testing.T) {
        rr := httptest.NewRecorder()
        os.Setenv("FUNCTION_BASE_URL", "http://localhost:8080")
        hash := siteNotified.ConfirmationHash
        cloudFunctionHandlerUrl := fmt.Sprintf("%v/%v-handler?hash=%v", os.Getenv("FUNCTION_BASE_URL"), "resolved", hash)
        request, _ := http.NewRequest("GET", cloudFunctionHandlerUrl, nil)
        log.Println(request)
        var database = pkg.CreateLocalDatabase([]pkg.SiteInfo{siteUnexpected, siteRunning, siteUnavailable, siteConfirmed, siteNotified})
        ConfirmationHandlerWithDatabase(rr, request, database)
        assert.Equal(t, http.StatusOK, rr.Code)
}