package main

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-cdk-go/cdktf"
)

func TestShouldContainContainer(t *testing.T) {
	stackRef := NewMyStack(cdktf.Testing_App(nil), "stack")

	synth := cdktf.Testing_FullSynth(stackRef)
	assertion := cdktf.Testing_ToBeValidTerraform(synth)
	if !*assertion {
		t.Error("Not valid Terraform")
	}

	expectedProps := map[string]interface{}{
		"bucket": "terraform-cdk-nerja",
	}
	runValidations := true
	json := cdktf.Testing_Synth(stackRef, &runValidations)
	resourceType := "aws_s3_bucket"
	fmt.Println(*json)
	if !*cdktf.Testing_ToHaveResourceWithProperties(json, &resourceType, &expectedProps) {
		t.Error("Failed to find resource with properties")
	}
}
