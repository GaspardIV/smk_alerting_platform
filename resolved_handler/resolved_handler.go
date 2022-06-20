package resolved_handler

import (
        "fmt"
        "log"
        "net/http"
        "smk_alerting_platform/pkg"
        "time"
)

func ResolvedHandler(w http.ResponseWriter, r *http.Request) {
        database, done := getDatabase(w)
        if done {
                return
        }
        ResolvedHandlerWithDatabase(w, r, database)
}

func ResolvedHandlerWithDatabase(w http.ResponseWriter, r *http.Request, database pkg.Database) {
        hash, done := getHashFromQuery(w, r)
        if done {
                return
        }

        siteInfo, done := getSiteInfoFromHash(w, database, hash)
        if done {
                return
        }

        setStateResolved(w, siteInfo, database)
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
        siteInfo, err := database.GetSiteInfo("resolved_hash", hash)

        if err != nil || (siteInfo.State != pkg.Notified && siteInfo.State != pkg.Confirmed) {
                log.Println("hash outdated")
                w.WriteHeader(http.StatusBadRequest)
                w.Write([]byte("Link has expired. Cannot mark as resolved."))
                return pkg.SiteInfo{}, true
        }
        return siteInfo, false
}

func setStateResolved(w http.ResponseWriter, siteInfo pkg.SiteInfo, database pkg.Database) {
        siteInfo.State = pkg.Running
        siteInfo.StateChangeTimestamp = time.Now()
        err := database.UpdateSite(siteInfo)

        if err != nil {
                log.Printf("Error while updating site %v info %v", siteInfo.Url, err)
                w.WriteHeader(http.StatusInternalServerError)
                w.Write([]byte("Error while updating site info. Please try again, or contact us."))
        } else {
                w.WriteHeader(http.StatusOK)
                w.Write([]byte("Problem has been marked as Resolved."))
                log.Printf("resolved %v", siteInfo.Url)
        }
}