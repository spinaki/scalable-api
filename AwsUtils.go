package main

import (
        "github.com/aws/aws-sdk-go/aws"
        "github.com/aws/aws-sdk-go/aws/session"
        "os"
        "bytes"
        "github.com/aws/aws-sdk-go/aws/awsutil"
        "fmt"
        "log"
        "github.com/aws/aws-sdk-go/aws/credentials"
        "github.com/aws/aws-sdk-go/service/s3"
        "github.com/aws/aws-sdk-go/service/ses"
)

var (
        s3Service *s3.S3
        sesService *ses.SES
)

// initializes the external hooks like SES, S3.
func setupAWS(workspace string) {
        aws_access_key_id := "AWS_ACCESS_KEY_ID"
        aws_secret_access_key := "AWS_SECRET_KEY"
        token := ""
        creds := credentials.NewStaticCredentials(aws_access_key_id, aws_secret_access_key, token)
        _, err := creds.Get()
        if err != nil {
                log.Printf("bad credentials: %s \n", err)
        }
        cfg := aws.NewConfig().WithCredentials(creds)
        sess := session.New(&aws.Config{Region: aws.String("us-west-2")})
        s3Service = s3.New(sess, cfg)
        sesService = ses.New(sess, cfg)
}

// workspacePath: top lever directory containing the android project. It should have the main / project gradle file.

func UploadToS3(s3Key string, fileToUpload string)  (*string, error) {
        file, err := os.Open(fileToUpload)
        if err != nil {
                log.Printf("Error Opening  file %s : %s\n", fileToUpload, err)
                return nil, err
        }
        defer file.Close()
        fileInfo, _ := file.Stat()
        size := fileInfo.Size()
        buffer := make([]byte, size) // read file content to buffer

        file.Read(buffer)
        fileBytes := bytes.NewReader(buffer)
        key := s3Key
        params := &s3.PutObjectInput{
                Bucket: aws.String("YOUR_BUCKET_NAME"), // TODO: add your own bucket name
                Key: aws.String(key),
                Body: fileBytes,
                ContentLength: aws.Int64(size),
                ContentType: aws.String("YOUR_CONTENT_TYPE"), // TODO: add your own content type
                ACL:aws.String(s3.ObjectCannedACLPublicRead), // Public READ ?
        }
        resp, err := s3Service.PutObject(params)
        if err != nil {
                log.Printf("S3 Error Response: %s", err)
                return nil, err
        }
        log.Printf("S3 Upload Successful Response %s \n", awsutil.StringValue(resp))
        urlString := "" // TODO: create the URL based on domain and the S3 key
        return &urlString, nil
}

func SendSuccessEmail(toAddress string, url string) error {
        body := "Success! " + url +
                "<br> SOME SUCCESS MESSAGE"
        subject := "SOME SUCCESS SUBJECT"
        return SendSESEmail(toAddress, body, subject)
}

func SendFailureEmail(toAddress string, job Job, err error) {
        body := fmt.Sprint("Failed Job: ", job.Name, "<br> NextParallelJobs", job.NextParallelJobs,
                "<br> Cause:", err.Error(), "<br><br> Payload:", "<br>", job.JobPayload.Time, "<br>",
                job.JobPayload.Email, "<br>", job.JobPayload.Id, "<br>", job.JobPayload.Data, "<br>",
                job.JobPayload.RemoteUrl)
        subject:= "SOME FAILURE SUBJECT"
        SendSESEmail(toAddress, body, subject)
}
// http://docs.aws.amazon.com/sdk-for-go/api/service/ses/
func SendSESEmail(toAddress string, body string, subject string) error {
        params := &ses.SendEmailInput{
                Destination: &ses.Destination{ // Required
                        ToAddresses: []*string{
                                aws.String(toAddress), // Required
                        },
                },
                Message: &ses.Message{ // Required
                        Body: &ses.Body{ // Required
                                Html: &ses.Content{
                                        Data:    aws.String(body), // Required
                                },
                                Text: &ses.Content{
                                        Data:    aws.String(body), // Required
                                },
                        },
                        Subject: &ses.Content{ // Required
                                Data:    aws.String(subject), // Required
                        },
                },
                Source: aws.String("BuildBot <buildbot@email.com>"), // Required TODO: ADD YOUR OWN SOURCE ADDRESS
                ReplyToAddresses: []*string{
                        aws.String("do-not-reply@email.com"), // Required TODO: ADD YOUR OWN RE ADDRESS
                },
        }
        resp, err := sesService.SendEmail(params)

        if err != nil {
                // Print the error, cast err to awserr.Error to get the Code and
                // Message from an error.
                log.Println("SES SendEmail Error: " + err.Error())
                return err
        }
        // Pretty-print the response data.
        log.Println("SES Success Response: " + resp.String())
        return nil
}
