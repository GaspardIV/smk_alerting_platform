package pkg

import (
        "context"

        "cloud.google.com/go/firestore"
        "github.com/fatih/structs"
        "google.golang.org/api/iterator"
)

type firestoreDatabase struct {
        client *firestore.Client
}

func CreateDatabase() (Database, error) {
        client, err := firestore.NewClient(context.Background(), projectID)
        if err != nil {
                return firestoreDatabase{}, err
        }
        return firestoreDatabase{client: client}, nil
}

func (db firestoreDatabase) GetSiteInfo(field string, value interface{}) (SiteInfo, error) {
        iter := db.client.Collection("sites").Where(field, "==", value).Documents(context.Background())
        doc, err := iter.Next()
        if err != nil {
                return SiteInfo{}, err
        }

        var siteInfo SiteInfo
        err = doc.DataTo(&siteInfo)
        if err != nil {
                return SiteInfo{}, err
        }
        siteInfo.ID = doc.Ref.ID
        return siteInfo, nil
}

func (db firestoreDatabase) GetAllSites() ([]SiteInfo, error) {
        iter := db.client.Collection("sites").Documents(context.Background())
        docs, err := iter.GetAll()
        if err != nil {
                return nil, err
        }

        siteInfos := make([]SiteInfo, len(docs))
        for i, doc := range docs {
                var siteInfo SiteInfo
                err = doc.DataTo(&siteInfo)
                if err != nil {
                        return nil, err
                }
                siteInfo.ID = doc.Ref.ID

                siteInfos[i] = siteInfo
        }

        return siteInfos, nil
}

func (db firestoreDatabase) UpdateSite(siteInfo SiteInfo) error {
        m := structs.Map(siteInfo)
        m["last_change_timestamp"] = firestore.ServerTimestamp
        _, err := db.client.Collection("sites").Doc(siteInfo.ID).Set(context.Background(), m)
        return err
}

func (db firestoreDatabase) Clear() error {
        batchSize := 20
        ref := db.client.Collection("sites")
        ctx := context.Background()
        for {
                // Get a batch of documents
                iter := ref.Limit(batchSize).Documents(ctx)
                numDeleted := 0

                // Iterate through the documents, adding
                // a delete operation for each one to a
                // WriteBatch.
                batch := db.client.Batch()
                for {
                        doc, err := iter.Next()
                        if err == iterator.Done {
                                break
                        }
                        if err != nil {
                                return err
                        }

                        batch.Delete(doc.Ref)
                        numDeleted++
                }

                // If there are no documents to delete,
                // the process is over.
                if numDeleted == 0 {
                        return nil
                }

                _, err := batch.Commit(ctx)
                if err != nil {
                        return err
                }
        }
}

func (db firestoreDatabase) AddSite(info SiteInfo) error {
        m := structs.Map(info)
        m["last_change_timestamp"] = firestore.ServerTimestamp
        _, _, err := db.client.Collection("sites").Add(context.Background(), m)
        return err
}