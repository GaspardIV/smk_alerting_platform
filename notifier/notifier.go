package notifier

import (
        "smk_alerting_platform/pkg"
        "time"

        "crypto/sha256"
        "encoding/json"
        "fmt"
        "log"
        "math/rand"
        "net/http"
        "os"

        "github.com/sendgrid/sendgrid-go"
        "github.com/sendgrid/sendgrid-go/helpers/mail"

        secretmanager "cloud.google.com/go/secretmanager/apiv1"
        "context"
        secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

type messageType string

const (
        Confirmation messageType = "confirmation"
        Resolved                 = "resolved"
)

type notifierPayload struct {
        Url string `json:"url"`
}

func RetrieveApiKey() (string, error) {
        ctx := context.Background()
        client, err := secretmanager.NewClient(ctx)
        if err != nil {
                return "", fmt.Errorf("Failed to create secretmanager client: %v", err)
        }
        defer client.Close()
        req := &secretmanagerpb.AccessSecretVersionRequest{
                Name: fmt.Sprintf("projects/%v/secrets/sendgrid-api-key/versions/latest", pkg.ProjectNumber),
        }

        result, err := client.AccessSecretVersion(ctx, req)
        if err != nil {
                return "", fmt.Errorf("Failed to access secret version: %v", err)
        }
        return string(result.Payload.Data), nil
}

func generateHash() string {
        salt := make([]byte, 8)
        rand.Read(salt)
        data := []byte(string(salt))
        hash := sha256.Sum256(data)
        return fmt.Sprintf("%x", hash)
}

func generateMessage(url string, siteInfo *pkg.SiteInfo, database pkg.Database, msgType messageType, hash string) string {
        var messageSuffix string
        if msgType == Confirmation {
                siteInfo.ConfirmationHash = hash
                messageSuffix = "to confirm you are working on the issue.\n"
        } else {
                siteInfo.ResolvedHash = hash
                messageSuffix = "once you resolve the issue.\n"
        }
        cloudFunctionHandlerUrl := fmt.Sprintf("%v/%v-handler?hash=%v", os.Getenv("FUNCTION_BASE_URL"), msgType, hash)
        message := fmt.Sprintf("Visit the following link: %v %v", cloudFunctionHandlerUrl, messageSuffix)
        return message
}

func sendEmail(url string, adminEmailAddress string, messageSuffix string) error {
        subject := fmt.Sprintf("smk-alerting-platform: Your site %v is unavailable.", url)
        message := fmt.Sprintf("%v\n%v", subject, messageSuffix)
        log.Printf("Sending email to %v with the content:\n%v", adminEmailAddress, message)
        if os.Getenv("TURN_ON_EMAIL_SENDING") != "true" {
                log.Printf("Email sending turned off. Set TURN_ON_EMAIL_SENDING=true variable in to turn on the email sending.")
                return nil
        }

        from := mail.NewEmail("smk-alerting-platform", "ka.konecki@student.uw.edu.pl")
        to := mail.NewEmail(adminEmailAddress, adminEmailAddress)
        email := mail.NewSingleEmailPlainText(from, subject, to, message)
        apiKey, apiKeyErr := RetrieveApiKey()
        if apiKeyErr != nil {
                return apiKeyErr
        }
        client := sendgrid.NewSendClient(apiKey)
        response, err := client.Send(email)
        if err != nil || response.StatusCode != http.StatusAccepted {
                return fmt.Errorf("error while sending an email regarding site %v info, err: %v, response: %v", url, err, response)
        } else {
                log.Printf("Email has been sent successfully.")
        }
        return nil
}

func notifyAdmin(url string, database pkg.Database) error {
        siteInfo, err := database.GetSiteInfo("url", url)
        if err != nil {
                return fmt.Errorf("Error while getting site %v info %v", url, err)
        }

        if pkg.Unavailable == siteInfo.State {
                log.Printf("Site %v is unavailable. Notifying primary administrator %v...", url, siteInfo.PrimaryAdministratorEmail)
                emailText := generateMessage(url, &siteInfo, database, Confirmation, generateHash())
                emailText += generateMessage(url, &siteInfo, database, Resolved, generateHash())
                err := sendEmail(url, siteInfo.PrimaryAdministratorEmail, emailText)
                if err != nil {
                        return err
                }
                err = database.UpdateSite(siteInfo)
                if err != nil {
                        return fmt.Errorf("Error while updating site %v info %v", url, err)
                }
                siteInfo.State = pkg.Notified
                siteInfo.StateChangeTimestamp = time.Now()
                err = database.UpdateSite(siteInfo)
                if err != nil {
                        return fmt.Errorf("Error while updating site %v info %v", url, err)
                }
        } else if pkg.Notified == siteInfo.State {
                log.Printf("Primary administrator %v of site %v has already been notified. Notifying secondary administrator... %v", siteInfo.PrimaryAdministratorEmail, url, siteInfo.SecondaryAdministratorEmail)
                emailText := generateMessage(url, &siteInfo, database, Resolved, siteInfo.ResolvedHash)
                err := sendEmail(url, siteInfo.SecondaryAdministratorEmail, emailText)
                if err != nil {
                        return err
                }
                err = database.UpdateSite(siteInfo)
                if err != nil {
                        return fmt.Errorf("Error while updating site %v info %v", url, err)
                }
        }
        return nil
}

func NotifierWithDatabase(w http.ResponseWriter, r *http.Request, database pkg.Database) {
        var payload notifierPayload
        if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
                message := fmt.Sprintf("json.NewDecoder: %v", err)
                log.Println(message)
                w.Write([]byte(message))
                w.WriteHeader(http.StatusBadRequest)
                return
        }
        err := notifyAdmin(payload.Url, database)
        if err != nil {
                log.Println(err.Error())
                w.Write([]byte(err.Error()))
                w.WriteHeader(http.StatusInternalServerError)
                return
        }
        w.WriteHeader(http.StatusOK)
}

func Notifier(w http.ResponseWriter, r *http.Request) {
        database, err := pkg.CreateDatabase()
        if err != nil {
                w.Write([]byte(fmt.Sprintf("Could not connect to database: %v", err)))
                w.WriteHeader(http.StatusInternalServerError)
        }
        NotifierWithDatabase(w, r, database)
}