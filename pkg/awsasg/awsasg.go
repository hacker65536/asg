package awsasg

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling/types"
	"github.com/hacker65536/asg/pkg/utils"
)

type AWSAsg struct {
	svc *autoscaling.Client
}

type Asg struct {
	Max     int32    `json:"max"`
	Min     int32    `json:"min"`
	Name    string   `json:"name"`
	Tg      []string `json:"tg"`
	Lb      []string `json:"lb"`
	Ec2s    []Ec2    `json:"ec2s"`
	Desired int32    `json:"desired"`
	//Tags    []map[string]interface{}    `json:"desired"`
}

type Ec2 struct {
	InstanceId   string
	HealthStatus string
	InstanceType string
}

func New() *AWSAsg {

	cfg, err := config.LoadDefaultConfig()
	if err != nil {
		panic("unable to load SDK config, " + err.Error())
	}
	svc := autoscaling.NewFromConfig(cfg)
	return &AWSAsg{
		svc: svc,
	}
}

func (a *AWSAsg) Ls(args []string) {
	params := &autoscaling.DescribeAutoScalingGroupsInput{}
	if len(args) != 0 {
		params.AutoScalingGroupNames = args
	}
	//	LsOutputFull(a.describeAsg(params))
	a.LsNameOnly()
}
func (a *AWSAsg) LsNameOnly() {

	svc := a.svc
	params := &autoscaling.DescribeAutoScalingGroupsInput{}
	p := autoscaling.NewDescribeAutoScalingGroupsPaginator(svc, params)
	var i int
	for p.HasMorePages() {
		i++
		page, err := p.NextPage(context.TODO())

		if err != nil {
			log.Fatal("anything is wrong")
		}
		for _, v := range page.AutoScalingGroups {
			fmt.Println(aws.ToString(v.AutoScalingGroupName))
		}
	}
}

func LsOutputFull(asgs []Asg) {

	for _, v := range asgs {
		w := tabwriter.NewWriter(os.Stdout, 0, 1, 1, ' ', tabwriter.DiscardEmptyColumns)
		fmt.Fprintln(w, strings.Join([]string{
			utils.Normal(v.Name),
			utils.Normal("num"),
		}, " "))

		// title row
		fmt.Fprintln(w, strings.Join([]string{
			utils.Normal("status"),
			utils.Normal("id"),
			utils.Normal("class"),
		}, "\t"))

		// instance data rows
		if len(v.Ec2s) > 0 {
			for _, v2 := range v.Ec2s {
				fmt.Fprintln(w, strings.Join([]string{
					func() string {
						if v2.HealthStatus == "Healthy" {
							return utils.Green(v2.HealthStatus)
						}
						return utils.Yellow(v2.HealthStatus)
					}(),
					utils.Normal(v2.InstanceId),
					utils.Normal(v2.InstanceType),
				}, "\t"))
			}
		}
		w.Flush()
	}

}

func (a *AWSAsg) describeAsg(params *autoscaling.DescribeAutoScalingGroupsInput) []Asg {

	svc := a.svc
	p := autoscaling.NewDescribeAutoScalingGroupsPaginator(svc, params)

	asgs := []Asg{}
	var i int
	for p.HasMorePages() {
		i++

		page, err := p.NextPage(context.TODO())

		if err != nil {
			log.Fatal("anything is wrong")
		}

		/*
			j, err := json.Marshal(page)
			if err != nil {
				log.Fatal("anything is wrong")

			}
			fmt.Println(string(j))
		*/

		for _, v := range page.AutoScalingGroups {

			asgs = append(asgs, Asg{
				Max:     aws.ToInt32(v.MaxSize),
				Min:     aws.ToInt32(v.MinSize),
				Name:    aws.ToString(v.AutoScalingGroupName),
				Desired: aws.ToInt32(v.DesiredCapacity),
				Ec2s: func(is []types.Instance) []Ec2 {
					ec2s := []Ec2{}
					for _, v := range is {
						ec2s = append(ec2s, Ec2{
							InstanceId:   aws.ToString(v.InstanceId),
							HealthStatus: aws.ToString(v.HealthStatus),
							InstanceType: aws.ToString(v.InstanceType),
						})
					}
					return ec2s
				}(v.Instances),
			})
		}

	}
	return asgs
}
