package gomodule

import (
	"fmt"
	"github.com/google/blueprint"
	"github.com/roman-mazur/bood"
	"path"
)

var (
	pctx = blueprint.NewPackageContext("github.com/MaryLynJuana/KPI_Assembly_System/build/gomodule")

	goBuild = pctx.StaticRule("binaryBuild", blueprint.RuleParams{
		Command:     "cd $workDir && go build -o $outputPath $pkg",
		Description: "build go command $pkg",
	}, "workDir", "outputPath", "pkg")

	goVendor = pctx.StaticRule("vendor", blueprint.RuleParams{
		Command:     "cd $workDir && go mod vendor",
		Description: "vendor dependencies of $name",
	}, "workDir", "name")
)

type goTestedBinaryModuleType struct {
	blueprint.SimpleName

	properties struct {
		Pkg string
		Srcs []string
		SrcsExclude []string
		VendorFirst bool

		Deps []string
	}
}

func (gtb *goTestedBinaryModuleType) DynamicDependencies(blueprint.DynamicDependerModuleContext) []string {
	return gtb.properties.Deps
}

func (gtb *goTestedBinaryModuleType) GenerateBuildActions(ctx blueprint.ModuleContext) {
	name := ctx.ModuleName()
	config := bood.ExtractConfig(ctx)
	config.Debug.Printf("Adding build actions for go binary module '%s'", name)

	outputPath := path.Join(config.BaseOutputDir, "bin", name)

	var inputs []string
	inputErors := false
	for _, src := range gtb.properties.Srcs {
		if matches, err := ctx.GlobWithDeps(src, gtb.properties.SrcsExclude); err == nil {
			inputs = append(inputs, matches...)
		} else {
			ctx.PropertyErrorf("srcs", "Cannot resolve files that match pattern %s", src)
			inputErors = true
		}
	}
	if inputErors {
		return
	}

	if gtb.properties.VendorFirst {
		vendorDirPath := path.Join(ctx.ModuleDir(), "vendor")
		ctx.Build(pctx, blueprint.BuildParams{
			Description: fmt.Sprintf("Vendor dependencies of %s", name),
			Rule:        goVendor,
			Outputs:     []string{vendorDirPath},
			Implicits:   []string{path.Join(ctx.ModuleDir(), "go.mod")},
			Optional:    true,
			Args: map[string]string{
				"workDir": ctx.ModuleDir(),
				"name":    name,
			},
		})
		inputs = append(inputs, vendorDirPath)
	}

	ctx.Build(pctx, blueprint.BuildParams{
		Description: fmt.Sprintf("Build %s as Go binary", name),
		Rule:        goBuild,
		Outputs:     []string{outputPath},
		Implicits:   inputs,
		Args: map[string]string{
			"outputPath": outputPath,
			"workDir":    ctx.ModuleDir(),
			"pkg":        gtb.properties.Pkg,
		},
	})

}

func SimpleBinFactory() (blueprint.Module, []interface{}) {
	mType := &goTestedBinaryModuleType{}
	return mType, []interface{}{&mType.SimpleName.Properties, &mType.properties}
}