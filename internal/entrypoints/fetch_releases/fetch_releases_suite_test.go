package fetchreleases_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestFetchReleases(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "FetchReleases Suite")
}
