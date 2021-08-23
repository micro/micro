package s3

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	sthree "github.com/aws/aws-sdk-go/service/s3"
	store2 "github.com/micro/micro/v3/service/events/store"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"
)

type backupImpl struct {
	client *sthree.S3
	opts   Options
}

func NewBackup(opts ...Option) store2.Backup {
	// parse the options
	options := Options{Secure: true}
	for _, o := range opts {
		o(&options)
	}

	sess := session.Must(session.NewSession(&aws.Config{
		Endpoint:    &options.Endpoint,
		Region:      &options.Region,
		Credentials: credentials.NewStaticCredentials(options.AccessKeyID, options.SecretAccessKey, ""),
	}))
	client := sthree.New(sess)

	return &backupImpl{client: client, opts: options}
}

func (i *backupImpl) Snapshot(st store.Store) error {
	// find latest S3 backup file
	out, err := i.client.ListObjects(&sthree.ListObjectsInput{
		Bucket: aws.String(i.opts.Bucket),
		Prefix: aws.String("micro/eventsBackup/"),
	})
	if err != nil {
		logger.Errorf("Error retrieving objects from bucket %s", err)
		return err
	}
	latest := time.Now().Add(-12 * time.Hour).Format("2006010215") // default to now
	// list operations results are returned in utf8 binary order
	for _, obj := range out.Contents {
		parts := strings.Split(*obj.Key, "/")
		objName := parts[len(parts)-1]
		if match, _ := regexp.MatchString("^[0-9]{10}", objName); !match {
			continue
		}
		latest = objName
	}
	// Work out what the next backup period should be, backup
	t, err := time.Parse("2006010215", latest[:10])
	if err != nil {
		ferr := fmt.Errorf("error parsing time %s", err)
		logger.Error(ferr.Error())
		return ferr
	}
	// loop until all the eligible backup periods are completed
	rollBackNum := 0
dateLoop:
	for {
		t = t.Add(1 * time.Hour)
		if time.Since(t) < 1*time.Hour {
			// our work is done here
			logger.Info("No more backups to be completed right now")
			break dateLoop
		}
		next := t.Format("2006010215")
		logger.Infof("Performing backup for %s", next)
		// read all the entries
		uploadNum := 0
		offset := 0
		limit := 100
		buf := &bytes.Buffer{}
		for {
			recs, err := st.Read("/"+next, store.ReadSuffix(), store.ReadLimit(uint(limit)), store.ReadOffset(uint(offset)))
			if err == store.ErrNotFound || len(recs) == 0 {
				logger.Errorf("No records found")
				break
			} else if err != nil {
				logger.Errorf("Error reading recs %s", err)
				return err
			}
			logger.Infof("Processing %d records", len(recs))
			// process
			for _, rec := range recs {
				// 1 record per line
				buf.Write(rec.Value)
				buf.WriteString("\n")
				if buf.Len() >= 50000000 { // 50MB to help constrain the process size
					uploadNum++
					if err := i.uploadToS3(fmt.Sprintf("%s-%d", next, uploadNum), buf); err != nil {
						rollBackNum = uploadNum
						break dateLoop
					}
				}
			}
			offset += limit
		}
		// write anything leftover to file and upload
		uploadNum++
		if err := i.uploadToS3(fmt.Sprintf("%s-%d", next, uploadNum), buf); err != nil {
			rollBackNum = uploadNum
			break dateLoop
		}
	}
	if rollBackNum > 0 {
		// rollback this date so we can retry later
		d := t.Format("2006010215")
		for num := 1; num <= rollBackNum; num++ {
			if _, err := i.client.DeleteObject(&sthree.DeleteObjectInput{
				Bucket: aws.String(i.opts.Bucket),
				Key:    aws.String(fmt.Sprintf("micro/eventsBackup/%s-%s", d, num)),
			}); err != nil {
				logger.Errorf("Error during rollback %s", err)
			}
		}
	}
	logger.Infof("Backup complete")

	return nil
}

func (i *backupImpl) uploadToS3(key string, buf *bytes.Buffer) error {
	if buf.Len() == 0 {
		return nil
	}
	_, err := i.client.PutObject(&sthree.PutObjectInput{
		Bucket: aws.String(i.opts.Bucket),
		Body:   bytes.NewReader(buf.Bytes()),
		Key:    aws.String(fmt.Sprintf("micro/eventsBackup/%s", key)),
		ACL:    aws.String("private"),
	})
	if err != nil {
		logger.Errorf("Error uploading %s", err)
		return err
	}
	buf.Reset()
	return nil
}
