package passwordhashing_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPasswordHashing(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "PasswordHashing Suite")
}
