package common

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	wfv1 "github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
)

type DagContext interface {
	GetTask(taskName string) *wfv1.DAGTask
	GetTaskDependencies(taskName string) []string
	GetTaskFinishedAtTime(taskName string) time.Time
}

type TaskResult string

const (
	TaskResultSucceeded TaskResult = "Succeeded"
	TaskResultFailed    TaskResult = "Failed"
	TaskResultErrored   TaskResult = "Errored"
	TaskResultSkipped   TaskResult = "Skipped"
	TaskResultCompleted TaskResult = "Completed"
	TaskResultDaemoned  TaskResult = "Daemoned"
)

var (
	// TODO: This should use validate.workflowFieldNameFmt, but we can't import it here because an import cycle would be created
	taskNameRegex   = regexp.MustCompile(`([a-zA-Z0-9][-a-zA-Z0-9]*?\.[A-Z][a-z]+)|([a-zA-Z0-9][-a-zA-Z0-9]*)`)
	taskResultRegex = regexp.MustCompile(`([a-zA-Z0-9][-a-zA-Z0-9]*?\.[A-Z][a-z]+)`)
)

type expansionMatch struct {
	taskName string
	start    int
	end      int
}

func GetTaskDependencies(task *wfv1.DAGTask, ctx DagContext) ([]string, string) {
	depends := getTaskDependsLogic(task, ctx)
	matches := taskNameRegex.FindAllStringSubmatchIndex(depends, -1)
	var expansionMatches []expansionMatch
	dependencies := make(map[string]bool)
	for _, matchGroup := range matches {
		// We have matched a taskName.TaskResult
		if matchGroup[2] != -1 {
			match := depends[matchGroup[2]:matchGroup[3]]
			split := strings.Split(match, ".")
			dependencies[split[0]] = true
		} else if matchGroup[4] != -1 {
			match := depends[matchGroup[4]:matchGroup[5]]
			dependencies[match] = true
			expansionMatches = append(expansionMatches, expansionMatch{taskName: match, start: matchGroup[4], end: matchGroup[5]})
		}
	}

	var out []string
	for dependency := range dependencies {
		out = append(out, dependency)
	}

	if len(expansionMatches) == 0 {
		return out, depends
	}

	sort.Slice(expansionMatches, func(i, j int) bool {
		// Sort in descending order
		return expansionMatches[i].start > expansionMatches[j].start
	})
	for _, match := range expansionMatches {
		matchTask := ctx.GetTask(match.taskName)
		depends = depends[:match.start] + expandDependency(match.taskName, matchTask) + depends[match.end:]
	}

	return out, depends
}

func TopologicalSort(leafTasks []string, ctx DagContext) [][]string {
	maxLevel := 0
	levels := make(map[string]int)

	var visit func(string, int)
	visit = func(taskName string, level int) {
		if level > maxLevel {
			maxLevel = level
		}

		if prevLevel, visited := levels[taskName]; visited {
			if level > prevLevel {
				levels[taskName] = level
			}
			return
		}

		for _, task := range ctx.GetTaskDependencies(taskName) {
			visit(task, level+1)
		}

		levels[taskName] = level
	}

	for _, taskName := range leafTasks {
		visit(taskName, 0)
	}

	topologicalSort := make([][]string, maxLevel+1)
	for taskName, level := range levels {
		topologicalSort[maxLevel-level] = append(topologicalSort[maxLevel-level], taskName)
	}
	return topologicalSort
}

func ValidateTaskResults(dagTask *wfv1.DAGTask) error {
	// If a user didn't specify a depends expression, there are no task results to validate
	if dagTask.Depends == "" {
		return nil
	}

	matches := taskResultRegex.FindAllStringSubmatch(dagTask.Depends, -1)
	for _, matchGroup := range matches {
		split := strings.Split(matchGroup[1], ".")
		taskName, taskResult := split[0], TaskResult(split[1])
		switch taskResult {
		case TaskResultSucceeded, TaskResultFailed, TaskResultSkipped, TaskResultCompleted, TaskResultErrored, TaskResultDaemoned:
			// Do nothing
		default:
			return fmt.Errorf("task result '%s' for task '%s' is invalid", taskResult, taskName)
		}
	}
	return nil
}

func getTaskDependsLogic(dagTask *wfv1.DAGTask, ctx DagContext) string {
	if dagTask.Depends != "" {
		return dagTask.Depends
	}

	// For backwards compatibility, "dependencies: [A, B]" is equivalent to "depends: (A.Successful || A.Skipped || A.Daemoned)) && (B.Successful || B.Skipped || B.Daemoned)"
	var dependencies []string
	for _, dependency := range dagTask.Dependencies {
		depTask := ctx.GetTask(dependency)
		dependencies = append(dependencies, expandDependency(dependency, depTask))
	}
	return strings.Join(dependencies, " && ")
}

func expandDependency(depName string, depTask *wfv1.DAGTask) string {
	resultForTask := func(result TaskResult) string { return fmt.Sprintf("%s.%s", depName, result) }

	taskDepends := []string{resultForTask(TaskResultSucceeded), resultForTask(TaskResultSkipped), resultForTask(TaskResultDaemoned)}
	if depTask.ContinueOn != nil {
		if depTask.ContinueOn.Error {
			taskDepends = append(taskDepends, resultForTask(TaskResultErrored))
		}
		if depTask.ContinueOn.Failed {
			taskDepends = append(taskDepends, resultForTask(TaskResultFailed))
		}
	}
	return "(" + strings.Join(taskDepends, " || ") + ")"
}

// GetTaskAncestry returns a list of taskNames which are ancestors of this task.
// The list is ordered by the tasks finished time.
func GetTaskAncestry(ctx DagContext, taskName string) []string {
	visited := make(map[string]time.Time)

	var getAncestry func(currTask string)
	getAncestry = func(currTask string) {
		if _, seen := visited[currTask]; seen {
			return
		}
		for _, depTask := range ctx.GetTaskDependencies(currTask) {
			getAncestry(depTask)
		}
		if currTask != taskName {
			if _, ok := visited[currTask]; !ok {
				visited[currTask] = ctx.GetTaskFinishedAtTime(currTask)
			}
		}
	}

	getAncestry(taskName)

	ancestry := make([]string, len(visited))
	for newTask, newFinishedAt := range visited {
		insertTask(visited, ancestry, newTask, newFinishedAt)
	}

	return ancestry
}

// insertTask inserts the newTaskName at the right position ordered by time into the ancestry list.
func insertTask(visited map[string]time.Time, ancestry []string, newTaskName string, finishedAt time.Time) {
	for i, taskName := range ancestry {
		if taskName == "" {
			ancestry[i] = newTaskName
			return
		}
		if finishedAt.Before(visited[taskName]) {
			// insert at position i and shift others
			copy(ancestry[i+1:], ancestry[i:])
			ancestry[i] = newTaskName
			return
		}
	}
}
