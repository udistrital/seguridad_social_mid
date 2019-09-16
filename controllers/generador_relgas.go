package controllers

import (
	"encoding/json"
	"strings"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/manucorporat/try"
	"github.com/udistrital/seguridad_social_mid/models"
)

// GeneradorRelgasController operations for Generador_relgas
type GeneradorRelgasController struct {
	beego.Controller
}

// URLMapping ...
func (c *GeneradorRelgasController) URLMapping() {
	c.Mapping("ObtenerHechosCalculo", c.ObtenerHechosCalculo)
}

// RegistrarNuevosHechos ...
// @Title RegistrarNuevosHechos
// @Description Actualiza los nuevos conceptos y resoluciones de seguridad social
// @Success 200 {object} models.Alert
// @Failure 404 error
// @router /RegistrarNuevosHechos [post]
func (c *GeneradorRelgasController) RegistrarNuevosHechos() {
	var conceptosAporte []models.ConceptoAporte
	try.This(func() {
		err := json.Unmarshal(c.Ctx.Input.RequestBody, &conceptosAporte)
		if err != nil {
			panic(err)
		}
		predicados, err := construirHechosConceptos(conceptosAporte)
		if err != nil {
			panic(err)
		}
		err = RegistrarHechos(predicados)
		if err != nil {
			panic(err)
		}
		fmt.Println("predicados: ", predicados)
		//beego.Info("predicados: ", predicados)
		c.Data["json"] = models.Alert{Type: "success", Code: "1", Body: nil}
	}).Catch(func(e try.E) {
		ImprimirError("error en ObtenerHechosCalculo()", e.(error))
		c.Data["json"] = models.Alert{Type: "error", Code: "000", Body: nil}
	})
	c.ServeJSON()
}

// construirHechosConceptos construye los predicados a partir de un arreglo de conceptosAporte
func construirHechosConceptos(conceptosAporte []models.ConceptoAporte) (predicados []models.Predicado, err error) {
	estado := "inactivo"
	for _, value := range conceptosAporte {
		// Predicado para el concepto_aporte
		predicado := models.Predicado{
			Id:            value.Id,
			Nombre:        "concepto_aporte(" + value.NombreAporte + "," + value.Porcentaje + "," + value.Nomina + "," + value.NombreResolucion + ").",
			Dominio:       models.Dominio{Id: 19},
			TipoPredicado: models.TipoPredicado{Id: 1},
		}
		predicados = append(predicados, predicado)

		// Predicado para resolucion_aporte
		predicado.Id = value.Resolucion.Id
		if value.Resolucion.Estado {
			estado = "activo"
		}
		predicado.Nombre = "resolucion_aporte(" + value.Resolucion.NombreAporte + "," + value.Resolucion.Resolucion + "," + value.Resolucion.Vigencia + "," + estado + ")."
		predicados = append(predicados, predicado)
	}
	return
}

// ObtenerHechosCalculo ...
// @Title ObtenerHechosCalculo
// @Description Obtiene los hechos de calculo para seguridad social
// @Success 200 {object} models.Alert
// @Failure 404 error
// @router /ObtenerHechosCalculo [get]
func (c *GeneradorRelgasController) ObtenerHechosCalculo() {

	try.This(func() {

		conceptosAporte, err := obtenerConceptosAporte()
		if err != nil {
			panic(err)
		}

		c.Data["json"] = models.Alert{Type: "success", Code: "1", Body: conceptosAporte}
	}).Catch(func(e try.E) {
		ImprimirError("error en ObtenerHechosCalculo()", e.(error))
		c.Data["json"] = models.Alert{Type: "error", Code: "000", Body: nil}
	})
	c.ServeJSON()
}

// obtenerResolucionAporte devuelve un arreglo de tipo models.ResolucionAporte armado a partir de los hechos resolucion_aporte
func obtenerResolucionAporte(tipoAporte string) (resolucionAporte models.ResolucionAporte, err error) {
	var predicado []models.Predicado
	activo := false
	err = getJson("http://"+beego.AppConfig.String("rulerServicio")+"/predicado?"+
		"limit=1&"+
		"query=Dominio.Nombre:SeguridadSocial,"+
		"Nombre__startswith:resolucion_aporte(,"+
		"Nombre__contains:"+tipoAporte, &predicado)
	if err != nil {
		return
	}

	hecho := strings.Split(predicado[0].Nombre, ",")
	for i := range hecho {
		hecho[i] = ajustarString(hecho[i])
	}

	if hecho[3] == "activo" {
		activo = true
	}
	resolucionAporte = models.ResolucionAporte{
		Id:           predicado[0].Id,
		NombreAporte: hecho[0],
		Resolucion:   hecho[1],
		Vigencia:     hecho[2],
		Estado:       activo,
	}

	return
}

// obtenerConceptosAporte devuelve un arreglo de tipo models.ResolucionAporte armado a partir de los hechos concepto_aporte
func obtenerConceptosAporte() (conceptosAporte []models.ConceptoAporte, err error) {
	var predicados []models.Predicado

	err = getJson("http://"+beego.AppConfig.String("rulerServicio")+"/predicado?"+
		"limit=0&"+
		"query=Dominio.Nombre:SeguridadSocial,Nombre__startswith:concepto_aporte(", &predicados)
	if err != nil {
		return
	}

	for _, value := range predicados {
		hecho := strings.Split(value.Nombre, ",")
		for i := range hecho {
			hecho[i] = ajustarString(hecho[i])
		}

		resolucionAporte, errResolucion := obtenerResolucionAporte(hecho[0])
		if err != nil {
			return conceptosAporte, errResolucion
		}
		conceptosAporte = append(conceptosAporte, models.ConceptoAporte{
			Id:               value.Id,
			NombreAporte:     hecho[0],
			Porcentaje:       hecho[1],
			Nomina:           hecho[2],
			NombreResolucion: hecho[3],
			Resolucion:       resolucionAporte,
		})
	}
	return
}

// ajustarString ajusta el string s que recibe como parametro y elimna alguna de los strings establecidos en replace
func ajustarString(s string) string {
	if strings.Contains(s, ").") {
		s = strings.Replace(s, ").", "", 1)
	} else if strings.Contains(s, "concepto_aporte(") {
		s = strings.Replace(s, "concepto_aporte(", "", 1)
	} else if strings.Contains(s, "resolucion_aporte(") {
		s = strings.Replace(s, "resolucion_aporte(", "", 1)
	}
	return s
}
