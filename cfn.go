package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/awslabs/aws-sdk-go/service/cloudformation"

	"log"
)

func listCfn(region string, verbose bool) (cfnNumber int) {
	svc := cloudformation.New(session.New(&aws.Config{Region: aws.String(region)}))
	params := &cloudformation.ListStacksInput{
		StackStatusFilter: []*string{
			aws.String("CREATE_COMPLETE"),
			aws.String("CREATE_FAILED"),
			aws.String("CREATE_IN_PROGRESS"),
			aws.String("DELETE_FAILED"),
			aws.String("DELETE_IN_PROGRESS"),
			aws.String("ROLLBACK_FAILED"),
			aws.String("ROLLBACK_IN_PROGRESS"),
			aws.String("UPDATE_COMPLETE"),
			aws.String("UPDATE_COMPLETE_CLEANUP_IN_PROGRESS"),
			aws.String("UPDATE_IN_PROGRESS"),
			aws.String("UPDATE_ROLLBACK_COMPLETE"),
			aws.String("UPDATE_ROLLBACK_COMPLETE_CLEANUP_IN_PROGRESS"),
			aws.String("UPDATE_ROLLBACK_FAILED"),
			aws.String("UPDATE_ROLLBACK_IN_PROGRESS"),
		},
	}
	resp, err := svc.ListStacks(params)
	if err != nil {
		log.Fatal("Cannot get CFN data: ", err)
	}

	if verbose {
		for _, name := range resp.StackSummaries {
			log.Println("CFN Stack Name: ", *name.StackName, " Status: ", *name.StackStatus)
		}
	}
	cfnNumber = len(resp.StackSummaries)
	return cfnNumber
}
