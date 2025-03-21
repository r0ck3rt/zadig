/*
Copyright 2022 The KodeRover Authors.

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

package step

import (
	"fmt"

	codehostmodels "github.com/koderover/zadig/v2/pkg/microservice/systemconfig/core/codehost/repository/models"
	"github.com/koderover/zadig/v2/pkg/types"
)

type StepGitSpec struct {
	CodeHosts []*codehostmodels.CodeHost `bson:"codehosts"      json:"codehosts"  yaml:"codehosts"`
	Repos     []*types.Repository        `bson:"repos"          json:"repos"      yaml:"repos"`
	Proxy     *Proxy                     `bson:"proxy"          json:"proxy"      yaml:"proxy"`
}

const (
	//	Oauth prefix
	OauthTokenPrefix = "oauth2"

	FileName = "reaper.tar.gz"
)

func (p *Proxy) GetProxyURL() string {
	var uri string
	if p.NeedPassword {
		uri = fmt.Sprintf("%s://%s:%s@%s:%d",
			p.Type,
			p.Username,
			p.Password,
			p.Address,
			p.Port,
		)
		return uri
	}

	uri = fmt.Sprintf("%s://%s:%d",
		p.Type,
		p.Address,
		p.Port,
	)
	return uri
}
