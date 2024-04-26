package web

import (
	"github.com/Duke1616/ecmdb/internal/model/internal/domain"
	"github.com/Duke1616/ecmdb/internal/model/internal/service"
	"github.com/Duke1616/ecmdb/internal/relation"
	"github.com/Duke1616/ecmdb/pkg/ginx"
	"github.com/Duke1616/ecmdb/pkg/tools"
	"github.com/ecodeclub/ekit/slice"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc   service.Service
	mgSvc service.MGService
	RMSvc relation.RMSvc
}

func NewHandler(svc service.Service, groupSvc service.MGService, RMSvc relation.RMSvc) *Handler {
	return &Handler{
		svc:   svc,
		mgSvc: groupSvc,
		RMSvc: RMSvc,
	}
}

func (h *Handler) RegisterRoutes(server *gin.Engine) {
	g := server.Group("/model")
	// 模型分组
	g.POST("/group/create", ginx.WrapBody[CreateModelGroupReq](h.CreateGroup))

	// 模型操作
	g.POST("/create", ginx.WrapBody[CreateModelReq](h.CreateModel))
	g.POST("/detail", ginx.WrapBody[DetailModelReq](h.DetailModel))
	g.POST("/list", ginx.WrapBody[Page](h.ListModels))

	// 模型关联关系
	//g.POST("/relation/diagram", ginx.WrapBody[Page](h.FindRelationModelDiagram))
}

func (h *Handler) CreateGroup(ctx *gin.Context, req CreateModelGroupReq) (ginx.Result, error) {
	id, err := h.mgSvc.CreateModelGroup(ctx.Request.Context(), domain.ModelGroup{
		Name: req.Name,
	})

	if err != nil {
		return systemErrorResult, err
	}
	return ginx.Result{
		Data: id,
		Msg:  "添加模型分组成功",
	}, nil
}

func (h *Handler) CreateModel(ctx *gin.Context, req CreateModelReq) (ginx.Result, error) {
	id, err := h.svc.CreateModel(ctx, domain.Model{
		Name:    req.Name,
		GroupId: req.GroupId,
		UID:     req.UID,
		Icon:    req.Icon,
	})

	if err != nil {
		return systemErrorResult, err
	}
	return ginx.Result{
		Data: id,
		Msg:  "添加模型成功",
	}, nil
}

func (h *Handler) DetailModel(ctx *gin.Context, req DetailModelReq) (ginx.Result, error) {
	m, err := h.svc.FindModelById(ctx, req.ID)
	if err != nil {
		return systemErrorResult, err
	}

	return ginx.Result{
		Data: m,
		Msg:  "模型查找成功",
	}, nil
}

func (h *Handler) ListModels(ctx *gin.Context, req Page) (ginx.Result, error) {
	models, total, err := h.svc.ListModels(ctx, req.Offset, req.Limit)
	if err != nil {
		return systemErrorResult, err
	}

	return ginx.Result{
		Data: RetrieveModelsListResp{
			Total: total,
			Models: slice.Map(models, func(idx int, m domain.Model) Model {
				return toModelVo(m)
			}),
		},
	}, nil
}

func (h *Handler) FindRelationModelDiagram(ctx *gin.Context, req Page) (ginx.Result, error) {
	// TODO 为了后续加入 label 概念进行过滤先查询所有的模型
	// 查询所有模型
	models, _, err := h.svc.ListModels(ctx, req.Offset, req.Limit)
	if err != nil {
		return systemErrorResult, err
	}

	// 取出所有的 uids
	modelUidS := slice.Map(models, func(idx int, src domain.Model) string {
		return src.UID
	})

	// 查询包含的数据
	ds, err := h.RMSvc.ListSrcModelByUIDs(ctx, modelUidS)
	if err != nil {
		return systemErrorResult, err
	}

	// 生成关联节点的map
	var mds map[string][]relation.ModelDiagram
	mds = tools.ToMapS(ds, func(m relation.ModelDiagram) string {
		return m.SourceModelUid
	})

	// 返回 vo，前端展示
	diagrams := toModelDiagramVo(models, mds)

	return ginx.Result{
		Data: RetrieveRelationModelDiagram{
			Diagrams: diagrams,
		},
	}, nil
}
