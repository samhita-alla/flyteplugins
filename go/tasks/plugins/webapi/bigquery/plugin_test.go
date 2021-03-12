package bigquery

import (
	"testing"
	"time"

	flyteIdlCore "github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/core"
	"github.com/flyteorg/flyteplugins/go/tasks/pluginmachinery/core"
	"github.com/stretchr/testify/assert"
	"google.golang.org/api/bigquery/v2"
	"google.golang.org/api/googleapi"
)

func TestFormatJobReference(t *testing.T) {
	t.Run("format job reference", func(t *testing.T) {
		jobReference := bigquery.JobReference{
			JobId:     "my-job-id",
			Location:  "EU",
			ProjectId: "flyte-test",
		}

		str := formatJobReference(jobReference)

		assert.Equal(t, "flyte-test:EU.my-job-id", str)
	})
}

func TestCreateTaskInfo(t *testing.T) {
	t.Run("create task info", func(t *testing.T) {
		resourceMeta := ResourceMetaWrapper{
			JobReference: bigquery.JobReference{
				JobId:     "my-job-id",
				Location:  "EU",
				ProjectId: "flyte-test",
			},
		}

		taskInfo := createTaskInfo(&resourceMeta)

		assert.Equal(t, 1, len(taskInfo.Logs))
		assert.Equal(t, flyteIdlCore.TaskLog{
			Uri:  "https://console.cloud.google.com/bigquery?project=flyte-test&j=bq:EU:my-job-id&page=queryresults",
			Name: "BigQuery Console",
		}, *taskInfo.Logs[0])
	})
}

func TestHandleCreateError(t *testing.T) {
	occurredAt := time.Now()
	taskInfo := core.TaskInfo{OccurredAt: &occurredAt}

	t.Run("handle 401", func(t *testing.T) {
		createError := googleapi.Error{
			Code:    401,
			Message: "user xxx is not authorized",
		}

		phase := handleCreateError(&createError, &taskInfo)

		assert.Equal(t, flyteIdlCore.ExecutionError{
			Code:    "http401",
			Message: "user xxx is not authorized",
			Kind:    flyteIdlCore.ExecutionError_USER,
		}, *phase.Err())
		assert.Equal(t, taskInfo, *phase.Info())
	})

	t.Run("handle 500", func(t *testing.T) {
		createError := googleapi.Error{
			Code:    500,
			Message: "oops",
		}

		phase := handleCreateError(&createError, &taskInfo)

		assert.Equal(t, flyteIdlCore.ExecutionError{
			Code:    "http500",
			Message: "oops",
			Kind:    flyteIdlCore.ExecutionError_SYSTEM,
		}, *phase.Err())
		assert.Equal(t, taskInfo, *phase.Info())
	})
}
