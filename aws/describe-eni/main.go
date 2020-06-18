package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"os"
)

func main() {
	var err error
	regions, err := fetchRegion()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	vpcs := []string{}
	regionmap := make(map[string][]string)
	for _, region := range regions {
		vpcs, err = fetchVPC(region)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		//fmt.Printf("Region %v, VPCs: %v\n\n", region, vpcs)
		regionmap[region] = vpcs
		//create region/vpc map
	}
	r, err := fetchENI(regionmap)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(r)
}

func fetchRegion() ([]string, error) {
	awsSession := session.Must(session.NewSession(&aws.Config{Region: aws.String("eu-west-1")}))

	svc := ec2.New(awsSession)
	awsRegions, err := svc.DescribeRegions(&ec2.DescribeRegionsInput{})
	if err != nil {
		return nil, err
	}

	regions := make([]string, 0, len(awsRegions.Regions))
	for _, region := range awsRegions.Regions {
		regions = append(regions, *region.RegionName)
	}

	return regions, nil
}

func fetchVPC(region string) ([]string, error){
	awsSession := session.Must(session.NewSession(&aws.Config{Region: aws.String(region)}))

	svc := ec2.New(awsSession)
	awsVPCs, err := svc.DescribeVpcs(&ec2.DescribeVpcsInput{})
	if err != nil {
		return nil, err
	}

	vpcs := make([]string, 0, len(awsVPCs.Vpcs))
	for vpc := 0; vpc < len(awsVPCs.Vpcs) ; vpc++ {
		vpcs = append(vpcs, *awsVPCs.Vpcs[vpc].VpcId)
	}

	if err != nil {
		return nil, err
	}

	return vpcs, nil
}

func fetchENI(regionmap map[string][]string) ([]string, error){

	for k, _ := range regionmap {
		awsSession := session.Must(session.NewSession(&aws.Config{Region: aws.String(k)}))
		svc := ec2.New(awsSession)
		awsENIs, err := svc.DescribeNetworkInterfaces(&ec2.DescribeNetworkInterfacesInput{})
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Region: %v\n", k)
		for eni := range awsENIs.NetworkInterfaces{
			fmt.Println(string(*awsENIs.NetworkInterfaces[eni].NetworkInterfaceId))
		}
	}

	return nil, nil
}