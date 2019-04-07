package moleculer_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSimple(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Simple Suite")
}
