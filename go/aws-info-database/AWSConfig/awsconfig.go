package awsconfig

type EC2Information struct {
	InstanceID       string
	PrivateIPAddress string
	PrivateDnsName   string
	SubnetId         string
	PublicIpAddress  string
	InstanceType     string
	Tags             map[string]string
}

type EC2Instances struct {
	EC2Array []EC2Information
}
