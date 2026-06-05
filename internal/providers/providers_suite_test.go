package providers_test

import (
	"testing"

	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestProviders(t *testing.T) {
	httpmock.Activate(t)

	RegisterFailHandler(Fail)
	RunSpecs(t, "Providers Suite")
}
