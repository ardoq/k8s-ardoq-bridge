package unit_test

import (
	"K8SArdoqBridge/app/subscriptions"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Resource Logic Tests", func() {
	type testCases []struct {
		namespaceLabels map[string]string
		resourceLabels  map[string]string
		expectation     bool
		description     string
	}

	cases := testCases{
		{
			namespaceLabels: map[string]string{
				"sync-to-ardoq": "enabled",
			},
			resourceLabels: map[string]string{},
			expectation:    true,
			description:    "tagged on the namespace level",
		}, {
			namespaceLabels: map[string]string{
				"sync-to-ardoq": "enabled",
			},
			resourceLabels: map[string]string{
				"sync-to-ardoq": "disabled",
			},
			expectation: false,
			description: "resource explicitly disabled",
		},
		{
			namespaceLabels: map[string]string{},
			resourceLabels:  map[string]string{},
			expectation:     false,
			description:     "neither tagged",
		},
		{
			namespaceLabels: map[string]string{},
			resourceLabels: map[string]string{
				"sync-to-ardoq": "enabled",
			},
			expectation: true,
			description: "tagged on the resource level",
		},
		{
			namespaceLabels: map[string]string{},
			resourceLabels: map[string]string{
				"sync-to-ardoq": "disabled",
			},
			expectation: false,
			description: "resource explicitly disabled on the resource level",
		},
	}
	Context("Resource sync candidate tests", func() {
		for _, x := range cases {
			It(x.description, func() {
				Expect(subscriptions.PerformSync(x.namespaceLabels, x.resourceLabels)).To(BeEquivalentTo(x.expectation))
			})
		}
	})
	Context("Resource cleanup candidates tests", func() {

		for _, x := range cases {
			It(x.description, func() {
				Expect(subscriptions.PerformSync(x.namespaceLabels, x.resourceLabels)).To(BeEquivalentTo(x.expectation))
			})
		}
	})
})
