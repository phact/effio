package effio

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// effio make -dev <file.json> -fio <dir> -path <dir>
func (cmd *Cmd) MakeSuite() {
	// the default device filename is <hostname>.json
	devfile, err := os.Hostname()
	if err != nil {
		devfile = "devices"
	}
	devfile = fmt.Sprintf("%s.json", devfile)

	// parse subcommand arguments
	var idFlag, devFlag, fioFlag, pathFlag string
	fs := cmd.FlagSet
	fs.StringVar(&idFlag, "id", "", "Id of the test suite")
	fs.StringVar(&devFlag, "dev", devfile, "JSON file containing device metadata")
	fs.StringVar(&fioFlag, "fio", "fio_configs/", "directory containing fio config templates")
	fs.StringVar(&pathFlag, "path", "./suites/", "generated suite is written to this path")
	fs.Parse(cmd.Args)

	if idFlag == "" {
		fs.Usage()
	}

	// load device data from json
	devs := LoadDevicesFile(mustAbs(devFlag))

	// load the fio config templates into memory
	templates := LoadFioConfDir(mustAbs(fioFlag))

	// use an absolute directory for pathFlag
	outDir := mustAbs(pathFlag)

	// build up a test suite of devs x templates
	suite := NewSuite(idFlag)
	suite.Populate(devs, templates)
	suite.WriteAll(outDir)
}

// mustAbs change any relative path to an absolute path
// any error from filepath.Abs is considered fatal
func mustAbs(p string) string {
	out, err := filepath.Abs(p)
	if err != nil {
		log.Fatalf("BUG: Required operation failed with error: %s\n", err)
	}

	return out
}
