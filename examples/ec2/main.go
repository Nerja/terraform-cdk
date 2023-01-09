package main

import (
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdktf/cdktf-provider-aws-go/aws/v10/iaminstanceprofile"
	"github.com/cdktf/cdktf-provider-aws-go/aws/v10/iampolicy"
	"github.com/cdktf/cdktf-provider-aws-go/aws/v10/iamrole"
	"github.com/cdktf/cdktf-provider-aws-go/aws/v10/iamrolepolicyattachment"
	"github.com/cdktf/cdktf-provider-aws-go/aws/v10/instance"
	awsprovider "github.com/cdktf/cdktf-provider-aws-go/aws/v10/provider"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
)

func NewMyStack(scope constructs.Construct, id string) cdktf.TerraformStack {
	stack := cdktf.NewTerraformStack(scope, &id)

	awsprovider.NewAwsProvider(stack, jsii.String("AWS"), &awsprovider.AwsProviderConfig{
		Region: jsii.String("eu-west-1"),
	})

	policy := iampolicy.NewIamPolicy(stack, jsii.String("ec2-policy"), &iampolicy.IamPolicyConfig{
		Name: jsii.String("ec2-s3-policy"),
		Policy: jsii.String(`{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Effect": "Allow",
					"Action": [
						"s3:*"
					],
					"Resource": "*"
				}
			]
		}`),
	})

	role := iamrole.NewIamRole(stack, jsii.String("ec2-role"), &iamrole.IamRoleConfig{
		Name: jsii.String("ec2-s3-role"),
		AssumeRolePolicy: jsii.String(`{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Effect": "Allow",
					"Principal": {
						"Service": "ec2.amazonaws.com"
					},
					"Action": "sts:AssumeRole"
				}
			]
		}`),
	})

	iamrolepolicyattachment.NewIamRolePolicyAttachment(stack, jsii.String("ec2-policy-role-attachment"), &iamrolepolicyattachment.IamRolePolicyAttachmentConfig{
		PolicyArn: policy.Arn(),
		Role:      role.Name(),
	})

	instanceProfile := iaminstanceprofile.NewIamInstanceProfile(stack, jsii.String("ec2-instance-profile"), &iaminstanceprofile.IamInstanceProfileConfig{
		Name: jsii.String("ec2-s3-instance-profile"),
		Role: role.Name(),
	})

	instance := instance.NewInstance(stack, jsii.String("ec2-instance"), &instance.InstanceConfig{
		Ami:                jsii.String("ami-0abe92d15a280b758"),
		InstanceType:       jsii.String("t4g.nano"),
		IamInstanceProfile: instanceProfile.Id(),
		SubnetId:           jsii.String("subnet-0c8839e5d1a9a9b61"),
	})

	cdktf.NewTerraformOutput(stack, jsii.String("instance_id"), &cdktf.TerraformOutputConfig{
		Value: instance.Id(),
	})

	return stack
}

func main() {
	app := cdktf.NewApp(nil)

	stack := NewMyStack(app, "ec2")
	cdktf.NewCloudBackend(stack, &cdktf.CloudBackendProps{
		Hostname:     jsii.String("app.terraform.io"),
		Organization: jsii.String("testnerja"),
		Workspaces:   cdktf.NewNamedCloudWorkspace(jsii.String("ec2-elb-asg")),
	})

	app.Synth()
}
