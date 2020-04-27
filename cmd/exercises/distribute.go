package exercises

import (
	"fmt"
	"github.com/krakowski/ilias"
	"github.com/krakowski/ilias-cli/util"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
)

type Tutor struct {
	Id string `yaml:"id"`
	Hours int `yaml:"hours"`
	Count int
}

const (
	TutorsFilename = ".tutors.yml"
)

var (
	seed int64
)

var submissionsDistributeCommand = &cobra.Command{
	Use:   "distribute [exercise] [assignment]",
	Short: "Distributes submissions on a set of tutors",
	SilenceErrors: true,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		distribute(args[0], args[1])
	},
}

func distribute(exerciseId string, assignmentId string) {
	tutors, err := readTutors()
	if err != nil {
		log.Fatal(err)
	}

	client := util.NewIliasClient()

	spin := util.StartSpinner("Fetching submissions")
	submissions, err := client.Exercise.List(&ilias.ListParams{
		Reference:    exerciseId,
		Assignment:   assignmentId,
		IncludeEmpty: false,
	})

	if err != nil {
		spin.StopError(err)
		os.Exit(1)
	}

	spin.StopSuccess(fmt.Sprintf("Received %d entries", len(submissions)))

	assignments := assignSubmissions(tutors, submissions)
	workSpace := util.WorkSpace{
		Exercise: util.Exercise{
			Reference:  exerciseId,
			Assignment: assignmentId,
		},

		Table: util.Table{
			Reference:  "",
			Identifier: "",
		},

		Corrections:       assignments,
	}

	output, err := yaml.Marshal(workSpace)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(util.WorkspaceFilename, output, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
}

func assignSubmissions(tutors []Tutor, submissions []ilias.SubmissionMeta) map[string][]string {
	rand.Seed(seed)
	submissionCount := len(submissions)
	totalHours := sumHours(tutors)

	var totalCount = 0
	for i, _ := range tutors {
		tutor := &tutors[i]
		tutor.Count = int(math.Round(float64(tutor.Hours) / float64(totalHours) * float64(submissionCount)))
		totalCount += tutor.Count
	}

	for totalCount != submissionCount {
		if totalCount < submissionCount {
			tutors[rand.Intn(len(tutors))].Count += 1
			totalCount += 1
		} else {
			tutors[rand.Intn(len(tutors))].Count -= 1
			totalCount -= 1
		}
	}

	printTutorTable(tutors)

	assignments := make(map[string][]string)
	for _, tutor := range tutors {
		for i := 0; i < tutor.Count; i++ {
			submission := submissions[0]
			submissions = submissions[1:]
			assignments[tutor.Id] = append(assignments[tutor.Id], submission.Identifier)
		}
	}

	return assignments
}

func sumHours(tutors []Tutor) int {
	sum := 0
	for _, tutor := range tutors {
		sum += tutor.Hours
	}
	return sum
}

func readTutors() ([]Tutor, error) {
	buffer, err := ioutil.ReadFile(TutorsFilename)
	if err != nil {
		return nil, err
	}

	var tutors []Tutor
	if err := yaml.Unmarshal(buffer, &tutors); err != nil {
		return nil, err
	}

	return tutors, nil
}

func printTutorTable(tutors []Tutor) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Tutor", "Abgaben"})

	for _, tutor := range tutors {
		table.Append([]string{tutor.Id, strconv.Itoa(tutor.Count)})
	}

	table.Render()
}

func init() {
	submissionsDistributeCommand.Flags().Int64Var(&seed, "seed", 0, "RNG seed")
}
