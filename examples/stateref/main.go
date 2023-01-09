package main

import (
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	awsprovider "github.com/cdktf/cdktf-provider-aws-go/aws/v10/provider"
	"github.com/cdktf/cdktf-provider-aws-go/aws/v10/s3bucket"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
)

var ec2Id string = "instance_id"

func NewMyStack(scope constructs.Construct, id string) cdktf.TerraformStack {
	stack := cdktf.NewTerraformStack(scope, &id)

	awsprovider.NewAwsProvider(stack, jsii.String("AWS"), &awsprovider.AwsProviderConfig{
		Region: jsii.String("eu-west-1"),
	})

	state := cdktf.NewDataTerraformRemoteState(stack, &id, &cdktf.DataTerraformRemoteStateRemoteConfig{
		Hostname:     jsii.String("app.terraform.io"),
		Organization: jsii.String("testnerja"),
		Workspaces:   cdktf.NewNamedCloudWorkspace(jsii.String("ec2-elb-asg")),
	})

	bucket := s3bucket.NewS3Bucket(stack, jsii.String("some-bucket"), &s3bucket.S3BucketConfig{
		Bucket: jsii.String("nerja-ref-output-test"),
		Tags: &map[string]*string{
			"SomeTag": state.GetString(&ec2Id),
		},
	})

	cdktf.NewTerraformOutput(stack, jsii.String("bucket-output"), &cdktf.TerraformOutputConfig{
		Value: bucket.Id(),
	})

	return stack
}

func main() {
	app := cdktf.NewApp(nil)

	stack := NewMyStack(app, "stateref")
	cdktf.NewCloudBackend(stack, &cdktf.CloudBackendProps{
		Hostname:     jsii.String("app.terraform.io"),
		Organization: jsii.String("testnerja"),
		Workspaces:   cdktf.NewNamedCloudWorkspace(jsii.String("stateref")),
	})

	app.Synth()
}
