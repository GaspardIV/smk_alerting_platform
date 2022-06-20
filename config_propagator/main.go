package main

import (
        "smk_alerting_platform/pkg"

        "log"
        "time"
)

func generateSiteInfos(config []serviceConfig) []pkg.SiteInfo {
        siteInfos := make([]pkg.SiteInfo, 0)
        for _, c := range config {
                siteInfo := pkg.SiteInfo{
                        Url:                         c.Url,
                        PrimaryAdministratorEmail:   c.PrimaryAdministratorEmail,
                        SecondaryAdministratorEmail: c.SecondaryAdministratorEmail,
                        StateChangeTimestamp:        time.Now(),
                        State:                       pkg.Running,
                        Frequency:                   c.Frequency,
                        TimeUntilReporting:          c.TimeUntilReporting,
                        AllowedResponseTime:         c.AllowedResponseTime,
                }
                siteInfos = append(siteInfos, siteInfo)
        }
        return siteInfos
}

func setSiteInfos(siteInfos []pkg.SiteInfo, database pkg.Database) {
        err := database.Clear()
        if err != nil {
                log.Fatalln("database.Clear", err)
        }
        for _, info := range siteInfos {
                err = database.AddSite(info)
                if err != nil {
                        log.Fatalln("database.AddSite", err)
                }
        }
}

func main() {
        config := readConfig("config.json")
        infos := generateSiteInfos(config)
        database, err := pkg.CreateDatabase()
        if err != nil {
                log.Fatalln("CreateDatabase", err)
        }
        setSiteInfos(infos, database)
}