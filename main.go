package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

const region = "ap-northeast-1"

func ec2StartInstances() {
	ec2InsIds := strings.Split(os.Getenv("EC2INS_IDS"), ",")
	if len(ec2InsIds) == 0 {
		fmt.Println("指定されたインスタンスはありません")
	} else {
		ec2Client := getEc2Client()
		insStatusMap := getEc2InstanceStatuses(ec2Client)
		for _, insId := range ec2InsIds {
			if insStatus, ok := insStatusMap[insId]; ok {
				if insStatus == "stopped" {
					startInstance(ec2Client, insId)
				} else {
					fmt.Printf("インスタンス[%s]は既に起動しています\n", insId)
				}
			} else {
				fmt.Printf("インスタンス[%s]は存在しません\n", insId)
			}
		}
	}
	return
}

func main() {
	lambda.Start(ec2StartInstances)
}

// EC2インスタンスを生成する
func getEc2Client() *ec2.EC2 {
	var config aws.Config
	config = aws.Config{Region: aws.String(region)}
	sess := session.New(&config)
	svc := ec2.New(sess)
	return svc
}

// InstanceIdに紐づくInstanceStateをmap型に格納して返す
func getEc2InstanceStatuses(ec2Client *ec2.EC2) map[string]string {
	resp, err := ec2Client.DescribeInstances(nil)
	if err != nil {
		panic(err)
	}
	instanceStates := map[string]string{}
	for _, r := range resp.Reservations {
		for _, i := range r.Instances {
			instanceStates[*i.InstanceId] = *i.State.Name
		}
	}
	return instanceStates
}

// 指定したInstanceIdを起動する
func startInstance(ec2Client *ec2.EC2, insId string) {
	insIds := []*string{aws.String(insId)}
	_, err := ec2Client.StartInstances(&ec2.StartInstancesInput{
		InstanceIds: insIds,
	})
	if err != nil {
		panic(err)
	}
}
