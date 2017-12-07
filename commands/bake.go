package commands

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	yaml "gopkg.in/yaml.v2"

	"github.com/pivotal-cf/jhanda/commands"
	"github.com/pivotal-cf/jhanda/flags"
	"github.com/pivotal-cf/kiln/builder"
)

type BakeConfig struct {
	EmbedPaths               flags.StringSlice `short:"e"    long:"embed"                      description:"path to files to include in the tile /embed directory"`
	FormDirectories          flags.StringSlice `short:"f"    long:"forms-directory"            description:"path to a directory containing forms"`
	IconPath                 string            `short:"i"    long:"icon"                       description:"path to icon file"`
	InstanceGroupDirectories flags.StringSlice `short:"ig"   long:"instance-groups-directory"  description:"path to a directory containing instance groups"`
	JobDirectories           flags.StringSlice `short:"j"    long:"jobs-directory"             description:"path to a directory containing jobs"`
	Metadata                 string            `short:"m"    long:"metadata"                   description:"path to the metadata file"`
	MigrationDirectories     flags.StringSlice `short:"md"   long:"migrations-directory"       description:"path to a directory containing migrations"`
	OutputFile               string            `short:"o"    long:"output-file"                description:"path to where the tile will be output"`
	PropertyDirectories      flags.StringSlice `short:"pd"   long:"properties-directory"       description:"path to a directory containing property blueprints"`
	ReleaseDirectories       flags.StringSlice `short:"rd"   long:"releases-directory"         description:"path to a directory containing release tarballs"`
	RuntimeConfigDirectories flags.StringSlice `short:"rcd"  long:"runtime-configs-directory"  description:"path to a directory containing runtime configs"`
	StemcellTarball          string            `short:"st"   long:"stemcell-tarball"           description:"path to a stemcell tarball"`
	StubReleases             bool              `short:"sr"   long:"stub-releases"              description:"skips importing release tarballs into the tile"`
	Variables                flags.StringSlice `short:"vr"   long:"variable"                   description:"key value pairs of variables to interpolate"`
	VariableDirectories      flags.StringSlice `short:"vd"   long:"variables-directory"        description:"path to a directory containing variables"`
	Version                  string            `short:"v"    long:"version"                    description:"version of the tile"`
}

type Bake struct {
	metadataBuilder metadataBuilder
	tileWriter      tileWriter
	logger          logger
	Options         BakeConfig
}

//go:generate counterfeiter -o ./fakes/tile_writer.go --fake-name TileWriter . tileWriter

type tileWriter interface {
	Write(productName string, generatedMetadataContents []byte, input builder.WriteInput) error
}

//go:generate counterfeiter -o ./fakes/metadata_builder.go --fake-name MetadataBuilder . metadataBuilder

type metadataBuilder interface {
	Build(input builder.BuildInput) (builder.GeneratedMetadata, error)
}

//go:generate counterfeiter -o ./fakes/logger.go --fake-name Logger . logger

type logger interface {
	Printf(format string, v ...interface{})
	Println(v ...interface{})
}

func NewBake(metadataBuilder metadataBuilder, tileWriter tileWriter, logger logger) Bake {
	return Bake{
		metadataBuilder: metadataBuilder,
		tileWriter:      tileWriter,
		logger:          logger,
	}
}

func (b Bake) Execute(args []string) error {
	config, err := b.parseArgs(args)
	if err != nil {
		return err
	}

	releaseTarballs, err := b.extractReleaseTarballFilenames(config)
	if err != nil {
		return err
	}

	b.logger.Printf("Creating metadata for %s...", config.OutputFile)

	variables, err := b.buildVariablesMap(config.Variables)
	if err != nil {
		return err
	}

	buildInput := builder.BuildInput{
		FormDirectories:          config.FormDirectories,
		IconPath:                 config.IconPath,
		InstanceGroupDirectories: config.InstanceGroupDirectories,
		JobDirectories:           config.JobDirectories,
		MetadataPath:             config.Metadata,
		PropertyDirectories:      config.PropertyDirectories,
		ReleaseTarballs:          releaseTarballs,
		RuntimeConfigDirectories: config.RuntimeConfigDirectories,
		StemcellTarball:          config.StemcellTarball,
		VariableDirectories:      config.VariableDirectories,
		Version:                  config.Version,
	}

	generatedMetadata, err := b.metadataBuilder.Build(buildInput)
	if err != nil {
		return err
	}

	b.logger.Println("Marshaling metadata file...")

	generatedMetadataYAML, err := yaml.Marshal(generatedMetadata)
	if err != nil {
		return err
	}

	writeInput := builder.WriteInput{
		OutputFile:           config.OutputFile,
		StubReleases:         config.StubReleases,
		MigrationDirectories: config.MigrationDirectories,
		ReleaseDirectories:   config.ReleaseDirectories,
		EmbedPaths:           config.EmbedPaths,
	}

	interpolatedMetadata, err := b.interpolateMetadata(variables, generatedMetadataYAML)
	if err != nil {
		return err
	}

	err = b.tileWriter.Write(generatedMetadata.Name, interpolatedMetadata, writeInput)
	if err != nil {
		return err
	}

	return nil
}

func (b Bake) Usage() commands.Usage {
	return commands.Usage{
		Description:      "Bakes tile metadata, stemcell, releases, and migrations into a format that can be consumed by OpsManager.",
		ShortDescription: "bakes a tile",
		Flags:            b.Options,
	}
}

func (b Bake) parseArgs(args []string) (BakeConfig, error) {
	config := BakeConfig{}

	args, err := flags.Parse(&config, args)
	if err != nil {
		panic(err)
	}

	if len(config.ReleaseDirectories) == 0 {
		return config, errors.New("Please specify release tarballs directory with the --releases-directory parameter")
	}

	if config.StemcellTarball == "" {
		return config, errors.New("--stemcell-tarball is a required parameter")
	}

	if config.IconPath == "" {
		return config, errors.New("--icon is a required parameter")
	}

	if config.Metadata == "" {
		return config, errors.New("--metadata is a required parameter")
	}

	if config.Version == "" {
		return config, errors.New("--version is a required parameter")
	}

	if config.OutputFile == "" {
		return config, errors.New("--output-file is a required parameter")
	}

	if len(config.InstanceGroupDirectories) == 0 && len(config.JobDirectories) > 0 {
		return config, errors.New("--jobs-directory flag requires --instance-groups-directory to also be specified")
	}

	return config, nil
}

func (b Bake) buildVariablesMap(flagVariables []string) (map[string]string, error) {
	variables := map[string]string{}
	for _, variable := range flagVariables {
		variablePair := strings.SplitN(variable, "=", 2)
		if len(variablePair) < 2 {
			return nil, errors.New("variable needs a key value in the form of key=value")
		}
		variables[variablePair[0]] = variablePair[1]
	}

	return variables, nil
}

func (b Bake) extractReleaseTarballFilenames(config BakeConfig) ([]string, error) {
	var releaseTarballs []string

	for _, releasesDirectory := range config.ReleaseDirectories {
		files, err := ioutil.ReadDir(releasesDirectory)
		if err != nil {
			return []string{}, err
		}

		for _, file := range files {
			matchTarballs, _ := regexp.MatchString("tgz$|tar.gz$", file.Name())
			if !matchTarballs {
				continue
			}

			releaseTarballs = append(releaseTarballs, filepath.Join(releasesDirectory, file.Name()))
		}
	}

	return releaseTarballs, nil
}

func (b Bake) interpolateMetadata(variables map[string]string, generatedMetadataYAML []byte) ([]byte, error) {
	templateHelpers := template.FuncMap{
		"variable": func(key string) (string, error) {
			val, ok := variables[key]
			if !ok {
				return "", fmt.Errorf("could not find variable with key '%s'", key)
			}
			return val, nil
		},
	}

	t, err := template.New("metadata").
		Delims("$(", ")").
		Funcs(templateHelpers).
		Parse(string(generatedMetadataYAML))

	if err != nil {
		return nil, fmt.Errorf("template parsing failed: %s", err)
	}

	var buffer bytes.Buffer
	err = t.Execute(&buffer, variables)
	if err != nil {
		return nil, fmt.Errorf("template execution failed: %s", err)
	}
	return buffer.Bytes(), nil
}
