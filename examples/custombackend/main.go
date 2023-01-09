package main

import (
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	awsprovider "github.com/cdktf/cdktf-provider-aws-go/aws/v10/provider"
	"github.com/cdktf/cdktf-provider-aws-go/aws/v10/s3bucket"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
)

func NewMyStack(scope constructs.Construct, id string) cdktf.TerraformStack {
	stack := cdktf.NewTerraformStack(scope, &id)

	bucket := s3bucket.NewS3Bucket(stack, jsii.String("example-s3"), &s3bucket.S3BucketConfig{
		Bucket: jsii.String("ntcdktest"),
		Tags: &map[string]*string{
			"SomeChangedTag": jsii.String("SomeChangedValue"),
		},
	})

	cdktf.NewTerraformOutput(stack, jsii.String("bucket_name"), &cdktf.TerraformOutputConfig{
		Value: bucket.Id(),
	})

	return stack
}

func main() {
	app := cdktf.NewApp(nil)

	stack := NewMyStack(app, "custombackend")

	awsprovider.NewAwsProvider(stack, jsii.String("AWS"), &awsprovider.AwsProviderConfig{
		Region: jsii.String("eu-west-1"),
	})

	cdktf.NewS3Backend(stack, &cdktf.S3BackendProps{
		Bucket: jsii.String("nerja-terraform-cdk-test-backend"),
		Key:    jsii.String("lab/custombackend/terraform.tfstate"),
		Region: jsii.String("eu-west-1"),
	})

	app.Synth()
}
