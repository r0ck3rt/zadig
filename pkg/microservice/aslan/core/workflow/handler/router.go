/*
Copyright 2021 The KodeRover Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package handler

import (
	"github.com/gin-gonic/gin"
)

type Router struct{}

func (*Router) Inject(router *gin.RouterGroup) {
	// ---------------------------------------------------------------------------------------
	// 对外公共接口
	// ---------------------------------------------------------------------------------------
	webhook := router.Group("webhook")
	{
		webhook.POST("", ProcessGitWebHook)
	}

	// ---------------------------------------------------------------------------------------
	// Server Sent Events 接口
	// ---------------------------------------------------------------------------------------
	sse := router.Group("sse")
	{
		sse.GET("/workflowTasks/running", RunningWorkflowTasksSSE)
		sse.GET("/workflowTasks/pending", PendingWorkflowTasksSSE)
	}

	// ---------------------------------------------------------------------------------------
	// Pipeline 任务管理接口
	// ---------------------------------------------------------------------------------------
	taskV2 := router.Group("v2/tasks")
	{
		// TODO: MOVE ME TO ANOTHER PLACE
		taskV2.GET("/workflow/:pipelineName/taskId/:taskId", GetArtifactFile)
	}

	// ---------------------------------------------------------------------------------------
	// Pipeline Favorite 接口
	// ---------------------------------------------------------------------------------------
	favorite := router.Group("favorite")
	{
		favorite.POST("", CreateFavoritePipeline)
		favorite.DELETE("/:productName/:name/:type", DeleteFavoritePipeline)
		favorite.GET("", ListFavoritePipelines)
	}

	// ---------------------------------------------------------------------------------------
	// 新版本 通用工作流（暂命名） 接口
	// ---------------------------------------------------------------------------------------
	workflowV4 := router.Group("v4")
	{
		workflowV4.POST("", CreateWorkflowV4)
		workflowV4.POST("/workflowtask/:workflowName/field", SetWorkflowTasksCustomFields)
		workflowV4.GET("/workflowtask/:workflowName/field", GetWorkflowTasksCustomFields)
		workflowV4.GET("", ListWorkflowV4)
		workflowV4.POST("/auto", AutoCreateWorkflow)
		workflowV4.GET("/trigger", ListWorkflowV4CanTrigger)
		workflowV4.POST("/lint", LintWorkflowV4)
		workflowV4.POST("/check/:name", CheckWorkflowV4Approval)
		workflowV4.POST("/output/:jobName", GetWorkflowGlobalVars)
		workflowV4.POST("/repo/:jobName", GetWorkflowRepoIndex)
		workflowV4.GET("/name/:name", FindWorkflowV4)
		workflowV4.PUT("/:name", UpdateWorkflowV4)
		workflowV4.DELETE("/:name", DeleteWorkflowV4)
		workflowV4.GET("/preset/:name", GetWorkflowV4Preset)
		workflowV4.POST("/dynamicVariable/available", GetAvailableWorkflowV4DynamicVariable)
		workflowV4.POST("/dynamicVariable/render", GetWorkflowV4DynamicVariableValues)
		workflowV4.GET("/webhook/preset", GetWebhookForWorkflowV4Preset)
		workflowV4.GET("/webhook", ListGithookForWorkflowV4)
		workflowV4.POST("/webhook/:workflowName", CreateGithookForWorkflowV4)
		workflowV4.PUT("/webhook/:workflowName", UpdateGithookForWorkflowV4)
		workflowV4.DELETE("/webhook/:workflowName/trigger/:triggerName", DeleteGithookForWorkflowV4)
		workflowV4.GET("/jirahook/preset", GetJiraHookForWorkflowV4Preset)
		workflowV4.GET("/jirahook/:workflowName", ListJiraHookForWorkflowV4)
		workflowV4.POST("/jirahook/:workflowName", CreateJiraHookForWorkflowV4)
		workflowV4.PUT("/jirahook/:workflowName", UpdateJiraHookForWorkflowV4)
		workflowV4.DELETE("/jirahook/:workflowName/:hookName", DeleteJiraHookForWorkflowV4)
		workflowV4.GET("/meegohook/preset", GetMeegoHookForWorkflowV4Preset)
		workflowV4.GET("/meegohook/:workflowName", ListMeegoHookForWorkflowV4)
		workflowV4.POST("/meegohook/:workflowName", CreateMeegoHookForWorkflowV4)
		workflowV4.PUT("/meegohook/:workflowName", UpdateMeegoHookForWorkflowV4)
		workflowV4.DELETE("/meegohook/:workflowName/:hookName", DeleteMeegoHookForWorkflowV4)
		workflowV4.GET("/generalhook/preset", GetGeneralHookForWorkflowV4Preset)
		workflowV4.GET("/generalhook/:workflowName", ListGeneralHookForWorkflowV4)
		workflowV4.POST("/generalhook/:workflowName", CreateGeneralHookForWorkflowV4)
		workflowV4.PUT("/generalhook/:workflowName", UpdateGeneralHookForWorkflowV4)
		workflowV4.DELETE("/generalhook/:workflowName/:hookName", DeleteGeneralHookForWorkflowV4)
		workflowV4.POST("/generalhook/:workflowName/:hookName/webhook", GeneralHookEventHandler)
		workflowV4.GET("/cron/preset", GetCronForWorkflowV4Preset)
		workflowV4.GET("/cron", ListCronForWorkflowV4)
		workflowV4.POST("/cron/:workflowName", CreateCronForWorkflowV4)
		workflowV4.PUT("/cron", UpdateCronForWorkflowV4)
		workflowV4.DELETE("/cron/:workflowName/trigger/:cronID", DeleteCronForWorkflowV4)
		workflowV4.POST("/patch", GetPatchParams)
		workflowV4.GET("/sharestorage", CheckShareStorageEnabled)
		workflowV4.GET("/all", ListAllAvailableWorkflows)
		workflowV4.POST("/mse/render", RenderMseServiceYaml)
		workflowV4.GET("/mse/offline", GetMseOfflineResources)
		workflowV4.GET("/mse/:envName/tag", GetMseTagsInEnv)
		workflowV4.GET("/bluegreen/:envName/:serviceName", GetBlueGreenServiceK8sServiceYaml)
		workflowV4.GET("/jenkins/:id/:jobName", GetJenkinsJobParams)
		workflowV4.POST("/sql/validate", ValidateSQL)
		workflowV4.POST("/deploy/mergeImage", HelmDeployJobMergeImage)
	}

	// ---------------------------------------------------------------------------------------
	// workflow v4 任务接口
	// ---------------------------------------------------------------------------------------
	taskV4 := router.Group("v4/workflowtask")
	{
		taskV4.POST("", CreateWorkflowTaskV4)
		taskV4.GET("/filter/workflow/:name", GetWorkflowTaskFilters)
		taskV4.GET("", ListWorkflowTaskV4ByFilter)
		taskV4.GET("/workflow/:workflowName/task/:taskID", GetWorkflowTaskV4)
		taskV4.DELETE("/workflow/:workflowName/task/:taskID", CancelWorkflowTaskV4)
		taskV4.GET("/clone/workflow/:workflowName/task/:taskID", CloneWorkflowTaskV4)
		taskV4.GET("/view/workflow/:workflowName/task/:taskID", ViewWorkflowTaskV4)
		taskV4.POST("/retry/workflow/:workflowName/task/:taskID", RetryWorkflowTaskV4)
		taskV4.POST("/manualexec/workflow/:workflowName/task/:taskID", ManualExecWorkflowTaskV4)
		taskV4.GET("/manualexec/workflow/:workflowName/task/:taskID", GetManualExecWorkflowTaskV4Info)
		taskV4.POST("/breakpoint/:workflowName/:jobName/task/:taskID/:position", SetWorkflowTaskV4Breakpoint)
		taskV4.POST("/debug/:workflowName/task/:taskID", EnableDebugWorkflowTaskV4)
		taskV4.DELETE("/debug/:workflowName/:jobName/task/:taskID/:position", StopDebugWorkflowTaskJobV4)
		taskV4.POST("/revert/:workflowName/:jobName/task/:taskID", RevertWorkflowTaskV4Job)
		taskV4.GET("/revert/:workflowName/:jobName/task/:taskID", GetWorkflowTaskV4JobRevert)
		taskV4.POST("/approve", ApproveStage)
		taskV4.POST("/handle/error", HandleJobError)
		taskV4.GET("/workflow/:workflowName/taskId/:taskId/job/:jobName", GetWorkflowV4ArtifactFileContent)
		taskV4.GET("/workflow/:workflowName/taskId/:taskId/job/:jobName/build", GetWorkflowV4BuildJobArtifactFile)
		taskV4.PUT("/workflow/:workflowName/taskId/:taskId/remark", UpdateWorkflowV4TaskRemark)
		taskV4.POST("/trigger", CreateWorkflowTaskV4ByBuildInTrigger)
	}

	// ---------------------------------------------------------------------------------------
	// workflow view 接口
	// ---------------------------------------------------------------------------------------
	view := router.Group("view")
	{
		view.POST("", CreateWorkflowView)
		view.GET("", ListWorkflowViewNames)
		view.GET("/preset", GetWorkflowViewPreset)
		view.DELETE("", DeleteWorkflowView)
		view.PUT("", UpdateWorkflowView)
	}

	// ---------------------------------------------------------------------------------------
	// plugin repo 接口
	// ---------------------------------------------------------------------------------------
	plugin := router.Group("plugin")
	{
		plugin.GET("/template", ListPluginTemplates)
		plugin.POST("", UpsertUserPluginRepository)
		plugin.POST("/enterprise", UpsertEnterprisePluginRepository)
		plugin.GET("", ListUnofficalPluginRepositories)
		plugin.DELETE("/:id", DeletePluginRepo)
	}
}

type OpenAPIRouter struct{}

func (*OpenAPIRouter) Inject(router *gin.RouterGroup) {
	common := router.Group("")
	{
		common.GET("", OpenAPIGetWorkflowV4List)
	}

	// custom workflow apis
	custom := router.Group("custom")
	{
		custom.POST("/task", CreateCustomWorkflowTask)
		custom.GET("/task", OpenAPIGetWorkflowTaskV4)
		custom.DELETE("/task", OpenAPICancelWorkflowTaskV4)
		custom.POST("/task/approve", OpenAPIApproveStage)
		custom.DELETE("", OpenAPIDeleteCustomWorkflowV4)
		custom.GET("/:name/detail", OpenAPIGetCustomWorkflowV4)
		custom.POST("/:name/task/:taskID", OpenAPIRetryCustomWorkflowTaskV4)
		custom.PUT("/:name/task/:taskID", OpenAPIUpdateWorkflowV4TaskRemark)
		custom.GET("/:name/tasks", OpenAPIGetCustomWorkflowTaskV4)

	}

	view := router.Group("view")
	{
		view.POST("", OpenAPICreateWorkflowView)
		view.GET("", OpenAPIGetWorkflowViews)
		view.PUT("/:name", OpenAPIUpdateWorkflowView)
		view.DELETE("/:name", OpenAPIDeleteWorkflowView)
	}
}
