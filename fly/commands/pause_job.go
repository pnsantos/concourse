package commands

import (
	"fmt"

	"github.com/concourse/concourse/fly/commands/internal/flaghelpers"
	"github.com/concourse/concourse/fly/rc"
)

type PauseJobCommand struct {
	Job flaghelpers.JobFlag `short:"j" long:"job" required:"true" value-name:"PIPELINE/JOB" description:"Name of a job to pause"`
	TeamFlag
}

func (command *PauseJobCommand) Execute(args []string) error {
	pipelineName, jobName := command.Job.PipelineName, command.Job.JobName
	target, err := rc.LoadTarget(Fly.Target, Fly.Verbose)
	if err != nil {
		return err
	}

	err = target.Validate()
	if err != nil {
		return err
	}

	team := command.TeamFlag.GetTeamTarget(target)
	found, err := team.PauseJob(pipelineName, jobName)
	if err != nil {
		return err
	}

	if !found {
		if command.Team == "" {
			fmt.Println("hint: are you missing '--team' to specify the team for the build?")
		}
		return fmt.Errorf("%s/%s not found\n", pipelineName, jobName)
	}

	fmt.Printf("paused '%s'\n", jobName)

	return nil
}
