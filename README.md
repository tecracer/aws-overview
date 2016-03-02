# aws-overview

``aws-overview`` is a small script to get AWS account overview.
Available details are (number of):
* EC2, EC2 Running instances, EC2 Windows running instances
* ELB, ELB without assigned EC2 instances
* RDS, RDS MySQL/MSSQL/Oracle
* CFN
* Lambda functions
* S3 buckets
* Total of all above in all regions

## Installing

* ``go get github.com/partamonov/aws-overview``
* ``go install github.com/partamonov/aws-overview``

For cross platform compilation:
* ``env GOOS=windows GOARCH=amd64 go build``

## Usage

AWS Credentials expected in ``$HOME/.aws/credentials`` or as environment variables
``AWS_ACCESS_KEY_ID``
``AWS_SECRET_ACCESS_KEY``

```
Usage: aws-overview [-h] [-log-file=path] [other options]
 -h, --help
 -log-file=<PATH>: Log file location, if skipped logs to STDOUT
 -verbose=true/false: [bool], if true prints details information about objects
 -machine-readable=true/false: [bool], if true convert output to Logstash format, false print json output
```
