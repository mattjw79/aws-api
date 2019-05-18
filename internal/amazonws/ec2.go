package amazonws

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type Instance struct {
	AMI              string    `json:"ami"`
	AvailabilityZone string    `json:"availability-zone"`
	ID               string    `json:"id"`
	LaunchTime       time.Time `json:"launch-time"`
	Region           string    `json:"region"`
	State            string    `json:"state"`
	Tags             Tags      `json:"tags"`
	Type             string    `json:"type"`
	VPC              string    `json:"vpc"`
}

func NewInstanceFromEc2(instance *ec2.Instance, region string) Instance {
	i := Instance{
		AMI:              aws.StringValue(instance.ImageId),
		AvailabilityZone: aws.StringValue(instance.Placement.AvailabilityZone),
		ID:               aws.StringValue(instance.InstanceId),
		LaunchTime:       aws.TimeValue(instance.LaunchTime),
		Region:           region,
		State:            aws.StringValue(instance.State.Name),
		Tags:             NewTagsFromEc2(instance.Tags),
		Type:             aws.StringValue(instance.InstanceType),
		VPC:              aws.StringValue(instance.VpcId),
	}

	return i
}
