package cli

import (
	"fmt"
	"go/build"
	"html/template"
	"net/http"
	"os"
	"path"
	"strings"

	"go.etcd.io/bbolt"

	"github.com/stefanm8/zyla/core"
)

// Command interface
type Command interface {
	run(*CLI)
	usage() string
}

//CommandBase facilitates usage prints for commands that embed it
type CommandBase struct {
	Name     string
	Args     []string
	ArgsName []string
}

// NewCommand returns CommandBase instance
func NewCommand(name string, args []string, argsname []string) *CommandBase {
	return &CommandBase{name, args, argsname}
}

func (c *CommandBase) usage() string {
	text := c.Name
	if len(c.ArgsName) > 0 {
		text := c.Name + " " + strings.Join(c.ArgsName, " ")
		return text
	}
	return text

}

/**
* Commands
**/

// Starts up the server
type runserver struct {
	*CommandBase
}

func (cmd *runserver) run(cli *CLI) {
	var err error
	cli.Initializer.DB.Bolt.Update(func(tx *bbolt.Tx) error {
		fmt.Println("Creatings")
		for _, app := range cli.Initializer.Project.Applications {
			bucketName := strings.ToLower(app.Model.BucketName)
			tx.CreateBucketIfNotExists([]byte(strings.ToLower(bucketName)))
		}
		return nil
	})
	checkError(err)
	// cli.Initializer.Router.Register()
	fmt.Println("Serving")
	http.ListenAndServe(":8080", cli.Initializer.Router)
}

// Builds the project from yaml
type buildproject struct {
	*CommandBase
}

func (cmd *buildproject) run(cli *CLI) {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}

	dependencies := map[string]*core.Dependency{
		"zyla": &core.Dependency{
			Name:   "zyla",
			Import: "github.com/stefanm8/zyla",
			User:   "stefanm8",
		},
		"storm": &core.Dependency{
			Name:   "storm",
			Import: "github.com/asdine/storm",
			User:   "asdine",
		},
	}

	appFolderPath := path.Join(cli.WorkDir, ApplicationsFolderName)
	createDirIfNotExists(appFolderPath)
	for _, app := range cli.Project.Applications {
		app.Dependecies = dependencies
		createDirIfNotExists(path.Join(appFolderPath, app.Name))
		generateApp(cli, app)
	}
	generateMain(cli)
}

// Creates a super user
type createsuperuser struct {
	*CommandBase
}

func (cmd *createsuperuser) run(cli *CLI) {
	// reader := bufio.NewReader(os.Stdin)
	// fmt.Print("Enter username: ")
	// username, _ := reader.ReadString('\n')
	// fmt.Print("Enter e-mail: ")
	// email, _ := reader.ReadString('\n')
	// fmt.Print("Enter password: ")
	// password, _ := reader.ReadString('\n')

	// fmt.Println(username)
	// fmt.Println(email)
	// fmt.Println(password)
	// superUser := &User{
	// 	Username:  username,
	// 	Email:     email,
	// 	Password:  password,
	// 	Superuser: true,
	// 	// BoltBucket:
	// }
	// op := superUser.Register()
	// fmt.Println(op)

}

func generateApp(cli *CLI, app *core.Application) {

	appFolderPath := path.Join(ApplicationsFolderName, app.Name)
	templatesToParse := []string{
		"view",
		"model",
		// "tests"
	}

	funcMap := template.FuncMap{
		"Title": strings.Title,
	}

	for _, tplName := range templatesToParse {
		tplNameFile := tplName + ".tgo"
		destFile := createFileIfNotExists(path.Join(appFolderPath, tplName+".go"))
		tplRootPath := path.Join(cli.TplRoot, tplNameFile)
		tmpl, e := template.New(tplNameFile).Funcs(funcMap).ParseFiles(tplRootPath)
		err := tmpl.Execute(destFile, app)
		checkError(e)
		checkError(err)
	}
}

func generateMain(cli *CLI) {
	funcMap := template.FuncMap{
		"Title": strings.Title,
	}
	destFile := createFileIfNotExists(path.Join(cli.WorkDir, "main.go"))
	tplRootPath := path.Join(cli.TplRoot, "main.tgo")
	tmpl, e := template.New("main.tgo").Funcs(funcMap).ParseFiles(tplRootPath)
	err := tmpl.Execute(destFile, cli.Project)
	checkError(e)
	checkError(err)
}
