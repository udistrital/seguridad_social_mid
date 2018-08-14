package controllers

import (
	"github.com/astaxie/beego"
	"github.com/manucorporat/try"
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
// @Param	documento		query	string false		"documento de la persona"
// @Success 200 {object} interface{}
// @Failure 403
// @router / [get]
func (c *IncapacidadesController) GetPersonas() {
	var personasNaturales []map[string]interface{}
	documento := c.GetString("documento")
	try.This(func() {
		beego.Info(documento)
		if err := getJson("http://"+beego.AppConfig.String("administrativaService")+"/informacion_proveedor?"+
			"limit=10&query=NumDocumento__icontains:"+documento+",TipoPersona:NATURAL", &personasNaturales); err == nil {
			// beego.Info("personas: ", personasNaturales)
			c.Data["json"] = personasNaturales
		} else {
			panic(err)
		}
	}).Catch(func(e try.E) {
		beego.Error("Error en GetPersonas() ", e)
		c.Data["json"] = e
	})
	c.ServeJSON()
}
