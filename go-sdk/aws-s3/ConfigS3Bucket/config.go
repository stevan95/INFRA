package ConfigS3Bucket

type S3Config struct {
	BucketName     string
	RegionName     string
	PathToUpload   string
	FileToDownload string
}

var S3Conf S3Config

func (c *S3Config) InitS3Config(BucketName, RegionName string, PathToUpload, FileToDownload *string) {
	c.BucketName = BucketName
	c.RegionName = RegionName
	if PathToUpload != nil {
		//Optional parameter is provided
		c.PathToUpload = *PathToUpload
	}

	if FileToDownload != nil {
		//Optional parameter is provided
		c.FileToDownload = *FileToDownload
	}
}
