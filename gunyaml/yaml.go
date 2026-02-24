package gunyaml

import (
	"io/fs"

	"sigs.k8s.io/yaml"
)

// Decoder of YAML files for aconfig.
type Decoder struct {
	fsys fs.FS
}

// New YAML decoder for aconfig.
func New() *Decoder { return &Decoder{} }

// Format of the decoder.
func (d *Decoder) Format() string {
	return "yaml"
}

// DecodeFile implements aconfig.FileDecoder.
func (d *Decoder) DecodeFile(filename string) (map[string]interface{}, error) {
	b, err := fs.ReadFile(d.fsys, filename)
	if err != nil {
		return nil, err
	}
	var raw map[string]interface{}
	if err := yaml.Unmarshal(b, &raw); err != nil {
		return nil, err
	}
	return raw, nil
}

// Init implements aconfig.FileDecoder.
func (d *Decoder) Init(fsys fs.FS) {
	d.fsys = fsys
}
