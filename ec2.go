package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"

	"log"
)

func listEC2(region string, verbose bool) (ec2Number, ec2RunningNumber, ec2RunningWindows int) {
	ec2client := ec2.New(session.New(&aws.Config{Region: aws.String(region)}))

	resp, err := ec2client.DescribeInstances(nil)
	if err != nil {
		log.Fatal("Cannot get EC2 data: ", err)
	}

	// resp has all of the response data, pull out instance IDs:
	for idx := range resp.Reservations {
		for _, inst := range resp.Reservations[idx].Instances {
			ec2Number++
			if *inst.State.Name == "running" {
				ec2RunningNumber++
				if p2s(inst.Platform) == "windows" {
					ec2RunningWindows++
				}
			}
			if verbose {
				log.Println("Instance ID: ", *inst.InstanceId)
				for _, tag := range inst.Tags {
					if *tag.Key == "Name" {
						log.Println("-->Instance Name: ", *tag.Value)
					}
				}
				log.Println("-->Instance State: ", *inst.State.Name)
				if p2s(inst.Platform) == "" {
					log.Println("-->Instance Platform: Linux like")
				} else {
					log.Println("-->Instance Platform: ", p2s(inst.Platform))
				}
			}
		}
	}
	return ec2Number, ec2RunningNumber, ec2RunningWindows
}
