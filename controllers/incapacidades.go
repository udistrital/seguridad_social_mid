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
	var proveedores, personaNatural, respuesta []map[string]interface{}
	documento := c.GetString("documento")
	try.This(func() {
		if err := getJson("http://"+beego.AppConfig.String("administrativaService")+"/informacion_proveedor?"+
			"limit=6&query=NumDocumento__icontains:"+documento+",TipoPersona:NATURAL", &proveedores); err != nil {
			panic(err)
		}

		for index, proveedor := range proveedores {
			if err := getJson("http://"+beego.AppConfig.String("administrativaService")+"/informacion_persona_natural?"+
				"limit=1&query=Id:"+proveedor["NumDocumento"].(string), &personaNatural); err != nil {
				panic(err)
			} else {
				resp := map[string]interface{}{
					"display":       proveedor["NomProveedor"],
					"value":         proveedor["NumDocumento"],
					"documento":     proveedor["NumDocumento"],
					"id":            proveedor["Id"],
					"nominas":       index,
					"tipoDocumento": personaNatural[0]["TipoDocumento"].(map[string]interface{})["Abreviatura"],
				}

				respuesta = append(respuesta, resp)
			}
		}
		c.Data["json"] = respuesta

	}).Catch(func(e try.E) {
		beego.Error("Error en GetPersonas() ", e)
		c.Data["json"] = e
	})
	c.ServeJSON()
}
