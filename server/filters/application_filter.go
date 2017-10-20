package filters

import (
	"fmt"
	"github.com/flosch/pongo2"
	_ "github.com/flosch/pongo2-addons"
	_ "github.com/russross/blackfriday" // pongo2-addonsが依存するblackfridayが古いため明示的な依存を書く必要がある
	"strconv"
	"syscall"
)

func SuffixAssetsUpdate(in *pongo2.Value, param *pongo2.Value) (out *pongo2.Value, err *pongo2.Error) {
	assetsFile, ok := in.Interface().(string)
	if !ok {
		return nil, &pongo2.Error{
			Sender:    "suffixStylesheet",
			OrigError: fmt.Errorf("Data must be string %T ('%v')", in, in),
		}
	}

	var file syscall.Stat_t
	syscall.Stat("./public/assets"+assetsFile, &file)
	timestamp, _ := file.Mtim.Unix()
	return pongo2.AsValue(assetsFile + "?update=" + strconv.FormatInt(timestamp, 10)), nil
}
