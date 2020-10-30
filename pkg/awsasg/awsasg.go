package awsasg

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling/types"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/hacker65536/ec2/pkg/awsec2"
	"github.com/logrusorgru/aurora"
	color "github.com/logrusorgru/aurora"
)

type AwsAsg struct {
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
}

type Ec2 struct {
	Id        string `json:"id"`
	Status    string `json:"status"`
	Class     string `json:"class"`
	LifeCycle string `json:"lifecycle"`
}

type Asgs []Asg
type Ec2s []Ec2

// New is initial package
func New() *AwsAsg {

	cfg, err := config.LoadDefaultConfig()
	if err != nil {
		panic("unable to load SDK config, " + err.Error())
	}
	svc := autoscaling.NewFromConfig(cfg)
	return &AwsAsg{
		svc: svc,
	}
}

// LS is listasg
func (a *AwsAsg) Ls(filter string) Asgs {

	reg := regexp.MustCompile(filter)
	r := a.describeAsg()
	asglist := Asgs{}

	for _, v := range r.AutoScalingGroups {

		if reg.MatchString(aws.ToString(v.AutoScalingGroupName)) {
			es := Ec2s{}
			for _, v2 := range v.Instances {
				e := Ec2{
					Id:        *v2.InstanceId,
					Status:    *v2.HealthStatus,
					Class:     *v2.InstanceType,
					LifeCycle: string(v2.LifecycleState),
				}
				es = append(es, e)
			}

			asg := Asg{
				Max:     aws.ToInt32(v.MaxSize),
				Min:     aws.ToInt32(v.MinSize),
				Name:    *v.AutoScalingGroupName,
				Tg:      aws.ToStringSlice(v.TargetGroupARNs),
				Lb:      aws.ToStringSlice(v.LoadBalancerNames),
				Ec2s:    es,
				Desired: aws.ToInt32(v.DesiredCapacity),
			}
			asglist = append(asglist, asg)
		}
	}

	/*
		j, _ := json.Marshal(asglist)
		fmt.Println(string(j))

		//fmt.Printf("%#v", asglist)
	*/
	return asglist
}

func (a *AwsAsg) LsOutput(f string) {
	l := a.Ls(f)
	for _, v := range l {

		fmt.Printf("%s [%d/%d]\n", color.Bold(v.Name), len(v.Ec2s), v.Desired)
		if len(v.Ec2s) > 0 {
			r := awsec2.New().Ls(
				&ec2.DescribeInstancesInput{
					InstanceIds: func() []*string {
						ids := []string{}
						for _, vv := range v.Ec2s {
							if vv.LifeCycle == "InService" {
								ids = append(ids, vv.Id)
							}
						}
						return aws.StringSlice(ids)
					}(),
				},
			)

			r2 := map[string]*time.Time{}

			for _, vv := range r {
				r2[vv.Id] = vv.LaunchTime
			}
			jst, _ := time.LoadLocation("Asia/Tokyo")

			fmt.Println("status\t\tid\t\t\tclass\t\tcreated\t\t\tuptime")
			for _, v2 := range v.Ec2s {
				fmt.Printf("%s%s\t%s\t%s\t[%s]\t%s\n",
					func() aurora.Value {
						if v2.Status == "Healthy" {
							return color.Green(v2.Status)
						}
						return color.Yellow(v2.Status)
					}(),
					func() string {
						if v2.LifeCycle != "InService" {
							return "*"
						}
						return " "
					}(),
					v2.Id,
					v2.Class,
					func() string {
						if v2.LifeCycle == "InService" {
							return r2[v2.Id].In(jst).Format("2006-01-02 15:04:05")
						}
						return ""
					}(),

					func() string {
						if v2.LifeCycle == "InService" {
							return time.Since(aws.ToTime(r2[v2.Id])).Round(60 * time.Second).String()
						}
						return ""
					}(),
				)
			}
		}
		fmt.Println()
		//		fmt.Println(r)
	}

}

func (a *AwsAsg) describeAsg() *autoscaling.DescribeAutoScalingGroupsOutput {
	svc := a.svc
	params := &autoscaling.DescribeAutoScalingGroupsInput{}
	maxKeys := 100
	p := NewASGDescribeAutoScalingGroupsPaginator(svc, params, func(o *ASGDescribeAutoScalingGroupsPaginatorOptions) {
		if v := int32(maxKeys); v != 0 {
			o.Limit = v
		}
	})

	r := &autoscaling.DescribeAutoScalingGroupsOutput{
		AutoScalingGroups: []*types.AutoScalingGroup{},
	}
	var i int
	for p.HasMorePages() {
		i++
		page, err := p.NextPage((context.TODO()))
		if err != nil {
			log.Fatalf("aaaaa")
		}
		/*
			for _, obj := range page.AutoScalingGroups {
				fmt.Println("asg:=", aws.ToString(obj.AutoScalingGroupName))
			}
		*/

		r.AutoScalingGroups = append(r.AutoScalingGroups, page.AutoScalingGroups...)
	}

	return r
}
