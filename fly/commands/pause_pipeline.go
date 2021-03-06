package commands

import (
	"fmt"

	"github.com/concourse/concourse/fly/commands/internal/displayhelpers"
	"github.com/concourse/concourse/fly/commands/internal/flaghelpers"
	"github.com/concourse/concourse/fly/rc"
)

type PausePipelineCommand struct {
	Pipeline flaghelpers.PipelineFlag `short:"p"  long:"pipeline" description:"Pipeline to pause"`
	All      bool                     `short:"a"  long:"all"      description:"Pause all pipelines"`
}

func (command *PausePipelineCommand) Validate() error {
	return command.Pipeline.Validate()
}

func (command *PausePipelineCommand) Execute(args []string) error {
	if string(command.Pipeline) == "" && !command.All {
		displayhelpers.Failf("Either a pipeline name or --all are required")
	}

	if string(command.Pipeline) != "" && command.All {
		displayhelpers.Failf("A pipeline and --all can not both be specified")
	}

	err := command.Validate()
	if err != nil {
		return err
	}

	target, err := rc.LoadTarget(Fly.Target, Fly.Verbose)
	if err != nil {
		return err
	}

	err = target.Validate()
	if err != nil {
		return err
	}

	var pipelineNames []string
	if string(command.Pipeline) != "" {
		pipelineNames = []string{string(command.Pipeline)}
	}

	if command.All {
		pipelines, err := target.Team().ListPipelines()
		if err != nil {
			return err
		}

		for _, pipeline := range pipelines {
			pipelineNames = append(pipelineNames, pipeline.Name)
		}
	}

	for _, pipelineName := range pipelineNames {
		found, err := target.Team().PausePipeline(pipelineName)
		if err != nil {
			return err
		}

		if found {
			fmt.Printf("paused '%s'\n", pipelineName)
		} else {
			displayhelpers.Failf("pipeline '%s' not found\n", pipelineName)
		}
	}

	return nil
}
