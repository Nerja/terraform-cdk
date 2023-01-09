package main

import (
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/hashicorp/terraform-cdk-go/cdktf"

	awsprovider "github.com/cdktf/cdktf-provider-aws-go/aws/v10/provider"
	"github.com/cdktf/cdktf-provider-aws-go/aws/v10/s3bucket"
)

func NewMyStack(scope constructs.Construct, id string) cdktf.TerraformStack {
	stack := cdktf.NewTerraformStack(scope, &id)

	awsprovider.NewAwsProvider(stack, jsii.String("AWS"), &awsprovider.AwsProviderConfig{
		Region: jsii.String("eu-west-1"),
	})

	bucket := s3bucket.NewS3Bucket(stack, jsii.String("my-bucket"), &s3bucket.S3BucketConfig{
		Bucket: jsii.String("terraform-cdk-nerja"),
	})

	cdktf.NewTerraformOutput(stack, jsii.String("bucket_name"), &cdktf.TerraformOutputConfig{
		Value: bucket.Bucket(),
	})

	return stack
}

func main() {
	app := cdktf.NewApp(nil)

	stack := NewMyStack(app, "terraform-cdk")
	cdktf.NewCloudBackend(stack, &cdktf.CloudBackendProps{
		Hostname:     jsii.String("app.terraform.io"),
		Organization: jsii.String("testnerja"),
		Workspaces:   cdktf.NewNamedCloudWorkspace(jsii.String("terraform-cdk")),
	})

	app.Synth()
}
