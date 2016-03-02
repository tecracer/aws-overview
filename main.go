package main

import (
	"flag"
	logs "github.com/Sirupsen/logrus"
	logstash "github.com/Sirupsen/logrus/formatters/logstash"
	"log"
	"sync"
	"time"
)

var (
	regions = []string{
		"us-east-1",
		"us-west-1",
		"us-west-2",
		"eu-west-1",
		"ap-southeast-1",
		"ap-southeast-2",
		"ap-northeast-1",
		"sa-east-1",
	}
	verbose, machineReadable                                             bool
	logfile                                                              string
	s3Number                                                             int
	totalEC2Number, totalEC2RunningNumber, totalEC2RunningWindowsNumber  int
	totalElbNumber, totalElbWithoutEC2Number                             int
	totalRdsNumber, totalOrRdsNumber, totalMyRdsNumber, totalMsRdsNumber int
	totalLambdaNumber                                                    int
	totalCfnNumber                                                       int
	err                                                                  error
)

var wg sync.WaitGroup

func init() {
	flag.BoolVar(&verbose, "verbose", false, "Show detailed output")
	flag.BoolVar(&machineReadable, "machine-readable", false, "Machine-readable output")
	flag.StringVar(&logfile, "log-file", "", "Log file location")
	flag.Parse()
	if machineReadable {
		logs.SetFormatter(&logstash.LogstashFormatter{Type: "aws_overview", TimestampFormat: time.RFC822})
	} else {
		logs.SetFormatter(&logs.JSONFormatter{TimestampFormat: time.RFC822})
	}
}

func main() {
	// Make sure the credentials exists
	checkConfig()

	// Make sure we can create log file
	checkLogFile(logfile)

	for _, region := range regions {
		region := region
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Getting EC2 data
			rTotal, rRunning, rWindows := listEC2(region, verbose)
			if rTotal > 0 {
				logs.WithFields(logs.Fields{
					"EC2":               rTotal,
					"EC2Running":        rRunning,
					"EC2RunningWindows": rWindows,
					"Region":            region,
				}).Info(msg("EC2"))
			}
			totalEC2Number += rTotal
			totalEC2RunningNumber += rRunning
			totalEC2RunningWindowsNumber += rWindows
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			// Getting ELB data
			rElbTotal, rElbWithoutEC2Total := listElb(region, verbose)
			if rElbTotal > 0 {
				logs.WithFields(logs.Fields{
					"ELB":      rElbTotal,
					"ELBwoEC2": rElbWithoutEC2Total,
					"Region":   region,
				}).Info(msg("ELB"))
			}
			totalElbNumber += rElbTotal
			totalElbWithoutEC2Number += rElbWithoutEC2Total
		}()

		// Getting RDS data
		wg.Add(1)
		go func() {
			defer wg.Done()
			rRdsTotal, rRdsOTotal, rRdsMyTotal, rRdsMsTotal := listRds(region, verbose)
			if rRdsTotal > 0 {
				logs.WithFields(logs.Fields{
					"RDS":        rRdsTotal,
					"RDS_Oracle": rRdsOTotal,
					"RDS_MySQL":  rRdsMyTotal,
					"RDS_MSSQL":  rRdsMsTotal,
					"Region":     region,
				}).Info(msg("RDS"))
			}
			totalRdsNumber += rRdsTotal
			totalOrRdsNumber += rRdsOTotal
			totalMyRdsNumber += rRdsMyTotal
			totalMsRdsNumber += rRdsMsTotal
		}()

		// Getting Lambda data
		wg.Add(1)
		go func() {
			defer wg.Done()
			rLambdaTotal := listLambda(region, verbose)
			if rLambdaTotal > 0 {
				logs.WithFields(logs.Fields{
					"Lambda": rLambdaTotal,
					"Region": region,
				}).Info(msg("Lambda"))
			}
			totalLambdaNumber += rLambdaTotal
		}()

		// Getting CFN data
		wg.Add(1)
		go func() {
			defer wg.Done()
			rCfnTotal := listCfn(region, verbose)
			if rCfnTotal > 0 {
				logs.WithFields(logs.Fields{
					"CFN":    rCfnTotal,
					"Region": region,
				}).Info(msg("CFN"))
			}
			totalCfnNumber += rCfnTotal
		}()
	}
	// We do not care about region here, as we will get all
	wg.Add(1)
	go func() {
		defer wg.Done()
		s3Number, err = listS3("eu-west-1", verbose)
		if err != nil {
			log.Fatal("Cannot get S3 data: ", err)
		}
	}()

	wg.Wait()
	logs.WithFields(logs.Fields{
		"S3":                s3Number,
		"EC2":               totalEC2Number,
		"EC2Running":        totalEC2RunningNumber,
		"EC2RunningWindows": totalEC2RunningWindowsNumber,
		"ELB":               totalElbNumber,
		"ELBwithoutEC2":     totalElbWithoutEC2Number,
		"RDS":               totalRdsNumber,
		"RDS_Oracle":        totalOrRdsNumber,
		"RDS_MySQL":         totalMyRdsNumber,
		"RDS_MSSQL":         totalMsRdsNumber,
		"Lambda":            totalLambdaNumber,
		"CFN":               totalCfnNumber,
	}).Info("Account overview data")
}
