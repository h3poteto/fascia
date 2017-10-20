package pongo2echo

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/flosch/pongo2"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

// Pongo2Echo implements custom pongo2 rendering engine for echo
type Pongo2Echo struct {
	dirs              []string
	filters           []string
	templates         *pongo2.TemplateSet
	contextProcessors []ContextProcessorFunc
}

// NewRenderer creates a new Pongo2Echo struct
func NewRenderer() *Pongo2Echo {
	p := &Pongo2Echo{}
	p.templates = pongo2.NewSet("templates", p)

	return p
}

// ContextProcessorFunc signature
type ContextProcessorFunc func(echoCtx echo.Context, pongoCtx pongo2.Context)

// UseContextProcessor adds context processor to the pipeline
func (p *Pongo2Echo) UseContextProcessor(processor ContextProcessorFunc) {
	p.contextProcessors = append(p.contextProcessors, processor)
}

// Abs returns absolute path to file requested
func (p *Pongo2Echo) Abs(base, name string) string {
	if filepath.IsAbs(name) {
		return name
	}

	for _, dir := range p.dirs {
		fullpath := filepath.Join(dir, name)
		_, err := os.Stat(fullpath)
		if err == nil {
			return fullpath
		}
	}

	return filepath.Join(filepath.Dir(base), name)
}

// Get reads the path's content from your local filesystem.
func (p *Pongo2Echo) Get(path string) (io.Reader, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(buf), nil
}

// AddDirectory adds a directory to the list of directories searched for templates
func (p *Pongo2Echo) AddDirectory(dir string) {
	p.dirs = append(p.dirs, dir)
}

// RegisterTag registers a custom tag
func (p *Pongo2Echo) RegisterTag(name string, parserFunc pongo2.TagParser) {
	pongo2.RegisterTag(name, parserFunc)
}

// RegisterFilter registers a custom filter
func (p *Pongo2Echo) RegisterFilter(name string, fn pongo2.FilterFunction) {
	pongo2.RegisterFilter(name, fn)
}

// Render renders the view
func (p *Pongo2Echo) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	tmpl, err := p.templates.FromCache(name)
	if err != nil {
		return err
	}
	d, ok := data.(map[string]interface{})
	if !ok {
		return errors.New("Incorrect data format. Should be map[string]interface{}")
	}

	// run context processors
	for _, processor := range p.contextProcessors {
		processor(c, d)
	}

	return tmpl.ExecuteWriter(pongo2.Context(d), w)
}
