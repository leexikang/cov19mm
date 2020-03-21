package main

import (
	//"context"
	"fmt"
	"log"
	"time"
  "os"
  "context"
  "github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudfront"
  "github.com/aws/aws-lambda-go/events"
)

func main(){
  lambda.Start(Handler)
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
  now := time.Now()
  conf := aws.Config{Region: aws.String(os.Getenv("AWS_REGION"))}
  sess := session.New(&conf)
  svc := cloudfront.New(sess)
  _, err := svc.CreateInvalidation(&cloudfront.CreateInvalidationInput{
    DistributionId: aws.String(os.Getenv("CLOUDFRON_ID")),
    InvalidationBatch: &cloudfront.InvalidationBatch{
      CallerReference: aws.String(
        fmt.Sprintf("goinvali%s", now.Format("2006/01/02,15:04:05"))),
        Paths: &cloudfront.Paths{
          Quantity: aws.Int64(1),
          Items: []*string{
            aws.String("/*.html"),
          },
        },
      },
    })
    if err != nil {
      log.Print("Error ", err)
      return events.APIGatewayProxyResponse{}, err
    }
    return events.APIGatewayProxyResponse{StatusCode: 200}, nil
  }


