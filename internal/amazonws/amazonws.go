package amazonws

import (
	"errors"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func Regions() ([]Region, error) {
	sess, err := session.NewSession(&aws.Config{Region: aws.String("us-east-1")})
	svc := ec2.New(sess)
	input := &ec2.DescribeRegionsInput{}
	result, err := svc.DescribeRegions(input)
	if err != nil {
		return nil, err
	}

	regions := make([]Region, len(result.Regions))
	for idx, region := range result.Regions {
		regions[idx] = NewRegionFromEc2(region)
	}

	return regions, nil
}

func FindInstanceByID(id string) ([]Instance, error) {
	var wgPub, wgSub sync.WaitGroup
	c := make(chan Instance, 10)

	regions, err := Regions()
	if err != nil {
		return nil, err
	}

	for _, region := range regions {
		wgPub.Add(1)
		go func(r Region) {
			defer wgPub.Done()
			instances, _ := r.InstanceByID(id)
			for _, instance := range instances {
				c <- instance
			}
		}(region)
	}

	wgSub.Add(1)
	instances := make([]Instance, 0, 1)
	go func() {
		defer wgSub.Done()
		for instance := range c {
			instances = append(instances, instance)
		}
	}()

	wgPub.Wait()
	close(c)

	wgSub.Wait()

	if len(instances) == 0 {
		return nil, errors.New("instance not found")
	}

	return instances, nil
}

func FindInstanceByName(name string) ([]Instance, error) {
	var wgPub, wgSub sync.WaitGroup
	c := make(chan Instance, 10)

	regions, err := Regions()
	if err != nil {
		return nil, err
	}

	for _, region := range regions {
		wgPub.Add(1)
		go func(r Region) {
			defer wgPub.Done()
			instances, _ := r.InstanceByName(name)
			for _, instance := range instances {
				c <- instance
			}
		}(region)
	}

	wgSub.Add(1)
	instances := make([]Instance, 0, 1)
	go func() {
		defer wgSub.Done()
		for instance := range c {
			instances = append(instances, instance)
		}
	}()

	wgPub.Wait()
	close(c)

	wgSub.Wait()

	if len(instances) == 0 {
		return nil, errors.New("instance not found")
	}

	return instances, nil
}
