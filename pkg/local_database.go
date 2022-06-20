package pkg

import (
        "errors"
        "reflect"
        "time"

        "github.com/fatih/structs"
)

type localDatabase struct {
        sites []SiteInfo
}

var siteNotFound = errors.New("site not found")

func CreateLocalDatabase(sites []SiteInfo) *localDatabase {
        return &localDatabase{sites: sites}
}

func (db *localDatabase) GetSiteInfo(field string, value interface{}) (SiteInfo, error) {
        for _, siteInfo := range db.sites {
                if m := structs.Map(siteInfo); m[field] == value {
                        return siteInfo, nil
                }
        }
        return SiteInfo{}, siteNotFound
}

func checkOp(value1 interface{}, op string, value2 interface{}) bool { // for now we only use "in" operator
        if op == "in" {
                switch reflect.TypeOf(value2).Kind() {
                case reflect.Slice:
                        s := reflect.ValueOf(value2)

                        for i := 0; i < s.Len(); i++ {
                                if s.Index(i).Interface() == value1 {
                                        return true
                                }
                        }
                }
                return false
        }
        return false
}

func (db *localDatabase) GetAllSites() ([]SiteInfo, error) {
        return db.sites, nil
}

func (db *localDatabase) UpdateSite(siteInfo SiteInfo) error {
        for i, site := range db.sites {
                if site.ID == siteInfo.ID {
                        siteInfo.LastChangeTimestamp = time.Now()
                        db.sites[i] = siteInfo
                        return nil
                }
        }
        return siteNotFound
}

func (db *localDatabase) Clear() error {
        db.sites = nil
        return nil
}

func (db *localDatabase) AddSite(info SiteInfo) error {
        db.sites = append(db.sites, info)
        return nil
}