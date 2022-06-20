package main

import (
        "encoding/json"
        "fmt"
        "io/ioutil"
        "log"
        "os"
)

type serviceConfig struct {
        Url                         string
        PrimaryAdministratorEmail   string
        SecondaryAdministratorEmail string
        Frequency                   int
        TimeUntilReporting          int
        AllowedResponseTime         int
}

func readConfig(filename string) []serviceConfig {
        fp, err := os.Open(filename)
        if err != nil {
                log.Fatalln(err)
        }
        defer fp.Close()

        bytes, _ := ioutil.ReadAll(fp)
        var config []serviceConfig
        if err := json.Unmarshal(bytes, &config); err != nil {
                log.Fatalln(err)
        }

        return config
}

func (s serviceConfig) String() string {
        return fmt.Sprintf(
                "{\n  Url: %v\n  Administrators: [%v, %v]\n  Frequency: %v\n  TimeUntilReporting: %v\n  AllowedResponseTime: %v\n}",
                s.Url,
                s.PrimaryAdministratorEmail,
                s.SecondaryAdministratorEmail,
                s.Frequency,
                s.TimeUntilReporting,
                s.AllowedResponseTime)
}