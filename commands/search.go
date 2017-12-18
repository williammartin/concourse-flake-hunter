package commands

import (
	"fmt"
	"os"

	"github.com/albertoleal/concourse-flake-hunter/fly"
	"github.com/albertoleal/concourse-flake-hunter/hunter"
	"github.com/urfave/cli"
)

var SearchCommand = cli.Command{
	Name:        "search",
	Usage:       "search <arguments>",
	Description: "Searches for flakes",

	Flags: []cli.Flag{
		cli.IntFlag{
			Name:  "limit, l",
			Usage: "Limit number of builds to check",
			Value: 50,
		},
		cli.StringFlag{
			Name:  "job, j",
			Usage: "Job name to search",
			Value: "",
		},
	},

	Action: func(ctx *cli.Context) error {
		if ctx.Args().First() == "" {
			return cli.NewExitError("need to provide a pattern", 1)
		}

		client := ctx.App.Metadata["client"].(fly.Client)

		searcher := hunter.NewSearcher(client)
		spec := hunter.SearchSpec{
			Pattern: ctx.Args().First(),
			Limit:   ctx.Int("limit"),
			Job:     ctx.String("job"),
		}
		builds, err := searcher.Search(spec)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		table := &Table{
			Content: [][]string{},
			Header:  []string{"pipeline/job", "build url"},
		}

		for _, build := range builds {
			line := []string{}
			line = append(line, fmt.Sprintf("%s/%s", build.PipelineName, build.JobName))
			line = append(line, fmt.Sprintf("%s", build.ConcourseURL))
			table.Content = append(table.Content, line)
		}

		context := &Context{Stdout: os.Stdout}
		table.Render(context)
		return nil
	},
}
