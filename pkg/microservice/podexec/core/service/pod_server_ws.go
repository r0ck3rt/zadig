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

package service

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/koderover/zadig/v2/pkg/tool/clientmanager"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"

	commonmodels "github.com/koderover/zadig/v2/pkg/microservice/aslan/core/common/repository/models"
	commonrepo "github.com/koderover/zadig/v2/pkg/microservice/aslan/core/common/repository/mongodb"
	internalhandler "github.com/koderover/zadig/v2/pkg/shared/handler"
	e "github.com/koderover/zadig/v2/pkg/tool/errors"
	"github.com/koderover/zadig/v2/pkg/tool/kube/getter"
	"github.com/koderover/zadig/v2/pkg/tool/log"
)

func ServeWs(c *gin.Context) {
	ctx, err := internalhandler.NewContextWithAuthorization(c)
	defer func() { internalhandler.JSONResponse(c, ctx) }()

	if err != nil {
		ctx.RespErr = fmt.Errorf("authorization Info Generation failed: err %s", err)
		ctx.UnAuthorized = true
		return
	}

	podName := c.Param("podName")
	containerName := c.Param("containerName")

	if podName == "" {
		ctx.RespErr = e.ErrInvalidParam.AddDesc("containerName can't be empty,please check!")
		return
	}
	log.Infof("exec containerName: %s, pod: %s", containerName, podName)

	productName := c.Query("projectName")
	envName := c.Param("envName")
	productInfo, err := commonrepo.NewProductColl().Find(&commonrepo.ProductFindOptions{Name: productName, EnvName: envName})
	if err != nil {
		ctx.RespErr = e.ErrInternalError.AddDesc(fmt.Sprintf("failed to find product %s/%s, err: %s", productName, envName, err))
		return
	}
	namespace, clusterID := productInfo.Namespace, productInfo.ClusterID

	pty, err := NewTerminalSession(c.Writer, c.Request, nil)
	if err != nil {
		log.Errorf("get pty failed: %v", err)
		ctx.RespErr = e.ErrInternalError.AddDesc(fmt.Sprintf("get pty failed: %v", err))
		return
	}
	defer func() {
		log.Info("close session.")
		_ = pty.Close()
	}()

	kubeCli, err := clientmanager.NewKubeClientManager().GetKubernetesClientSet(clusterID)
	if err != nil {
		msg := fmt.Sprintf("get kubecli err :%v", err)
		log.Errorf(msg)
		_, _ = pty.Write([]byte(msg))
		pty.Done()

		ctx.RespErr = e.ErrInternalError.AddDesc(fmt.Sprintf("get kubecli err :%v", err))
		return
	}

	ok, err := ValidatePod(kubeCli, namespace, podName, containerName)
	if !ok {
		msg := fmt.Sprintf("Validate pod error! err: %v", err)
		log.Errorf(msg)
		_, _ = pty.Write([]byte(msg))
		pty.Done()

		ctx.RespErr = e.ErrInternalError.AddDesc(fmt.Sprintf("Validate pod error! err: %v", err))
		return
	}

	err = ExecPod(clusterID, []string{"/bin/sh"}, pty, namespace, podName, containerName)
	if err != nil {
		msg := fmt.Sprintf("Exec to pod error! err: %v", err)
		log.Errorf(msg)
		_, _ = pty.Write([]byte(msg))
		pty.Done()

		ctx.RespErr = e.ErrInternalError.AddDesc(fmt.Sprintf("Exec to pod error! err: %v", err))
		return
	}
}

func DebugWorkflow(c *gin.Context) {
	ctx := internalhandler.NewContext(c)
	defer func() { internalhandler.JSONResponse(c, ctx) }()
	logger := ctx.Logger
	taskID, err := strconv.ParseInt(c.Param("taskID"), 10, 64)
	if err != nil {
		ctx.RespErr = e.ErrInvalidParam.AddDesc("无效 task ID")
		return
	}

	ctx.RespErr = debugWorkflow(c, c.Param("workflowName"), c.Param("jobName"), taskID, logger)
	return
}

func debugWorkflow(c *gin.Context, workflowName, jobName string, taskID int64, logger *zap.SugaredLogger) error {
	workflowTask, err := commonrepo.NewworkflowTaskv4Coll().Find(workflowName, taskID)
	if err != nil {
		return e.ErrStopDebugShell.AddDesc(fmt.Sprintf("failed to find task: %s", err))
	}
	if workflowTask.Finished() {
		return e.ErrStopDebugShell.AddDesc("task has been finished")
	}

	var task *commonmodels.JobTask
FOR:
	for _, stage := range workflowTask.Stages {
		for _, jobTask := range stage.Jobs {
			if jobTask.Name == jobName {
				task = jobTask
				break FOR
			}
		}
	}
	if task == nil {
		logger.Error("debug workflow failed: not found job")
		return e.ErrInvalidParam.AddDesc("Job不存在")
	}
	log.Infof("DebugWorkflow: %s, %s, %d", workflowName, jobName, taskID)

	jobTaskSpec := &commonmodels.JobTaskFreestyleSpec{}
	if err := commonmodels.IToi(task.Spec, jobTaskSpec); err != nil {
		logger.Errorf("debug workflow failed: IToi %v", err)
		return e.ErrGetDebugShell.AddDesc("启动调试终端意外失败")
	}

	pty, err := NewTerminalSession(c.Writer, c.Request, nil, &TerminalSessionOption{
		SecretEnvs: func() (secrets []string) {
			for _, v := range jobTaskSpec.Properties.Envs {
				if v.IsCredential {
					secrets = append(secrets, v.Value)
				}
			}
			return secrets
		}(),
		Type: Workflow,
	})
	if err != nil {
		log.Errorf("get pty failed: %v", err)
		return e.ErrGetDebugShell.AddDesc(fmt.Sprintf("get pty failed: %v", err))
	}
	defer func() {
		log.Info("close session.")
		_ = pty.Close()
	}()

	kubeClient, err := clientmanager.NewKubeClientManager().GetControllerRuntimeClient(jobTaskSpec.Properties.ClusterID)
	if err != nil {
		log.Errorf("debug workflow failed: get kube client error: %s", err)
		return e.ErrGetDebugShell.AddDesc("启动调试终端意外失败: get kube client")
	}

	pods, err := getter.ListPods(jobTaskSpec.Properties.Namespace, labels.Set{"job-name": task.K8sJobName}.AsSelector(), kubeClient)
	if err != nil {
		logger.Errorf("debug workflow failed: list pods %v", err)
		return e.ErrGetDebugShell.AddDesc("启动调试终端意外失败: ListPods")
	}
	if len(pods) == 0 {
		logger.Error("debug workflow failed: list pods num 0")
		return e.ErrGetDebugShell.AddDesc("启动调试终端意外失败: ListPods num 0")
	}
	pod := pods[0]
	switch pod.Status.Phase {
	case corev1.PodRunning:
	default:
		logger.Errorf("debug workflow failed: pod status is %s", pod.Status.Phase)
		return e.ErrGetDebugShell.AddDesc(fmt.Sprintf("Job 状态 %s 无法启动调试终端", pod.Status.Phase))
	}

	var envs []string
	for _, env := range jobTaskSpec.Properties.Envs {
		removeDquoteVal := strings.ReplaceAll(env.Value, `"`, `\"`)
		removeBquoteVal := strings.ReplaceAll(removeDquoteVal, "`", "\\`")
		envs = append(envs, fmt.Sprintf("%s=\"%s\"", env.Key, removeBquoteVal))
	}
	script := ""
	if len(envs) != 0 {
		script += "env " + strings.Join(envs, " ") + " "
	}
	script += "bash\n"

	err = ExecPod(jobTaskSpec.Properties.ClusterID, []string{"/bin/sh", "-c", script}, pty, jobTaskSpec.Properties.Namespace, pod.Name, pod.Spec.Containers[0].Name)
	if err != nil {
		msg := fmt.Sprintf("Exec to pod error! err: %v", err)
		log.Errorf(msg)
		_, _ = pty.Write([]byte(msg))
		pty.Done()

		return e.ErrGetDebugShell.AddDesc(fmt.Sprintf("Exec to pod error! err: %v", err))
	}
	return nil
}
