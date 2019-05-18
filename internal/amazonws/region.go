package amazonws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type Region struct {
	Name     string `json:"name"`
	Endpoint string `json:"endpoint"`
}

func NewRegion(name string) Region {
	r := Region{
		Name: name,
	}

	return r
}

func NewRegionFromEc2(region *ec2.Region) Region {
	r := Region{
		Name:     aws.StringValue(region.RegionName),
		Endpoint: aws.StringValue(region.Endpoint),
	}

	return r
}

func (r *Region) getInstancesByInput(input *ec2.DescribeInstancesInput) ([]Instance, error) {
	sess, err := session.NewSession(&aws.Config{Region: aws.String(r.Name)})
	if err != nil {
		return nil, err
	}
	svc := ec2.New(sess)

	result, err := svc.DescribeInstances(input)
	if err != nil {
		return nil, err
	}

	instances := make([]Instance, 0, 10)
	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			instances = append(instances, NewInstanceFromEc2(instance, r.Name))
		}
	}

	return instances, nil
}

func (r *Region) Instances() ([]Instance, error) {
	input := &ec2.DescribeInstancesInput{}
	return r.getInstancesByInput(input)
}

func (r *Region) InstanceByID(id string) ([]Instance, error) {
	input := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{
			aws.String(id),
		},
	}
	return r.getInstancesByInput(input)
}

func (r *Region) InstanceByName(name string) ([]Instance, error) {
	input := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name:   aws.String("tag:Name"),
				Values: []*string{aws.String(name)},
			},
		},
	}
	return r.getInstancesByInput(input)
}
