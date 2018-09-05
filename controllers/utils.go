package controllers

import (
	"time"

	"github.com/astaxie/beego"
)

// UtilsController operations for Beneficiarios
type UtilsController struct {
	beego.Controller
}

// URLMapping ...
func (c *UtilsController) URLMapping() {
	c.Mapping("GetActualDate", c.GetActualDate)
}

// GetActualTime ...
// @Title GetActualTime
// @Description obtiene la fecha actual del servidor
// @Param	body		body 	models.Beneficiarios	true		"body for Beneficiarios content"
// @Success 201 {string} fecha_actual
// @Failure 403 body is empty
// @router /GetActualDate [get]
func (c *UtilsController) GetActualDate() {
	t := time.Now()
	c.Data["json"] = t.Format("2006-01-02")
	c.ServeJSON()
}