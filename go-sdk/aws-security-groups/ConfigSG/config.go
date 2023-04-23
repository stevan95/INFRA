package ConfigSG

type SGDetails struct {
	VpcName                  string `json:"vpcName"`
	SgName                   string `json:"sgName"`
	SgDescription            string `json:"sgDescription"`
	SgPortsIPAddressProtocol string `json:"sgPortsIPAddressProtocol"`
}

var NewSG SGDetails
