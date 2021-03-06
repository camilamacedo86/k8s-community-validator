// Copyright 2021 The K8s Community Validator Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"errors"
	"fmt"
	"github.com/k8s-operatorhub/k8s-community-validator/pkg/validation"
	"os"

	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"

	"github.com/k8s-operatorhub/k8s-community-validator/pkg/result"
	apimanifests "github.com/operator-framework/api/pkg/manifests"
	apierrors "github.com/operator-framework/api/pkg/validation/errors"
)

func main() {

	var outputFormat string

	flag.StringVarP(&outputFormat, "output", "o", result.Text,
		"Result format for results. One of: [text, json-alpha1]. Note: output format types containing "+
			"\"alphaX\" are subject to change and not covered by guarantees of stable APIs.")

	flag.Parse()

	validate(outputFormat)
	results := runValidator()
	printResults(results, outputFormat)
}

func printResults(results []apierrors.ManifestResult, outputFormat string) {
	// Create Result to be output.
	res := result.NewResult()
	res.AddManifestResults(results...)

	if err := res.PrintWithFormat(outputFormat); err != nil {
		log.Fatal(err)
	}
}

func runValidator() []apierrors.ManifestResult {
	// Read the bundle
	bundle, err := apimanifests.GetBundleFromDir(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	objs := bundle.ObjectsToValidate()
	for _, obj := range bundle.Objects {
		objs = append(objs, obj)
	}

	// pass the objects to the validator
	results := validation.K8sCommunityValidator.Validate(objs...)
	return results
}

func validate(outputFormat string) {
	if len(os.Args) < 2 {
		log.Fatal(errors.New("an image tag or directory is a required argument"))
	}
	if outputFormat != result.JSONAlpha1 && outputFormat != result.Text {
		log.Fatal(fmt.Errorf("invalid value for output flag: %v", outputFormat))
	}
}