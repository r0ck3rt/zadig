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

package handler

import (
	"encoding/json"

	"github.com/gin-gonic/gin"

	"github.com/koderover/zadig/v2/pkg/microservice/aslan/core/collaboration/service"
	internalhandler "github.com/koderover/zadig/v2/pkg/shared/handler"
	e "github.com/koderover/zadig/v2/pkg/tool/errors"
	"github.com/koderover/zadig/v2/pkg/tool/log"
)

func GetCollaborationNew(c *gin.Context) {
	ctx := internalhandler.NewContext(c)
	defer func() { internalhandler.JSONResponse(c, ctx) }()
	projectName := c.Query("projectName")
	if projectName == "" {
		ctx.RespErr = e.ErrInvalidParam.AddDesc("projectName can not be empty")
		return
	}
	ctx.Resp, ctx.RespErr = service.GetCollaborationNew(projectName, ctx.UserID, ctx.IdentityType, ctx.Account, ctx.Logger)
}

// @Summary Sync Collaboration Instance
// @Description Sync Collaboration Instance
// @Tags 	collaboration
// @Accept 	json
// @Produce json
// @Param 	projectName		query		string									true	"project name"
// @Param 	body 			body 		service.SyncCollaborationInstanceArgs 	true 	"body"
// @Success 200
// @Router /api/aslan/collaboration/collaborations/sync [post]
func SyncCollaborationInstance(c *gin.Context) {
	ctx := internalhandler.NewContext(c)
	defer func() { internalhandler.JSONResponse(c, ctx) }()
	projectName := c.Query("projectName")
	args := new(service.SyncCollaborationInstanceArgs)
	data, err := c.GetRawData()
	if err != nil {
		log.Errorf("SyncCollaborationInstance c.GetRawData() err: %s", err)
		ctx.RespErr = e.ErrInvalidParam.AddDesc(err.Error())
		return
	}
	if err = json.Unmarshal(data, args); err != nil {
		log.Errorf("SyncCollaborationInstance json.Unmarshal err: %s", err)
		ctx.RespErr = e.ErrInvalidParam.AddDesc(err.Error())
		return
	}
	if projectName == "" {
		ctx.RespErr = e.ErrInvalidParam.AddDesc("projectName can not be empty")
		return
	}
	ctx.RespErr = service.SyncCollaborationInstance(args, projectName, ctx.UserID, ctx.IdentityType, ctx.Account, ctx.RequestID, ctx.Logger)
}

func CleanCIResources(c *gin.Context) {
	ctx := internalhandler.NewContext(c)
	defer func() { internalhandler.JSONResponse(c, ctx) }()
	ctx.RespErr = service.CleanCIResources(ctx.UserName, ctx.RequestID, ctx.Logger)
}
