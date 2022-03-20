package controller

import (
	"maintainman/model"
	"maintainman/service"

	"github.com/kataras/iris/v12"
)

func GetUser(ctx iris.Context) {
	id, _ := ctx.Values().GetUint("user_id")
	response := service.GetUserInfoByID(id)
	ctx.Values().Set("response", response)
}

func GetUserByID(ctx iris.Context) {
	id, _ := ctx.Params().GetUint("id")
	response := service.GetUserByID(id)
	ctx.Values().Set("response", response)
}

func GetAllUsers(ctx iris.Context) {
	req := &model.AllUserJson{}
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	response := service.GetAllUsers(req)
	ctx.Values().Set("response", response)
}

func UserLogin(ctx iris.Context) {
	aul := &model.LoginJson{}
	if err := ctx.ReadJSON(&aul); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	aul.LoginIP = ctx.RemoteAddr()
	response := service.UserLogin(aul)
	ctx.Values().Set("response", response)
}

func UserRenew(ctx iris.Context) {
	id, _ := ctx.Values().GetUint("user_id")
	response := service.UserRenew(id)
	ctx.Values().Set("response", response)
}

func UserRegister(ctx iris.Context) {
	aul := &model.ModifyUserJson{}
	if err := ctx.ReadJSON(&aul); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	aul.RoleName = ""
	aul.DivisionID = 0
	response := service.CreateUser(aul)
	ctx.Values().Set("response", response)
}

func CreateUser(ctx iris.Context) {
	aul := &model.ModifyUserJson{}
	if err := ctx.ReadJSON(&aul); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	aul.OperatorID, _ = ctx.Values().GetUint("user_id")
	response := service.CreateUser(aul)
	ctx.Values().Set("response", response)
}

func UpdateUser(ctx iris.Context) {
	aul := &model.ModifyUserJson{}
	if err := ctx.ReadJSON(&aul); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	aul.RoleName = ""
	aul.DivisionID = 0
	id, _ := ctx.Values().GetUint("user_id")
	aul.OperatorID = id
	response := service.UpdateUser(id, aul)
	ctx.Values().Set("response", response)
}

func UpdateUserByID(ctx iris.Context) {
	aul := &model.ModifyUserJson{}
	if err := ctx.ReadJSON(&aul); err != nil {
		ctx.Values().Set("response", model.ErrorInvalidData(err))
		return
	}
	aul.OperatorID, _ = ctx.Values().GetUint("user_id")
	id, _ := ctx.Params().GetUint("id")
	response := service.UpdateUser(id, aul)
	ctx.Values().Set("response", response)
}

func DeleteUserByID(ctx iris.Context) {
	id, _ := ctx.Params().GetUint("id")
	response := service.DeleteUser(id)
	ctx.Values().Set("response", response)
}
