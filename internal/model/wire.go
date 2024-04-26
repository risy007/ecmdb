//go:build wireinject

package model

import (
	"github.com/Duke1616/ecmdb/internal/model/internal/repository"
	"github.com/Duke1616/ecmdb/internal/model/internal/repository/dao"
	"github.com/Duke1616/ecmdb/internal/model/internal/service"
	"github.com/Duke1616/ecmdb/internal/model/internal/web"
	"github.com/Duke1616/ecmdb/internal/relation"
	"github.com/Duke1616/ecmdb/pkg/mongox"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	web.NewHandler,
	initMGProvider,
	initModelProvider)

func InitModule(db *mongox.Mongo, rmModule *relation.Module) (*Module, error) {
	wire.Build(
		ProviderSet,
		wire.FieldsOf(new(*relation.Module), "RMSvc"),
		wire.Struct(new(Module), "*"),
	)
	return new(Module), nil
}

var initMGProvider = wire.NewSet(
	service.NewMGService,
	repository.NewMGRepository,
	dao.NewModelGroupDAO,
)

var initModelProvider = wire.NewSet(
	service.NewModelService,
	repository.NewModelRepository,
	dao.NewModelDAO)
