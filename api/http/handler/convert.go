package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/relaxcloud-cn/html2md/internal/model"
	"github.com/relaxcloud-cn/html2md/internal/service"
)

// ConvertHandler HTML转换处理器
type ConvertHandler struct {
	service *service.ConvertService
}

// NewConvertHandler 创建转换处理器
func NewConvertHandler(service *service.ConvertService) *ConvertHandler {
	return &ConvertHandler{
		service: service,
	}
}

// Convert 转换HTML为Markdown
// @Summary 转换HTML为Markdown
// @Description 将HTML内容转换为Markdown格式，支持多种转换选项和插件
// @Tags 转换
// @Accept json
// @Produce json
// @Param request body model.ConvertRequest true "转换请求参数"
// @Success 200 {object} model.APIResponse{data=model.ConvertResponse} "转换成功"
// @Failure 400 {object} model.APIResponse{data=interface{}} "请求参数错误"
// @Failure 500 {object} model.APIResponse{data=interface{}} "内部服务器错误"
// @Router /api/v1/convert [post]
func (h *ConvertHandler) Convert(c *gin.Context) {
	var req model.ConvertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse(
			model.CodeBadRequest,
			"请求参数格式错误: "+err.Error(),
			nil,
		))
		return
	}

	result, err := h.service.Convert(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse(
			model.CodeInternalError,
			"转换失败: "+err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(result))
}

// ConvertBatch 批量转换HTML为Markdown
// @Summary 批量转换HTML为Markdown
// @Description 批量转换多个HTML内容为Markdown格式
// @Tags 转换
// @Accept json
// @Produce json
// @Param request body model.BatchConvertRequest true "批量转换请求参数"
// @Success 200 {object} model.APIResponse{data=model.BatchConvertResponse} "转换成功"
// @Failure 400 {object} model.APIResponse{data=interface{}} "请求参数错误"
// @Failure 500 {object} model.APIResponse{data=interface{}} "内部服务器错误"
// @Router /api/v1/convert/batch [post]
func (h *ConvertHandler) ConvertBatch(c *gin.Context) {
	var req model.BatchConvertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse(
			model.CodeBadRequest,
			"请求参数格式错误: "+err.Error(),
			nil,
		))
		return
	}

	result, err := h.service.ConvertBatch(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse(
			model.CodeInternalError,
			"批量转换失败: "+err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(result))
}

// Health 健康检查
// @Summary 健康检查
// @Description 检查服务健康状态和运行信息
// @Tags 系统
// @Produce json
// @Success 200 {object} model.APIResponse{data=model.HealthResponse} "服务正常"
// @Failure 500 {object} model.APIResponse{data=interface{}} "服务异常"
// @Router /api/v1/health [get]
func (h *ConvertHandler) Health(c *gin.Context) {
	result, err := h.service.Health()
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse(
			model.CodeInternalError,
			"健康检查失败: "+err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(result))
}

// GetConverterInfo 获取转换器信息
// @Summary 获取转换器信息
// @Description 获取转换器版本、支持的插件和功能信息
// @Tags 系统
// @Produce json
// @Success 200 {object} model.APIResponse{data=map[string]interface{}} "获取成功"
// @Router /api/v1/info [get]
func (h *ConvertHandler) GetConverterInfo(c *gin.Context) {
	info := h.service.GetConverterInfo()
	c.JSON(http.StatusOK, model.NewSuccessResponse(info))
}

// ConvertSimple 简单转换接口（GET方式）
// @Summary 简单HTML转换（GET方式）
// @Description 通过URL参数进行简单的HTML转换，适用于快速测试
// @Tags 转换
// @Accept plain
// @Produce json
// @Param html query string true "HTML内容"
// @Success 200 {object} model.APIResponse{data=model.ConvertResponse} "转换成功"
// @Failure 400 {object} model.APIResponse{data=interface{}} "请求参数错误"
// @Failure 500 {object} model.APIResponse{data=interface{}} "内部服务器错误"
// @Router /api/v1/convert/simple [get]
func (h *ConvertHandler) ConvertSimple(c *gin.Context) {
	html := c.Query("html")
	if html == "" {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse(
			model.CodeBadRequest,
			"HTML参数不能为空",
			nil,
		))
		return
	}

	req := &model.ConvertRequest{
		HTML: html,
	}

	result, err := h.service.Convert(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse(
			model.CodeInternalError,
			"转换失败: "+err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(result))
}
