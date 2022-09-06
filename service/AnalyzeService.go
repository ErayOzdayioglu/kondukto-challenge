package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/go-git/go-git/v5"
	"kondukto.com/challenge/domain"
	"kondukto.com/challenge/repository"
)

type DefaultAnalyzeService struct {
	Repository repository.ScanResultsRepository
}

type AnalyzeService interface {
	Analyze(string, string)
}

func (s DefaultAnalyzeService) Analyze(id string, url string) {
	clone(url, id)
	output := startBanditAndAnalyze(id)
	scanResult := parseOutput(output)
	scanResult.GithubUrl = url
	s.Repository.Update(id, *scanResult)

}

func clone(url string, id string) {
	fmt.Println(url)
	_, err := git.PlainClone("/tmp/src/"+id, false, &git.CloneOptions{
		URL: url,
	})

	if err != nil {
		log.Fatalln("error while git clone : " + err.Error())
	}
}

func startBanditAndAnalyze(id string) string {
	client, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	reader, err := client.ImagePull(context.Background(), "opensorcery/bandit", types.ImagePullOptions{})

	if err != nil {
		fmt.Println(err)
	} else {
		io.Copy(os.Stdout, reader)
	}

	resp, err := client.ContainerCreate(context.Background(), &container.Config{
		Image: "opensorcery/bandit",
		Cmd:   []string{"-r", "/code"},
	}, &container.HostConfig{
		Binds: []string{
			"/tmp/src/" + id + ":/code",
		}}, nil, nil, "")
	if err != nil {
		panic(err)
	}

	if err := client.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	statusCh, errCh := client.ContainerWait(context.Background(), resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			panic(err)
		}
	case <-statusCh:
	}

	out, err := client.ContainerLogs(context.Background(), resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		panic(err)
	}
	var buf bytes.Buffer

	stdcopy.StdCopy(&buf, &buf, out)

	return buf.String()
}

func parseOutput(str string) *domain.ScanResults {
	temp := strings.Split(str, "\n")

	arr := make([]string, 0)
	for _, s := range temp {
		str := strings.Replace(s, "\t", "", -1)
		arr = append(arr, str)
	}

	var scanResults domain.ScanResults
	var issuesByConfidence domain.IssuesByConfidence
	var issuesBySeverity domain.IssuesBySeverity

	issuesBySeverityIndex := indexOf("Total issues (by severity):", arr)
	issuesBySeverity.Undefined = parseNumberFromString(arr[issuesBySeverityIndex+1])
	issuesBySeverity.Low = parseNumberFromString(arr[issuesBySeverityIndex+2])
	issuesBySeverity.Medium = parseNumberFromString(arr[issuesBySeverityIndex+3])
	issuesBySeverity.High = parseNumberFromString(arr[issuesBySeverityIndex+4])

	issuesByConfidenceIndex := indexOf("Total issues (by confidence):", arr)
	issuesByConfidence.Undefined = parseNumberFromString(arr[issuesByConfidenceIndex+1])
	issuesByConfidence.Low = parseNumberFromString(arr[issuesByConfidenceIndex+2])
	issuesByConfidence.Medium = parseNumberFromString(arr[issuesByConfidenceIndex+3])
	issuesByConfidence.High = parseNumberFromString(arr[issuesByConfidenceIndex+4])

	scanResults.IssuesByConfidence = issuesByConfidence
	scanResults.IssuesBySeverity = issuesBySeverity

	return &scanResults
}

func indexOf(element string, data []string) int {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return -1 //not found.
}

func parseNumberFromString(str string) float64 {
	numberString := str[len(str)-3:]
	value, err := strconv.ParseFloat(numberString, 64)
	if err != nil {
		log.Fatalln("Error in conversion")
	}
	return value
}

func NewAnalyzeService(Repository repository.ScanResultsRepository) DefaultAnalyzeService {
	return DefaultAnalyzeService{Repository: Repository}
}
