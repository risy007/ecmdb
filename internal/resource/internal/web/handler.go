package web

import (
	"github.com/Duke1616/ecmdb/internal/attribute"
	"github.com/Duke1616/ecmdb/internal/resource/internal/domain"
	"github.com/Duke1616/ecmdb/internal/resource/internal/service"
	"github.com/Duke1616/ecmdb/pkg/ginx"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc          service.Service
	attributeSvc attribute.Service
}

func NewHandler(service service.Service, attributeSvc attribute.Service) *Handler {
	return &Handler{
		svc:          service,
		attributeSvc: attributeSvc,
	}
}

func (h *Handler) RegisterRoutes(server *gin.Engine) {
	g := server.Group("/resource")

	g.POST("/create/:model_identifies", ginx.WrapBody[CreateResourceReq](h.CreateResource))
	g.POST("/detail/:model_identifies", ginx.WrapBody[DetailResourceReq](h.DetailResource))
}

func (h *Handler) CreateResource(ctx *gin.Context, req CreateResourceReq) (ginx.Result, error) {
	modelIdentifies := ctx.Param("model_identifies")

	id, err := h.svc.CreateResource(ctx, domain.Resource{
		ModelIdentifies: modelIdentifies,
		Data:            req.Data,
	})

	if err != nil {
		return systemErrorResult, err
	}

	return ginx.Result{
		Data: id,
		Msg:  "创建资源成功",
	}, nil
}

func (h *Handler) DetailResource(ctx *gin.Context, req DetailResourceReq) (ginx.Result, error) {
	modelIdentifies := ctx.Param("model_identifies")
	attributes, err := h.attributeSvc.SearchAttributeByModelIdentifies(ctx, modelIdentifies)
	if err != nil {
		return systemErrorResult, err
	}

	var dmAttr domain.DetailResource
	dmAttr.Attributes = make([]domain.Attribute, 0, len(attributes))
	for _, v := range attributes {
		val := domain.Attribute{
			ID:              v.ID,
			ModelIdentifies: v.ModelIdentifies,
			Identifies:      v.Identifies,
			Name:            v.Name,
			FieldType:       v.FieldType,
			Required:        v.Required,
		}

		dmAttr.Attributes = append(dmAttr.Attributes, val)
	}
	dmAttr.ID = req.ID

	resp, err := h.svc.FindResourceById(ctx, dmAttr)
	if err != nil {
		return systemErrorResult, err
	}

	return ginx.Result{
		Data: resp,
		Msg:  "查看资源详情成功",
	}, nil
}
