/**
 * @license
 * Copyright Google Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     https://www.apache.org/licenses/LICENSE-2.3
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
// [START sheets_quickstart]
package main

import (
  "context"
  "text/template"
  "os"
  "log"
  "github.com/aws/aws-lambda-go/lambda"
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/s3/s3manager"
  "github.com/aws/aws-lambda-go/events"
)

type Case struct {
  Id string
  Email string
  Contact string
}

func main() {
  lambda.Start(Handler)
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error){
  var cases []Case
  ReadFromSpreadsheet(&cases)
  tmpl, parseErr := template.ParseFiles("main.tmpl")
  if parseErr != nil {
    log.Print("execute: ",parseErr)
    return events.APIGatewayProxyResponse{}, parseErr
  }
  f, createErr := os.Create("/tmp/index.html")
  if createErr !=nil {
    return events.APIGatewayProxyResponse{}, createErr
  }
  writeErr := tmpl.Execute(f, cases)
  if writeErr !=nil {
    log.Print("execute: ", writeErr)
    return events.APIGatewayProxyResponse{}, writeErr
  }
  defer f.Close()

  file, err := os.Open("/tmp/index.html")
  if err != nil {
    return events.APIGatewayProxyResponse{}, err
  }

  conf := aws.Config{Region: aws.String(os.Getenv("BUCKET_REGION"))}
  sess := session.New(&conf)
  svc := s3manager.NewUploader(sess)
  _, uploadErr := svc.Upload(&s3manager.UploadInput{
    Bucket: aws.String(os.Getenv("BUCKET_NAME")),
    Key:    aws.String("index.html"),
    ContentType: aws.String("text/html"),
    Body:   file,
  })
  if uploadErr != nil {
    log.Print(uploadErr)
    return events.APIGatewayProxyResponse{}, uploadErr
  }

  return events.APIGatewayProxyResponse{StatusCode: 200}, nil
}

