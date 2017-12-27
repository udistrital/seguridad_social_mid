package controllers

import (
	"strconv"

	"github.com/astaxie/beego"
)

// IncapacidadesController operations for Incapacidades
type IncapacidadesController struct {
	beego.Controller
}

// URLMapping ...
func (c *IncapacidadesController) URLMapping() {
	c.Mapping("GetPersonas", c.GetPersonas)
}

// GetPersonas ...
// @Title GetPersonas
// @Description obtiene todas las personas que pueden aplicar a cualquier nómina
// @Success 200 {object} interface{}
// @Failure 403
// @router / [get]
func (c *IncapacidadesController) GetPersonas() {
	var (
		nominasResult interface{}
		alerta        interface{}
	)
	funcionariosResult := make(map[string]interface{})
	if err := getJson("http://"+beego.AppConfig.String("titanServicio")+"/nomina", &nominasResult); err == nil {
		nominas := nominasResult.([]interface{})
		for i := range nominas {
			idNomina := strconv.Itoa(int(nominas[i].(map[string]interface{})["Id"].(float64)))
			if err = sendJson("http://"+beego.AppConfig.String("titanServicio")+"/funcionario_proveedor", "POST", &alerta, nominas[i]); err == nil {
				funcionariosResult[idNomina] = alerta
				c.Ctx.Output.SetStatus(200)
			} else {
				c.Ctx.Output.SetStatus(403)
				c.Data["json"] = alerta
			}
		}
		c.Data["json"] = funcionariosResult
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}
