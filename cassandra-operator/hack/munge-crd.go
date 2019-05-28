package main

import (
	"fmt"
	"io/ioutil"
	"sigs.k8s.io/yaml"

	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
)

func main() {
	bytes, err := ioutil.ReadFile("kubernetes-resources/cassandra-operator-crd.yml")
	if err != nil {
		panic(err.Error())
	}

	var crd apiextensionsv1beta1.CustomResourceDefinition
	err = yaml.Unmarshal(bytes, &crd)
	if err != nil {
		panic(err.Error())
	}

	crd.Spec.Scope = "Namespaced"
	crd.Spec.Version = "v1alpha1"

	// fix spec.validation.openAPIV3Schema.properties[metadata].properties[annotations].additionalProperties: Forbidden: additionalProperties cannot be set to false
	// fix spec.validation.openAPIV3Schema.properties[metadata].properties[labels].additionalProperties: Forbidden: additionalProperties cannot be set to false
	schemaProps := crd.Spec.Validation.OpenAPIV3Schema.Properties["metadata"]
	annotations := schemaProps.Properties["annotations"]
	annotations.AdditionalProperties = nil
	schemaProps.Properties["annotations"] = annotations
	labels := schemaProps.Properties["labels"]
	labels.AdditionalProperties = nil
	schemaProps.Properties["labels"] = labels
	crd.Spec.Validation.OpenAPIV3Schema.Properties["metadata"] = schemaProps

	y, err := yaml.Marshal(crd)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	fmt.Println(string(y))
}
