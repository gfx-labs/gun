package gun

import (
	"fmt"
	"os"
	"path"

	"github.com/cristalhq/aconfig"
	"github.com/gfx-labs/gun/gunyaml"
	"github.com/joho/godotenv"
)

func Load(i any) {
	LoadPrefix(i, "")
}

func loadPrefix(i any, prefix string) error {
	yamlDecoder := gunyaml.New()
	fileName := "config"
	if prefix != "" {
		fileName = prefix
	}
	homeDir, _ := os.UserHomeDir()
	godotenv.Load()
	loader := aconfig.LoaderFor(i, aconfig.Config{
		AllowUnknownFields: true,
		AllowUnknownEnvs:   true,
		AllowUnknownFlags:  true,
		SkipFlags:          true,
		SkipEnv:            false,
		DontGenerateTags:   true,
		MergeFiles:         true,
		EnvPrefix:          prefix,
		FlagPrefix:         prefix,
		Files: []string{
			fmt.Sprintf("/%s.yml", fileName),
			fmt.Sprintf("/%s.yaml", fileName),
			fmt.Sprintf("/%s.json", fileName),
			fmt.Sprintf("/config/%s.yml", fileName),
			fmt.Sprintf("/config/%s.yaml", fileName),
			fmt.Sprintf("/config/%s.json", fileName),
			path.Join(homeDir, fmt.Sprintf(".gun/%s.yml", fileName)),
			path.Join(homeDir, fmt.Sprintf(".gun/%s.yaml", fileName)),
			path.Join(homeDir, fmt.Sprintf(".gun/%s.json", fileName)),
			fmt.Sprintf("./%s.yml", fileName),
			fmt.Sprintf("./%s.yaml", fileName),
			fmt.Sprintf("./%s.json", fileName),
		},
		FileDecoders: map[string]aconfig.FileDecoder{
			".yaml": yamlDecoder,
			".yml":  yamlDecoder,
			".json": yamlDecoder,
		},
	})
	if err := loader.Load(); err != nil {
		return err
	}
	return nil

}

func LoadPrefix(i any, prefix string) {
	err := loadPrefix(i, prefix)
	if err != nil {
		panic(err)
	}
}
