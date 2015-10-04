package filters
import (
	"fmt"
	"syscall"
	"strconv"
	"github.com/flosch/pongo2"
	_ "github.com/flosch/pongo2-addons"
)
func SuffixAssetsUpdate(in *pongo2.Value, param *pongo2.Value) (out *pongo2.Value, err *pongo2.Error) {
	assetsFile, ok := in.Interface().(string)
	if !ok {
		return nil, &pongo2.Error{
			Sender: "suffixStylesheet",
			ErrorMsg: fmt.Sprintf("Data must be string %T ('%v')", in, in),
		}
	}

	var file syscall.Stat_t
	syscall.Stat("./public/assets" + assetsFile, &file)
	timestamp, _ := file.Mtim.Unix()
	return pongo2.AsValue(assetsFile + "?update=" + strconv.FormatInt(timestamp, 10)), nil
}
