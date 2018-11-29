package controllers

import (
	"encoding/json"

	"github.com/udistrital/ss_mid_api/models"

	"github.com/astaxie/beego"
	"github.com/manucorporat/try"
)

//  ConceptosIbcController ...controlador de tipo beego.Controller
type ConceptosIbcController struct {
	beego.Controller
}

// URLMapping ...
func (c *ConceptosIbcController) URLMapping() {
	c.Mapping("ActualizarConceptos", c.ActualizarConceptos)
}

// ActualizarConceptos ...
// @Title ActualizarConceptos
// @Description Actualiza los hechos conceptos_ibc en el ruler
// @Param	body		body 	models.ConceptosIbc	true		"body for ConceptosIbc content"
// @Success 201 {int} models.Alert
// @Failure 404 body is empty
// @router /ActualizarConceptos/ [post]
func (c *ConceptosIbcController) ActualizarConceptos() {
	var v []models.Predicado
	alerta := models.Alert{
		Type: "error",
		Code: "0",
		Body: nil,
	}
	try.This(func() {
		err := json.Unmarshal(c.Ctx.Input.RequestBody, &v)
		if err != nil {
			alerta.Body = err.Error()
			panic(alerta)
		}
		beego.Info("v: ", v)
		alerta.Type = "success"
		alerta.Code = "1"

		//beego.Info("nombreConceptos: ", nombreConceptos)
		hechos := construirHechos(v)
		beego.Info("hechos construidos: ", hechos)
		c.Data["json"] = alerta
		//c.Data["json"] = ":D"

	}).Catch(func(e try.E) {
		beego.Error("Error en conceptos_ibc.ActualizarConceptos(): ", e.(models.Alert).Body)
		c.Data["json"] = e
	})

	c.ServeJSON()
}

// construirHechos contruye hechos para conceptos_ibc()
func construirHechos(nombreConceptos []models.Predicado) (hechosContruidos []string) {
	for _, value := range nombreConceptos {
		hecho := "concepto_ibc(" + value.Nombre + ","
		if value.Estado {
			hecho += " activo)."
		} else {
			hecho += " inactivo)."
		}
		hechosContruidos = append(hechosContruidos, hecho)
	}
	return
}
