package archive_bin

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"

	"github.com/google/blueprint"
	"github.com/roman-mazur/bood"
)

var (
	pctx          = blueprint.NewPackageContext("github.com/MaryLynJuana/KPI_Assembly_System/tree/master/build/modules/archive_bin")
	goArchive_bin = pctx.StaticRule("archive", blueprint.RuleParams{
		Command:     "cd $workDir && go test ${pkg}",
		Description: "generating test archive files",
	}, "workDir", "pkg")
)

type archive_binMod struct {
	blueprint.SimpleName

	properties struct {
		name   string
		binary string
	}
}

func (tb *archive_binMod) GenerateBuildActions(ctx blueprint.ModuleContext) {
	config := bood.ExtractConfig(ctx)
	config.Debug.Printf("Adding a build acction to generate the zip file for", tb.properties.name)

	archivePath := path.Join(config.BaseOutputDir, fmt.Sprintf("archives/%s.zip", tb.properties.name))
	binaryPath := path.Join(config.BaseOutputDir, fmt.Sprintf("bin/%s", tb.properties.binary))

	// Generate the binary file first
	cmd := exec.Command("go build -o " + binaryPath + " ./" + tb.properties.name)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		config.Info.Fatal("Error building the file")
	}

	// Then archive it
	if err := ZipFiles(archivePath, []string{binaryPath}); err != nil {
		panic(err)
	}

	ctx.Build(pctx, blueprint.BuildParams{
		Description: fmt.Sprintf("Generating a zip archive", tb.properties.name),
		Rule:        goArchive_bin,
		Outputs:     []string{archivePath},
		Args: map[string]string{
			"workDir": ctx.ModuleDir(),
			"pkg":     tb.properties.name,
		},
	})
}

func Archive_binFactory() (blueprint.Module, []interface{}) {
	mType := &archive_binMod{}
	return mType, []interface{}{&mType.SimpleName.Properties, &mType.properties}
}

func ZipFiles(filename string, files []string) error {

	newZipFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer newZipFile.Close()

	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	// Add files to zip
	for _, file := range files {
		if err = AddFileToZip(zipWriter, file); err != nil {
			return err
		}
	}
	return nil
}

func AddFileToZip(zipWriter *zip.Writer, filename string) error {

	fileToZip, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fileToZip.Close()

	info, err := fileToZip.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	header.Name = filename
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, fileToZip)
	return err
}
