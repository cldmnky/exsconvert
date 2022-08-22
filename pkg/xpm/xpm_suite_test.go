package xpm

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestXPM(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "XPM Suite")
}
