package certificate

import (
	"fmt"
	"strconv"

	"casper_go/core/logger"
	"casper_go/core/response"

	"github.com/gin-gonic/gin"
)

// AddAPI godoc
// @Summary 添加证书监控
// @Description 添加一个新的 SSL 证书监控项
// @Tags 证书监控
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body AddRequest true "证书监控信息"
// @Success 200 {object} response.Response{data=Response}
// @Failure 400 {object} response.Response
// @Router /api/certificates [post]
func AddAPI(c *gin.Context) {
	var req AddRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Certificate.Warnf("[证书监控] 请求参数错误: %v", err)
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	logger.Certificate.Infof("[证书监控] 尝试添加证书: %s", req.URL)

	monitor, err := Add(req.URL)
	if err != nil {
		logger.Certificate.Errorf("[证书监控] 添加失败 URL=%s, 错误: %v", req.URL, err)
		response.Error(c, 400, err.Error())
		return
	}

	logger.Certificate.Infof("[证书监控] 添加成功: %s (ID=%d)", monitor.URL, monitor.ID)
	response.SuccessMsg(c, "添加成功", monitor.ToResponse())
}

// GetListAPI godoc
// @Summary 获取证书监控列表
// @Description 分页获取所有证书监控项
// @Tags 证书监控
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} response.Response{data=response.PageData}
// @Failure 401 {object} response.Response
// @Router /api/certificates [get]
func GetListAPI(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	monitors, total, err := GetList(page, pageSize)
	if err != nil {
		response.ServerError(c, "获取列表失败: "+err.Error())
		return
	}

	var responses []*Response
	for _, monitor := range monitors {
		responses = append(responses, monitor.ToResponse())
	}

	response.Page(c, responses, total, page, pageSize)
}

// GetDetailAPI godoc
// @Summary 获取证书监控详情
// @Description 根据 ID 获取证书监控的详细信息
// @Tags 证书监控
// @Security BearerAuth
// @Param id path int true "监控项 ID"
// @Success 200 {object} response.Response{data=Response}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/certificates/{id} [get]
func GetDetailAPI(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的 ID")
		return
	}

	monitor, err := GetByID(uint(id))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, monitor.ToResponse())
}

// UpdateAPI godoc
// @Summary 更新证书信息
// @Description 手动触发证书信息更新
// @Tags 证书监控
// @Security BearerAuth
// @Param id path int true "监控项 ID"
// @Success 200 {object} response.Response{data=Response}
// @Failure 400 {object} response.Response
// @Router /api/certificates/{id} [put]
func UpdateAPI(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的 ID")
		return
	}

	monitor, err := Update(uint(id))
	if err != nil {
		response.ServerError(c, "更新失败: "+err.Error())
		return
	}

	response.SuccessMsg(c, "更新成功", monitor.ToResponse())
}

// UpdateInfoAPI godoc
// @Summary 更新证书客户名和备注
// @Description 更新证书的客户名和备注信息
// @Tags 证书监控
// @Security BearerAuth
// @Param id path int true "监控项 ID"
// @Param request body UpdateRequest true "更新信息"
// @Success 200 {object} response.Response{data=Response}
// @Failure 400 {object} response.Response
// @Router /api/certificates/{id}/info [patch]
func UpdateInfoAPI(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的 ID")
		return
	}

	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	monitor, err := UpdateInfo(uint(id), req.CustomerName, req.Remark)
	if err != nil {
		response.ServerError(c, "更新失败: "+err.Error())
		return
	}

	response.SuccessMsg(c, "保存成功", monitor.ToResponse())
}

// DeleteAPI godoc
// @Summary 删除证书监控
// @Description 根据 ID 删除证书监控项
// @Tags 证书监控
// @Security BearerAuth
// @Param id path int true "监控项 ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/certificates/{id} [delete]
func DeleteAPI(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的 ID")
		return
	}

	if err := Delete(uint(id)); err != nil {
		response.ServerError(c, "删除失败: "+err.Error())
		return
	}

	response.SuccessMsg(c, "删除成功", nil)
}

// CheckAllAPI godoc
// @Summary 检查所有证书
// @Description 手动触发所有证书的检查更新（并发处理）
// @Tags 证书监控
// @Security BearerAuth
// @Success 200 {object} response.Response{data=CheckAllResult}
// @Failure 500 {object} response.Response
// @Router /api/certificates/check-all [post]
func CheckAllAPI(c *gin.Context) {
	logger.Certificate.Info("[API] 开始批量检查证书")

	result, err := CheckAll()
	if err != nil {
		response.ServerError(c, "检查失败: "+err.Error())
		return
	}

	message := fmt.Sprintf("检查完成：总数 %d，成功 %d，失败 %d，耗时 %s",
		result.Total, result.Success, result.Failed, result.Duration)

	logger.Certificate.Infof("[API] %s", message)
	response.SuccessMsg(c, message, result)
}

// GenerateCSRAPI godoc
// @Summary 生成 CSR
// @Description 在线生成 SSL 证书签名请求（CSR）和私钥
// @Tags 证书工具
// @Accept json
// @Produce json
// @Param request body GenerateCSRRequest true "CSR 生成参数"
// @Success 200 {object} response.Response{data=GenerateCSRResponse}
// @Failure 400 {object} response.Response
// @Router /api/certificates/tools/generate-csr [post]
func GenerateCSRAPI(c *gin.Context) {
	var req GenerateCSRRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Certificate.Warnf("[CSR生成] 请求参数错误: %v", err)
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	logger.Certificate.Infof("[CSR生成] 开始生成 CSR, CN=%s, 算法=%s, 密钥长度=%d",
		req.CommonName, req.KeyAlgorithm, req.KeySize)

	result, err := GenerateCSR(&req)
	if err != nil {
		logger.Certificate.Errorf("[CSR生成] 生成失败: %v", err)
		response.ServerError(c, "生成 CSR 失败: "+err.Error())
		return
	}

	logger.Certificate.Infof("[CSR生成] 生成成功, CN=%s", req.CommonName)
	response.SuccessMsg(c, "CSR 生成成功", result)
}

// ValidateCSRAPI godoc
// @Summary 验证 CSR
// @Description 验证 CSR 文件的合法性并提取信息
// @Tags 证书工具
// @Accept json
// @Produce json
// @Param request body ValidateCSRRequest true "CSR 内容"
// @Success 200 {object} response.Response{data=ValidateCSRResponse}
// @Failure 400 {object} response.Response
// @Router /api/certificates/tools/validate-csr [post]
func ValidateCSRAPI(c *gin.Context) {
	var req ValidateCSRRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Certificate.Warnf("[CSR验证] 请求参数错误: %v", err)
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	logger.Certificate.Infof("[CSR验证] 开始验证 CSR")

	result, err := ValidateCSR(req.CSRContent)
	if err != nil {
		logger.Certificate.Errorf("[CSR验证] 验证失败: %v", err)
		response.ServerError(c, "验证 CSR 失败: "+err.Error())
		return
	}

	if result.Valid {
		logger.Certificate.Infof("[CSR验证] 验证成功, CN=%s", result.CommonName)
	} else {
		logger.Certificate.Warnf("[CSR验证] 验证失败: %s", result.ErrorMessage)
	}

	response.Success(c, result)
}

// ValidateCertificateAPI godoc
// @Summary 验证证书
// @Description 在线验证 SSL 证书文件并提取详细信息
// @Tags 证书工具
// @Accept json
// @Produce json
// @Param request body ValidateCertRequest true "证书内容"
// @Success 200 {object} response.Response{data=ValidateCertResponse}
// @Failure 400 {object} response.Response
// @Router /api/certificates/tools/validate-cert [post]
func ValidateCertificateAPI(c *gin.Context) {
	var req ValidateCertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Certificate.Warnf("[证书验证] 请求参数错误: %v", err)
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	logger.Certificate.Infof("[证书验证] 开始验证证书")

	result, err := ValidateCertificate(req.CertContent, req.PrivateKeyContent)
	if err != nil {
		logger.Certificate.Errorf("[证书验证] 验证失败: %v", err)
		response.ServerError(c, "验证证书失败: "+err.Error())
		return
	}

	if result.Valid {
		if result.KeyPairChecked {
			if result.KeyPairMatched {
				logger.Certificate.Infof("[证书验证] 验证成功, CN=%s, 剩余%d天, 证书私钥配对✓", result.CommonName, result.DaysLeft)
			} else {
				logger.Certificate.Warnf("[证书验证] 验证成功, CN=%s, 剩余%d天, 证书私钥不匹配✗", result.CommonName, result.DaysLeft)
			}
		} else {
			logger.Certificate.Infof("[证书验证] 验证成功, CN=%s, 剩余%d天", result.CommonName, result.DaysLeft)
		}
	} else {
		logger.Certificate.Warnf("[证书验证] 验证失败: %s", result.ErrorMessage)
	}

	response.Success(c, result)
}
