package resources

import (
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/tasks"
	"github.com/stretchr/testify/require"
)

func TestNewTask(t *testing.T) {
	task := tasks.NewTask()
	require.NotNil(t, task)
	require.NotNil(t, task.Arguments)
}
