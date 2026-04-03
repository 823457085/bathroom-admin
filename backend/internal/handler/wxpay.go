package handler

import (
	"bytes"
	"crypto/md5"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"bathroom-admin/internal/model"

	"github.com/gin-gonic/gin"
)

type WxPayHandler struct {
	orderRepo *model.OrderRepository
	appID    string
	mchId    string
	apiKey   string
}

func NewWxPayHandler(or *model.OrderRepository, appID, mchId, apiKey string) *WxPayHandler {
	return &WxPayHandler{orderRepo: or, appId: appID, mchId: mchId, apiKey: apiKey}
}

type UnifiedOrderRequest struct {
	OrderID    int64  `json:"order_id"`
	OpenID    string `json:"open_id"`
	Amount    int    `json:"amount"` // 单位：分
	Desc      string `json:"description"`
}

type UnifiedOrderResponse struct {
	PrepayID string `json:"prepay_id"`
	CodeURL  string `json:"code_url"`
}

// 统一下单
func (h *WxPayHandler) UnifiedOrder(c *gin.Context) {
	var req UnifiedOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := h.orderRepo.FindByID(req.OrderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	nonceStr := fmt.Sprintf("%d", time.Now().UnixNano())
	totalFee := int(order.TotalAmount * 100)
	body := req.Desc
	if body == "" {
		body = "卫浴商品"
	}

	params := map[string]string{
		"appid":            h.appId,
		"mch_id":           h.mchId,
		"nonce_str":        nonceStr,
		"body":             body,
		"out_trade_no":     order.OrderNo,
		"total_fee":        strconv.Itoa(totalFee),
		"spbill_create_ip": c.ClientIP(),
		"notify_url":       "https://your-domain.com/api/v1/wxpay/notify",
		"trade_type":       "NATIVE",
	}
	if req.OpenID != "" {
		params["trade_type"] = "JSAPI"
		params["openid"] = req.OpenID
	}

	sign := h.sign(params)
	params["sign"] = sign

	xmlData, _ := xml.Marshal(mapToXML(params))
	resp, err := http.Post("https://api.mch.weixin.qq.com/pay/unifiedorder", "text/xml", bytes.NewBuffer(xmlData))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "wechat pay request failed"})
		return
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	var result map[string]string
	xml.Unmarshal(bodyBytes, &result)

	if result["return_code"] == "SUCCESS" && result["result_code"] == "SUCCESS" {
		c.JSON(http.StatusOK, gin.H{
			"prepay_id": result["prepay_id"],
			"code_url":  result["code_url"],
		})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result["err_code_des"]})
	}
}

// 支付回调
func (h *WxPayHandler) Notify(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)
	var result map[string]string
	xml.Unmarshal(body, &result)

	if result["return_code"] == "SUCCESS" && h.verifySign(result) {
		outTradeNo := result["out_trade_no"]
		transactionID := result["transaction_id"]

		// 更新订单为已支付
		order, err := h.orderRepo.FindByOrderNo(outTradeNo)
		if err == nil && order.Status == 1 {
			h.orderRepo.UpdateStatus(order.ID, 3) // 已完成
		}

		c.XML(http.StatusOK, map[string]string{"return_code": "SUCCESS", "return_msg": "OK"})
	} else {
		c.XML(http.StatusBadRequest, map[string]string{"return_code": "FAIL", "return_msg": "签名验证失败"})
	}
}

func (h *WxPayHandler) sign(params map[string]string) string {
	keys := make([]string, 0, len(params))
	for k := range params {
		if k != "sign" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	var str string
	for _, k := range keys {
		str += fmt.Sprintf("%s=%s&", k, params[k])
	}
	str += fmt.Sprintf("key=%s", h.apiKey)
	hash := md5.Sum([]byte(str))
	return fmt.Sprintf("%x", hash)
}

func (h *WxPayHandler) verifySign(params map[string]string) bool {
	sign := h.sign(params)
	return sign == params["sign"]
}

func mapToXML(params map[string]string) []byte {
	var sb strings.Builder
	sb.WriteString("<xml>")
	for k, v := range params {
		sb.WriteString(fmt.Sprintf("<%s><![CDATA[%s]]></%s>", k, v, k))
	}
	sb.WriteString("</xml>")
	return []byte(sb.String())
}
