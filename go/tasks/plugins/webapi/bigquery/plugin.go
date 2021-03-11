package bigquery

import (
	"context"
	"encoding/gob"
	"fmt"
	"github.com/flyteorg/flyteplugins/go/tasks/pluginmachinery/google"
	structpb "github.com/golang/protobuf/ptypes/struct"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"google.golang.org/api/bigquery/v2"
	"google.golang.org/api/googleapi"
	"time"

	flyteIdlCore "github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/core"
	pluginErrors "github.com/flyteorg/flyteplugins/go/tasks/errors"
	pluginsCore "github.com/flyteorg/flyteplugins/go/tasks/pluginmachinery/core"

	"google.golang.org/api/option"

	"github.com/flyteorg/flytestdlib/logger"

	"github.com/flyteorg/flytestdlib/promutils"

	"github.com/flyteorg/flyteplugins/go/tasks/pluginmachinery"
	"github.com/flyteorg/flyteplugins/go/tasks/pluginmachinery/core"
	"github.com/flyteorg/flyteplugins/go/tasks/pluginmachinery/webapi"
)

const (
	bigqueryQueryJobTask = "bigquery_query_job_task"
)

type Plugin struct {
	metricScope       promutils.Scope
	cfg               *Config
	googleTokenSource google.TokenSource
}

type ResourceWrapper struct {
	Status *bigquery.JobStatus
}

type ResourceMetaWrapper struct {
	K8sServiceAccount string
	Namespace         string
	JobReference      bigquery.JobReference
}

func (p Plugin) GetConfig() webapi.PluginConfig {
	return GetConfig().WebAPI
}

func (p Plugin) ResourceRequirements(_ context.Context, _ webapi.TaskExecutionContextReader) (
	namespace core.ResourceNamespace, constraints core.ResourceConstraintsSpec, err error) {

	// Resource requirements are assumed to be the same.
	return "default", p.cfg.ResourceConstraints, nil
}

func (p Plugin) Create(ctx context.Context, taskCtx webapi.TaskExecutionContextReader) (webapi.ResourceMeta,
	webapi.Resource, error) {
	return p.createImpl(ctx, taskCtx)
}

func (p Plugin) createImpl(ctx context.Context, taskCtx webapi.TaskExecutionContextReader) (*ResourceMetaWrapper,
	*ResourceWrapper, error) {

	taskTemplate, err := taskCtx.TaskReader().Read(ctx)
	jobId := taskCtx.TaskExecutionMetadata().GetTaskExecutionID().GetGeneratedName()

	if err != nil {
		return nil, nil, pluginErrors.Errorf(pluginErrors.BadTaskSpecification, "unable to fetch task specification [%v]", err.Error())
	}

	inputs, err := taskCtx.InputReader().Get(ctx)

	if err != nil {
		return nil, nil, pluginErrors.Errorf(pluginErrors.BadTaskSpecification, "unable to fetch task inputs [%v]", err.Error())
	}

	var job *bigquery.Job

	namespace := taskCtx.TaskExecutionMetadata().GetNamespace()
	k8sServiceAccount := taskCtx.TaskExecutionMetadata().GetK8sServiceAccount()
	identity := google.Identity{K8sNamespace: namespace, K8sServiceAccount: k8sServiceAccount}
	tokenSource, err := p.googleTokenSource.GetTokenSource(ctx, identity)

	if err != nil {
		return nil, nil, pluginErrors.Errorf(pluginErrors.TaskFailedUnknownError, "unable to get token source [%v]", err.Error())
	}

	client, err := newBigQueryClient(ctx, tokenSource)

	if err != nil {
		return nil, nil, pluginErrors.Errorf(pluginErrors.TaskFailedUnknownError, "unable to get bigquery client [%v]", err.Error())
	}

	if taskTemplate.Type == bigqueryQueryJobTask {
		job, err = createQueryJob(jobId, taskTemplate.GetCustom(), inputs)
	} else {
		err = errors.Errorf("unexpected task type [%v]", taskTemplate.Type)
	}

	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to create query job")
	}

	job.Configuration.Labels = taskCtx.TaskExecutionMetadata().GetLabels()

	resp, err := client.Jobs.Insert(job.JobReference.ProjectId, job).Do()

	if err != nil {
		apiError, ok := err.(*googleapi.Error)

		if ok && apiError.Code == 409 {
			resourceMeta := ResourceMetaWrapper{
				JobReference:      *job.JobReference,
				Namespace:         namespace,
				K8sServiceAccount: k8sServiceAccount,
			}

			job, err := client.Jobs.Get(resourceMeta.JobReference.ProjectId, resourceMeta.JobReference.JobId).Do()

			if err != nil {
				return nil, nil, errors.Wrapf(err, "failed to get job %v %v", resourceMeta.JobReference.ProjectId, resourceMeta.JobReference.JobId)
			}

			resource := ResourceWrapper{Status: job.Status}

			return &resourceMeta, &resource, nil
		} else {
			return nil, nil, errors.Wrap(err, "failed to create query job")
		}
	}

	resource := ResourceWrapper{Status: resp.Status}
	resourceMeta := ResourceMetaWrapper{
		JobReference:      *job.JobReference,
		Namespace:         namespace,
		K8sServiceAccount: k8sServiceAccount,
	}

	return &resourceMeta, &resource, nil
}

func createQueryJob(jobID string, custom *structpb.Struct, inputs *flyteIdlCore.LiteralMap) (*bigquery.Job, error) {
	queryConfig, err := unmarshalBigQueryQueryConfig(custom)

	if err != nil {
		return nil, pluginErrors.Errorf(pluginErrors.BadTaskSpecification, "can't unmarshall struct to BigQueryQueryJob")
	}

	jobConfigurationQuery, err := getJobConfigurationQuery(queryConfig, inputs)

	if err != nil {
		return nil, pluginErrors.Errorf(pluginErrors.BadTaskSpecification, "unable to fetch task inputs [%v]", err.Error())
	}

	jobReference := bigquery.JobReference{
		JobId:     jobID,
		Location:  queryConfig.Location,
		ProjectId: queryConfig.ProjectID,
	}

	return &bigquery.Job{
		Configuration: &bigquery.JobConfiguration{
			Query: jobConfigurationQuery,
		},
		JobReference: &jobReference,
	}, nil
}

func (p Plugin) Get(ctx context.Context, taskCtx webapi.GetContext) (latest webapi.Resource, err error) {
	return p.getImpl(ctx, taskCtx)
}

func (p Plugin) getImpl(ctx context.Context, taskCtx webapi.GetContext) (wrapper *ResourceWrapper, err error) {
	resourceMeta := taskCtx.ResourceMeta().(*ResourceMetaWrapper)

	tokenSource, err := p.googleTokenSource.GetTokenSource(ctx, google.Identity{
		K8sNamespace:      resourceMeta.Namespace,
		K8sServiceAccount: resourceMeta.K8sServiceAccount,
	})

	if err != nil {
		return nil, pluginErrors.Errorf(pluginErrors.TaskFailedUnknownError, "unable to get token source [%v]", err.Error())
	}

	client, err := newBigQueryClient(ctx, tokenSource)

	if err != nil {
		return nil, pluginErrors.Errorf(pluginErrors.TaskFailedUnknownError, "unable to get client [%v]", err.Error())
	}

	job, err := client.Jobs.Get(resourceMeta.JobReference.ProjectId, resourceMeta.JobReference.JobId).Do()

	if err != nil {
		return nil, errors.Wrapf(err, "failed to get job %v %v", resourceMeta.JobReference.ProjectId, resourceMeta.JobReference.JobId)
	}

	// Only cache fields we want to keep in memory instead of the potentially huge execution closure.
	return &ResourceWrapper{
		Status: job.Status,
	}, nil
}

func (p Plugin) Delete(ctx context.Context, taskCtx webapi.DeleteContext) error {
	if taskCtx.ResourceMeta() == nil {
		return nil
	}

	resourceMeta := taskCtx.ResourceMeta().(*ResourceMetaWrapper)
	tokenSource, err := p.googleTokenSource.GetTokenSource(ctx, google.Identity{
		K8sNamespace:      resourceMeta.Namespace,
		K8sServiceAccount: resourceMeta.K8sServiceAccount,
	})

	if err != nil {
		return pluginErrors.Errorf(pluginErrors.TaskFailedUnknownError, "unable to get token source [%v]", err.Error())
	}

	client, err := newBigQueryClient(ctx, tokenSource)

	if err != nil {
		return err
	}

	_, err = client.Jobs.Cancel(resourceMeta.JobReference.ProjectId, resourceMeta.JobReference.JobId).Do()

	if err != nil {
		return err
	}

	logger.Info(ctx, "Cancelled job [%v:%v.%v]", resourceMeta.JobReference.ProjectId, resourceMeta.JobReference.Location, resourceMeta.JobReference.JobId)

	return nil
}

func (p Plugin) Status(_ context.Context, tCtx webapi.StatusContext) (phase core.PhaseInfo, err error) {
	resourceMeta := tCtx.ResourceMeta().(*ResourceMetaWrapper)
	resource := tCtx.Resource().(*ResourceWrapper)
	version := pluginsCore.DefaultPhaseVersion

	if resource == nil {
		return core.PhaseInfoUndefined, nil
	}

	taskInfo := createTaskInfo(resourceMeta)

	switch resource.Status.State {
	case "PENDING":
		return core.PhaseInfoRunning(version, taskInfo), nil

	case "RUNNING":
		return core.PhaseInfoRunning(version, taskInfo), nil

	case "DONE":
		// TODO handle error_result
		if resource.Status.ErrorResult != nil {
			switch resource.Status.ErrorResult.Reason {
			case "":
				// fallthrough
			default:
				// TODO handle all codes as retryable/non-retriable properly

				execError := &flyteIdlCore.ExecutionError{
					Message: resource.Status.ErrorResult.Message,
					Kind:    flyteIdlCore.ExecutionError_USER,
					Code:    resource.Status.ErrorResult.Reason,
				}

				return pluginsCore.PhaseInfoFailed(pluginsCore.PhaseRetryableFailure, execError, taskInfo), nil
			}
		}

		return pluginsCore.PhaseInfoSuccess(taskInfo), nil
	}

	return core.PhaseInfoUndefined, pluginErrors.Errorf(pluginsCore.SystemErrorCode, "unknown execution phase [%v].", resource.Status.State)
}

func createTaskInfo(resourceMeta *ResourceMetaWrapper) *core.TaskInfo {
	timeNow := time.Now()

	j := fmt.Sprintf("bq:%v:%v", resourceMeta.JobReference.Location, resourceMeta.JobReference.JobId)

	return &core.TaskInfo{
		OccurredAt: &timeNow,
		Logs: []*flyteIdlCore.TaskLog{
			{
				Uri: fmt.Sprintf("https://console.cloud.google.com/bigquery?project=%v&j=%v&page=queryresults",
					resourceMeta.JobReference.ProjectId,
					j),
				Name: "BigQuery Console",
			},
		},
	}
}

func newBigQueryClient(ctx context.Context, tokenSource oauth2.TokenSource) (*bigquery.Service, error) {
	options := []option.ClientOption{
		option.WithScopes("https://www.googleapis.com/auth/bigquery"),
		// FIXME how do I access current version?
		option.WithUserAgent(fmt.Sprintf("%s/%s", "flytepropeller", "LATEST")),
		option.WithTokenSource(tokenSource),
	}

	return bigquery.NewService(ctx, options...)
}

func NewPlugin(cfg *Config, metricScope promutils.Scope) (*Plugin, error) {
	googleTokenSource, err := google.NewTokenSource(cfg.GoogleTokenSource)

	if err != nil {
		return nil, pluginErrors.Wrapf(pluginErrors.PluginInitializationFailed, err, "failed to get google token source")
	}

	return &Plugin{
		metricScope:       metricScope,
		cfg:               cfg,
		googleTokenSource: googleTokenSource,
	}, nil
}

func init() {
	gob.Register(ResourceMetaWrapper{})
	gob.Register(ResourceWrapper{})

	pluginmachinery.PluginRegistry().RegisterRemotePlugin(webapi.PluginEntry{
		ID:                 "bigquery",
		SupportedTaskTypes: []core.TaskType{bigqueryQueryJobTask},
		PluginLoader: func(ctx context.Context, iCtx webapi.PluginSetupContext) (webapi.AsyncPlugin, error) {
			return NewPlugin(GetConfig(), iCtx.MetricsScope())
		},
	})
}
