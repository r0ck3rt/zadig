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

package task

import (
	"fmt"

	"github.com/koderover/zadig/pkg/microservice/aslan/config"
	"github.com/koderover/zadig/pkg/microservice/aslan/core/common/repository/models"
	"github.com/koderover/zadig/pkg/setting"
)

type JenkinsIntegration struct {
	URL      string `bson:"url"                   json:"url"`
	Username string `bson:"username"              json:"username"`
	Password string `bson:"password"              json:"password"`
}

type JenkinsBuildArgs struct {
	JobName            string               `bson:"job_name"            json:"job_name"`
	JenkinsBuildParams []*JenkinsBuildParam `bson:"jenkins_build_param" json:"jenkins_build_params"`
}

type JenkinsBuildParam struct {
	Name         string      `bson:"name,omitempt"                     json:"name,omitempt"`
	Value        interface{} `bson:"value,omitempt"                    json:"value,omitempt"`
	Type         string      `bson:"type,omitempty"                    json:"type,omitempty"`
	AutoGenerate bool        `bson:"auto_generate,omitempty"           json:"auto_generate,omitempty"`
	ChoiceOption []string    `bson:"choice_option,omitempty"           json:"choice_option,omitempty"`
}

// JenkinsBuild ...
type JenkinsBuild struct {
	TaskType           config.TaskType             `bson:"type"                    json:"type"`
	Enabled            bool                        `bson:"enabled"                 json:"enabled"`
	TaskStatus         config.Status               `bson:"status"                  json:"status"`
	ServiceName        string                      `bson:"service_name"            json:"service_name"`
	Service            string                      `bson:"service"                 json:"service"`
	OnSetup            string                      `bson:"setup,omitempty"         json:"setup,omitempty"`
	Timeout            int                         `bson:"timeout,omitempty"       json:"timeout,omitempty"`
	Error              string                      `bson:"error,omitempty"         json:"error,omitempty"`
	ResReq             setting.Request             `bson:"res_req"                 json:"res_req"`
	ResReqSpec         setting.RequestSpec         `bson:"res_req_spec"            json:"res_req_spec"`
	JenkinsIntegration *JenkinsIntegration         `bson:"-"                       json:"jenkins_integration"`
	JenkinsBuildArgs   *JenkinsBuildArgs           `bson:"-"                       json:"jenkins_build_args"`
	Registries         []*models.RegistryNamespace `bson:"registries"              json:"registries"`
	StartTime          int64                       `bson:"start_time,omitempty"    json:"start_time,omitempty"`
	EndTime            int64                       `bson:"end_time,omitempty"      json:"end_time,omitempty"`
	LogFile            string                      `bson:"log_file"                json:"log_file"`
	Image              string                      `bson:"image,omitempty"         json:"image,omitempty"`
	IsRestart          bool                        `bson:"is_restart"              json:"is_restart"`
}

// ToSubTask ...
func (j *JenkinsBuild) ToSubTask() (map[string]interface{}, error) {
	var task map[string]interface{}
	if err := IToi(j, &task); err != nil {
		return nil, fmt.Errorf("convert JenkinsBuildTask to interface error: %v", err)
	}
	return task, nil
}

// SetIntegration ...
func (j *JenkinsBuild) SetIntegration(jenkinsIntegration *JenkinsIntegration) {
	if jenkinsIntegration != nil {
		j.JenkinsIntegration = jenkinsIntegration
	}
}

// SetIntegration ...
func (j *JenkinsBuild) SetJenkinsBuildArgs(jenkinsBuildArgs *JenkinsBuildArgs) {
	if jenkinsBuildArgs != nil {
		j.JenkinsBuildArgs = jenkinsBuildArgs
	}
}
