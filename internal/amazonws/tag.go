package amazonws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type Tag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Tags []Tag

func (t Tags) Get(key string) string {
	for _, tag := range t {
		if tag.Key == key {
			return tag.Value
		}
	}
	return ""
}

func NewTagFromEc2(tag *ec2.Tag) Tag {
	t := Tag{
		Key:   aws.StringValue(tag.Key),
		Value: aws.StringValue(tag.Value),
	}

	return t
}

func NewTagsFromEc2(tags []*ec2.Tag) Tags {
	t := make(Tags, 0, 10)
	for _, tag := range tags {
		t = append(t, NewTagFromEc2(tag))
	}

	return t
}
