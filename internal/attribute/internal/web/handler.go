package web

import (
	"github.com/Duke1616/ecmdb/internal/attribute/internal/domain"
	"github.com/Duke1616/ecmdb/internal/attribute/internal/service"
	"github.com/Duke1616/ecmdb/pkg/ginx"
	"github.com/ecodeclub/ekit/slice"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc service.Service
}

func NewHandler(svc service.Service) *Handler {
	return &Handler{
		svc: svc,
	}
}

func (h *Handler) RegisterRoutes(server *gin.Engine) {
	g := server.Group("/attribute")

	g.POST("/create", ginx.WrapBody[CreateAttributeReq](h.CreateAttribute))
	g.POST("/list", ginx.WrapBody[ListAttributeReq](h.ListAttributes))

	g.POST("/detail", ginx.WrapBody[DetailAttributeReq](h.DetailAttributeFields))
}

func (h *Handler) CreateAttribute(ctx *gin.Context, req CreateAttributeReq) (ginx.Result, error) {
	id, err := h.svc.CreateAttribute(ctx.Request.Context(), toDomain(req))

	if err != nil {
		return systemErrorResult, err
	}
	return ginx.Result{
		Data: id,
		Msg:  "添加模型属性成功",
	}, nil
}

func (h *Handler) DetailAttributeFields(ctx *gin.Context, req DetailAttributeReq) (ginx.Result, error) {
	fields, err := h.svc.SearchAttributeFieldsByModelUid(ctx, req.ModelUid)
	if err != nil {
		return systemErrorResult, err
	}

	return ginx.Result{
		Data: fields,
	}, nil
}

func (h *Handler) ListAttributes(ctx *gin.Context, req ListAttributeReq) (ginx.Result, error) {
	attrs, total, err := h.svc.ListAttributes(ctx, req.ModelUid)
	if err != nil {
		return systemErrorResult, err
	}

	return ginx.Result{
		Data: RetrieveAttributeList{
			Total: total,
			Attribute: slice.Map(attrs, func(idx int, src domain.Attribute) Attribute {
				return toAttributeVo(src)
			}),
		},
	}, nil
}
