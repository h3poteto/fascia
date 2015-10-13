package controllers_test

import (
	"io/ioutil"
	"net/http"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/flosch/pongo2"

	"testing"

	"../filters"
)

func TestControllers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Controllers Suite")
}

var _ = BeforeSuite(func() {
	pongo2.RegisterFilter("suffixAssetsUpdate", filters.SuffixAssetsUpdate)
	pongo2.DefaultSet = pongo2.NewSet("test", pongo2.MustNewLocalFileSystemLoader("../views"))
})

func ParseResponse(res *http.Response) (string, int) {
	defer res.Body.Close()
	contents, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	return string(contents), res.StatusCode
}
