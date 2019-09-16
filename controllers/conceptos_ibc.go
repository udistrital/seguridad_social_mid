package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/udistrital/seguridad_social_mid/models"

	"github.com/astaxie/beego"
	"github.com/manucorporat/try"
)

// ConceptosIbcController ...controlador de tipo beego.Controller
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

		alerta.Type = "success"
		alerta.Code = "1"

		construirHechos(v)
		err = RegistrarHechos(v)
		if err != nil {
			alerta.Body = err.Error()
			panic(alerta)
		}
		c.Data["json"] = alerta

	}).Catch(func(e try.E) {
		fmt.Println("Error en conceptos_ibc.ActualizarConceptos(): ", e.(models.Alert).Body)
		//beego.Error("Error en conceptos_ibc.ActualizarConceptos(): ", e.(models.Alert).Body)
		c.Data["json"] = e
	})

	c.ServeJSON()
}

// RegistrarHechos hace llamados recursivos para actualizar los hechos
func RegistrarHechos(nombreConceptos []models.Predicado) (err error) {
	var apiResponse interface{}
	for _, value := range nombreConceptos {
		err = sendJson("http://"+beego.AppConfig.String("rulerServicio")+"/predicado/"+strconv.Itoa(value.Id), "PUT", &apiResponse, value)
		if err != nil {
			return
		}
	}

	return
}

// construirHechos recorre el arreglo de predicados, se crea un hecho y luego modifica el nombre del objeto con ese hecho
func construirHechos(nombreConceptos []models.Predicado) {
	for i, value := range nombreConceptos {
		hecho := "concepto_ibc(" + value.Nombre + ","
		if value.Estado {
			hecho += " activo)."
		} else {
			hecho += " inactivo)."
		}
		nombreConceptos[i].Nombre = hecho
		nombreConceptos[i].Descripcion = nombreConceptos[i].DescripcionHecho
	}
}
