// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package event

import (
	engine2 "github.com/Bunny3th/easy-workflow/workflow/engine"
	"github.com/Duke1616/ecmdb/internal/engine"
	"github.com/Duke1616/ecmdb/internal/event/easyflow"
	"github.com/Duke1616/ecmdb/internal/event/producer"
	"github.com/Duke1616/ecmdb/internal/order"
	"github.com/Duke1616/ecmdb/internal/task"
	"github.com/Duke1616/ecmdb/internal/template"
	"github.com/ecodeclub/mq-api"
	"github.com/larksuite/oapi-sdk-go/v3"
	"gorm.io/gorm"
	"log"
	"sync"
)

// Injectors from wire.go:

func InitModule(q mq.MQ, db *gorm.DB, engineModule *engine.Module, taskModule *task.Module, orderModule *order.Module, templateModule *template.Module, lark2 *lark.Client) (*Module, error) {
	service := engineModule.Svc
	orderStatusModifyEventProducer, err := producer.NewOrderStatusModifyEventProducer(q)
	if err != nil {
		return nil, err
	}
	serviceService := taskModule.Svc
	service2 := templateModule.Svc
	service3 := orderModule.Svc
	processEvent := InitWorkflowEngineOnce(db, service, orderStatusModifyEventProducer, serviceService, service2, service3, lark2)
	module := &Module{
		Event: processEvent,
	}
	return module, nil
}

// wire.go:

var engineOnce = sync.Once{}

func InitWorkflowEngineOnce(db *gorm.DB, engineSvc engine.Service, producer2 producer.OrderStatusModifyEventProducer,
	taskSvc task.Service, templateSvc template.Service, orderSvc order.Service, lark2 *lark.Client) *easyflow.ProcessEvent {
	notify, err := easyflow.NewNotify(engineSvc, templateSvc, orderSvc, lark2)
	if err != nil {
		panic(err)
	}

	event, err := easyflow.NewProcessEvent(producer2, engineSvc, taskSvc, notify, orderSvc)
	if err != nil {
		panic(err)
	}

	engineOnce.Do(func() {
		engine2.DB = db
		if err = engine2.DatabaseInitialize(); err != nil {
			log.Fatalln("easy workflow 初始化数据表失败，错误:", err)
		}
		engine2.IgnoreEventError = false
	})

	return event
}
