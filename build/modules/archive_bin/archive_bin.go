package archive_bin

import (
	"fmt"
	"os"
	"path"

	"github.com/google/blueprint"
	"github.com/roman-mazur/bood"
)

var (
	pctx = blueprint.NewPackageContext("github.com/MaryLynJuana/KPI_Assembly_System/tree/master/build/modules/archive_bin")

	// goBuild = pctx.StaticRule("binaryBuild", blueprint.RuleParams{
	// 	Command:     "cd $workDir && go build -o $outputPath ./$pkg",
	// 	Description: "build go command $pkg",
	// }, "workDir", "outputPath", "pkg")

	goArchive = pctx.StaticRule("ZipArchiving", blueprint.RuleParams{
		Command:     "cd $workDir && zip -r $zipName $fileName",
		Description: "archiving command",
	}, "workDir", "zipName", "fileName")
)

type goArchive_binModule struct {
	blueprint.SimpleName

	properties struct {
		Name   string
		Binary string
	}
}

func (tb *goArchive_binModule) GenerateBuildActions(ctx blueprint.ModuleContext) {
	config := bood.ExtractConfig(ctx)
	config.Debug.Print("Adding a build acction to generate the zip file for", tb.properties.Name)

	archivePath := path.Join(config.BaseOutputDir, fmt.Sprintf("archives/%s.zip", tb.properties.Name))
	archiveDir := path.Join(config.BaseOutputDir, "archives")
	binaryPath := path.Join(config.BaseOutputDir, fmt.Sprintf("bin/%s", tb.properties.Binary))
	binaryDir := path.Join(config.BaseOutputDir, "bin")

	// Troubleshoot Linux/GoLang badcode
	os.MkdirAll(archiveDir, os.ModePerm)
	os.MkdirAll(binaryDir, os.ModePerm)
	os.Create(binaryPath)

	ctx.Build(pctx, blueprint.BuildParams{
		Description: fmt.Sprint("Generating a zip archive for module ", tb.properties.Name),
		Rule:        goArchive,
		Outputs:     []string{archivePath},
		Args: map[string]string{
			"workDir":  ctx.ModuleDir(),
			"zipName":  archivePath,
			"fileName": binaryPath,
		},
	})
}

func Archive_binFactory() (blueprint.Module, []interface{}) {
	mType := &goArchive_binModule{}
	return mType, []interface{}{&mType.SimpleName.Properties, &mType.properties}
}
