/*
Copyright 2023 The K8sGPT Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Some parts of this file have been modified to make it functional in Zadig

package analysis

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/segmentio/encoding/json"

	"github.com/magiconair/properties/assert"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

// sub-function
func analysis_RunAnalysisFilterTester(t *testing.T, filterFlag string, activceFilters []string) []Result {
	clientset := fake.NewSimpleClientset(
		&v1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "example",
				Namespace: "default",
			},
			Status: v1.PodStatus{
				Phase: v1.PodPending,
				Conditions: []v1.PodCondition{
					{
						Type:    v1.PodScheduled,
						Reason:  "Unschedulable",
						Message: "0/1 nodes are available: 1 node(s) had taint {node-role.kubernetes.io/master: }, that the pod didn't tolerate.",
					},
				},
			},
		},
		&v1.Endpoints{
			ObjectMeta: metav1.ObjectMeta{
				Name:        "example",
				Namespace:   "default",
				Annotations: map[string]string{},
			},
		},
		&v1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:        "example",
				Namespace:   "default",
				Annotations: map[string]string{},
			},
			Spec: v1.ServiceSpec{
				Selector: map[string]string{
					"app": "example",
				},
			},
		},
		&networkingv1.Ingress{
			ObjectMeta: metav1.ObjectMeta{
				Name:        "example",
				Namespace:   "default",
				Annotations: map[string]string{},
			},
		},
	)

	analysis := Analysis{
		Context:        context.Background(),
		Results:        []Result{},
		Namespace:      "default",
		MaxConcurrency: 1,
		Client: &Client{
			Client: clientset,
		},
	}
	if len(filterFlag) > 0 {
		// `--filter` is explicitly given
		analysis.Filters = strings.Split(filterFlag, ",")
	}
	analysis.RunAnalysis(activceFilters)
	return analysis.Results

}

// Test: Filter logic with running different Analyzers
func TestAnalysis_RunAnalysisWithFilter(t *testing.T) {
	var results []Result
	var filterFlag string

	//1. Neither --filter flag Nor active filter is specified, only the "core analyzers"
	results = analysis_RunAnalysisFilterTester(t, "", nil)
	assert.Equal(t, len(results), 3) // all built-in resource will be analyzed

	//2. When the --filter flag is specified

	filterFlag = "Pod" // --filter=Pod
	results = analysis_RunAnalysisFilterTester(t, filterFlag, nil)
	assert.Equal(t, len(results), 1)
	assert.Equal(t, results[0].Kind, filterFlag)

	filterFlag = "Ingress,Pod" // --filter=Ingress,Pod
	results = analysis_RunAnalysisFilterTester(t, filterFlag, nil)
	assert.Equal(t, len(results), 2)
}

// Test:  Filter logic with Active Filter
func TestAnalysis_RunAnalysisActiveFilter(t *testing.T) {

	//When the --filter flag is not specified but has actived filter in config
	var results []Result

	activeFilters := []string{"Ingress"}
	results = analysis_RunAnalysisFilterTester(t, "", activeFilters)
	assert.Equal(t, len(results), 1)

	activeFilters = []string{"Ingress", "Service"}
	results = analysis_RunAnalysisFilterTester(t, "", activeFilters)
	assert.Equal(t, len(results), 2)

	activeFilters = []string{"Ingress", "Service", "Pod"}
	results = analysis_RunAnalysisFilterTester(t, "", activeFilters)
	assert.Equal(t, len(results), 3)
}

func TestAnalysis_NoProblemJsonOutput(t *testing.T) {

	analysis := Analysis{
		Results:   []Result{},
		Namespace: "default",
	}

	expected := JsonOutput{
		Status:   StateOK,
		Problems: 0,
		Results:  []Result{},
	}

	gotJson, err := analysis.PrintOutput("json")
	if err != nil {
		t.Error(err)
	}

	got := JsonOutput{}
	err = json.Unmarshal(gotJson, &got)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(got)
	fmt.Println(expected)

	require.Equal(t, got, expected)

}

func TestAnalysis_ProblemJsonOutput(t *testing.T) {
	analysis := Analysis{
		Results: []Result{
			{
				Kind: "Deployment",
				Name: "test-deployment",
				Error: []Failure{
					{
						Text:      "test-problem",
						Sensitive: []Sensitive{},
					},
				},
				Details:      "test-solution",
				ParentObject: "parent-resource"},
		},
		Namespace: "default",
	}

	expected := JsonOutput{
		Status:   StateProblemDetected,
		Problems: 1,
		Results: []Result{
			{
				Kind: "Deployment",
				Name: "test-deployment",
				Error: []Failure{
					{
						Text:      "test-problem",
						Sensitive: []Sensitive{},
					},
				},
				Details:      "test-solution",
				ParentObject: "parent-resource"},
		},
	}

	gotJson, err := analysis.PrintOutput("json")
	if err != nil {
		t.Error(err)
	}

	got := JsonOutput{}
	err = json.Unmarshal(gotJson, &got)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(got)
	fmt.Println(expected)

	require.Equal(t, got, expected)
}

func TestAnalysis_MultipleProblemJsonOutput(t *testing.T) {
	analysis := Analysis{
		Results: []Result{
			{
				Kind: "Deployment",
				Name: "test-deployment",
				Error: []Failure{
					{
						Text:      "test-problem",
						Sensitive: []Sensitive{},
					},
					{
						Text:      "another-test-problem",
						Sensitive: []Sensitive{},
					},
				},
				Details:      "test-solution",
				ParentObject: "parent-resource"},
		},
		Namespace: "default",
	}

	expected := JsonOutput{
		Status:   StateProblemDetected,
		Problems: 2,
		Results: []Result{
			{
				Kind: "Deployment",
				Name: "test-deployment",
				Error: []Failure{
					{
						Text:      "test-problem",
						Sensitive: []Sensitive{},
					},
					{
						Text:      "another-test-problem",
						Sensitive: []Sensitive{},
					},
				},
				Details:      "test-solution",
				ParentObject: "parent-resource"},
		},
	}

	gotJson, err := analysis.PrintOutput("json")
	if err != nil {
		t.Error(err)
	}

	got := JsonOutput{}
	err = json.Unmarshal(gotJson, &got)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(got)
	fmt.Println(expected)

	require.Equal(t, got, expected)
}
