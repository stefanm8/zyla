package cli

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/asdine/storm"
	"github.com/gorilla/mux"

	"github.com/stefanm8/zyla/core"
	"gopkg.in/yaml.v2"
)

const (
	ApplicationsFolderName = "applications"
	ViewFileName           = "view.go"
	ModelFileName          = "model.go"
	composeFileName        = "zyla-compose.yaml"
	packageName            = "github.com/stefanm8/zyla"
	commandNameIdx         = 0
	commandArgsIdx         = 1
	dbName                 = "bolt.db"
)

var (
	ErrorComposeNotFound = errors.New("Compose file not found")
	ErrorInvalidCompose  = errors.New("Invalid configuration")
	ErrorInvalidYaml     = errors.New("Invalid syntax")
)

//CLI helper
type CLI struct {
	Commands    map[string]Command
	Args        []string
	WorkDir     string
	TplRoot     string
	Initializer *core.Initializer
	Project     *core.Project
}

// NewCLI constructor
func RunFromCLI() *CLI {
	// app := Initializer(Config)
	cli := &CLI{Args: os.Args[1:]}
	// cli.App.Initialize()

	buildCommand := &buildproject{
		NewCommand("build", cli.Args, []string{})}
	createsuperuserCommand := &createsuperuser{
		NewCommand("createsuperuser", cli.Args, []string{})}
	runserverCommand := &runserver{
		NewCommand("runserver", cli.Args, []string{})}

	cli.Commands = map[string]Command{
		"build":           buildCommand,
		"createsuperuser": createsuperuserCommand,
		"runserver":       runserverCommand,
	}

	cli.setup()

	return cli
}

func (cli *CLI) Handle() {
	if len(cli.Args) < 1 {
		cli.Usage()
		os.Exit(0)
	}

	cmd := cli.Commands[cli.Args[commandNameIdx]]
	if cmd != nil {
		cmd.run(cli)
		return
	}
	fmt.Print("Invalid command\n")
	cli.Usage()
}

func (cli *CLI) setup() {
	cli.loadProjectFromYaml()
	cli.setPaths()
	cli.Initializer = &core.Initializer{
		Project: cli.Project,
		Router:  mux.NewRouter(),
	}
	var err error
	cli.Initializer.DB, err = storm.Open(core.DatabaseName)
	fmt.Println(err)
	// defer cli.Initializer.DB.Close()
}

// Usage prints all the available commands and an usage example
func (cli *CLI) Usage() {
	fmt.Printf("\nUsage:\n	%s [command]\n", getScriptName())
	fmt.Print("\nAvailable Commands:\n")
	for _, c := range cli.Commands {
		fmt.Print(fmt.Sprintf("	%s\n", c.usage()))
	}
}

/**
* Utility Functions
**/

func (cli *CLI) parseConfig() []byte {
	config := openFile(path.Join(cli.WorkDir, composeFileName))
	if len(config) < 1 {
		fmt.Println(ErrorComposeNotFound)
		os.Exit(0)
	}
	// fmt.Println(config)
	return config

}

func (cli *CLI) loadProjectFromYaml() {
	config := cli.parseConfig()
	var project core.Project
	err := yaml.Unmarshal(config, &project)
	cli.Project = &project
	if err != nil {
		fmt.Println(ErrorInvalidYaml)
		fmt.Println(err)
		os.Exit(0)
	}
}

func (cli *CLI) setPaths() {

	gopath := os.Getenv("GOPATH")
	currentDir, _ := os.Getwd()
	cli.WorkDir = currentDir
	cli.TplRoot = gopath + "/src/" + packageName + "/tpl"

}

func createDirIfNotExists(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, os.ModePerm)
	}
}

func createFileIfNotExists(path string) *os.File {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println("creating")
		f, e := os.Create(path)
		checkError(e)
		return f
	}
	fmt.Println("erro")
	return nil
}

func openFile(path string) []byte {
	b, _ := ioutil.ReadFile(path)
	return b
}

func parseTpl(model string, fields []string, tplPath string, tplName string, destFile *os.File) {
	// funcMap := template.FuncMap{
	// 	"Title": strings.Title,
	// }
	// app := NewApp(model, fields)
	// tmpl, e := template.New(tplName).Funcs(funcMap).ParseFiles(tplPath)
	// checkError(e)
	// err := tmpl.Execute(destFile, app)
	// checkError(err)
}

func getScriptName() string {
	path := strings.Split(os.Args[0], "/")
	return path[len(path)-1]
}

func checkError(e error) {
	if e != nil {
		fmt.Println(e)
		os.Exit(0)
	}
}
