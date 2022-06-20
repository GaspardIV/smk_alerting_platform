package confirmation_handler

import (
        "fmt"
        "log"
        "net/http"
        "smk_alerting_platform/pkg"
        "time"
)

func ConfirmationHandler(w http.ResponseWriter, r *http.Request) {
        database, done := getDatabase(w)
        if done {
                return
        }
        ConfirmationHandlerWithDatabase(w, r, database)
}

func ConfirmationHandlerWithDatabase(w http.ResponseWriter, r *http.Request, database pkg.Database) {
        hash, done := getHashFromQuery(w, r)
        if done {
                return
        }

        siteInfo, done := getSiteInfoFromHash(w, database, hash)
        if done {
                return
        }

        setStateConfirmed(w, siteInfo, database)
}

func getDatabase(w http.ResponseWriter) (pkg.Database, bool) {
        database, err := pkg.CreateDatabase()
        if err != nil {
                w.Write([]byte(fmt.Sprintf("create database: %v", err)))
                w.WriteHeader(http.StatusInternalServerError)
                return nil, true
        }
        return database, false
}

func getHashFromQuery(w http.ResponseWriter, r *http.Request) (string, bool) {
        hashes, ok := r.URL.Query()["hash"]
        if !ok || len(hashes) < 1 || len(hashes[0]) < 1 {
                log.Println("Url Param 'hash' is missing")
                w.WriteHeader(http.StatusBadRequest)
                return "", true
        }
        hash := hashes[0]
        return hash, false
}

func getSiteInfoFromHash(w http.ResponseWriter, database pkg.Database, hash string) (pkg.SiteInfo, bool) {
        siteInfo, err := database.GetSiteInfo("confirmation_hash", hash)

        if err != nil || (siteInfo.State != pkg.Notified) {
                log.Println("hash outdated")
                w.WriteHeader(http.StatusBadRequest)
                w.Write([]byte("Link has expired. Cannot mark as confirmed."))
                return pkg.SiteInfo{}, true
        }
        return siteInfo, false
}

func setStateConfirmed(w http.ResponseWriter, siteInfo pkg.SiteInfo, database pkg.Database) {
        siteInfo.State = pkg.Confirmed
        siteInfo.StateChangeTimestamp = time.Now()
        err := database.UpdateSite(siteInfo)
        if err != nil {
                log.Printf("Error while updating site %v info %v", siteInfo.Url, err)
                w.WriteHeader(http.StatusInternalServerError)
                w.Write([]byte("Error while updating site info. Please try again, or contact us."))
        } else {
                w.WriteHeader(http.StatusOK)
                w.Write([]byte("You have confirmed that you have been informed about problem with your site. "))
                w.Write([]byte("Mark problem as resolved to start site checking again."))
                log.Printf("confirmed %v", siteInfo.Url)
        }
}