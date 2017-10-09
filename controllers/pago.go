package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/udistrital/ss_mid_api/golog"
	"github.com/udistrital/ss_mid_api/models"
)

// PagoController operations for Pago
type PagoController struct {
	beego.Controller
}

// URLMapping ...
func (c *PagoController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
	c.Mapping("CalcularSegSocial", c.CalcularSegSocial)
	c.Mapping("ConceptosIbc", c.ConceptosIbc)
	c.Mapping("GetNovedadesPorPersona", c.NovedadesPorPersona)
}

func (c *PagoController) ConceptosIbc() {
	fmt.Println("GetConceptosIbc")
	var predicados []models.Predicado
	var conceptos []models.Concepto
	var conceptosIbc []models.ConceptosIbc
	err := getJson("http://"+beego.AppConfig.String("rulerServicio")+
		"/predicado?limit=-1&query=Nombre__startswith:conceptos_ibc,Dominio.Id:4", &predicados)
	errConceptoTitan := getJson("http://"+beego.AppConfig.String("titanServicio")+
		"/concepto_nomina?limit=-1", &conceptos)
	if err != nil && errConceptoTitan != nil {
		c.Data["json"] = err.Error() + errConceptoTitan.Error()
	} else {
		nombres := golog.GetString(FormatoReglas(predicados), "conceptos_ibc(X).", "X")
		for i := 0; i < len(predicados); i++ {
			for j := 0; j < len(conceptos); j++ {
				if nombres[i] == conceptos[j].NombreConcepto {
					aux := models.ConceptosIbc{
						Id:          predicados[i].Id,
						Nombre:      nombres[i],
						Descripcion: conceptos[j].AliasConcepto}
					conceptosIbc = append(conceptosIbc, aux)
					break
				}
			}

		}
		c.Data["json"] = conceptosIbc
	}
	c.ServeJSON()
}

func (c *PagoController) NovedadesPorPersona() {
	personaStr := c.Ctx.Input.Param(":persona")
	_, err := strconv.Atoi(personaStr)
	var alertas []string
	var errores string
	var concepto models.ConceptoNominaPorPersona
	var detallePreliquidacion []models.DetallePreliquidacion
	var conceptoPorPersona []models.ConceptoNominaPorPersona
	var novedadesPersonaSs []models.NovedadesPersonaSS

	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		err := getJson("http://"+beego.AppConfig.String("titanServicio")+
			"/detalle_preliquidacion"+
			"?limit=0"+
			"&query=Persona:"+personaStr+",Concepto.Naturaleza:seguridad_social"+
			"&fields=Concepto,ValorCalculado", &detallePreliquidacion)
		errConcepto := getJson("http://"+beego.AppConfig.String("titanServicio")+
			"concepto_nomina_por_persona"+
			"?limit=0"+
			"&query=Persona:"+personaStr+",Concepto.Naturaleza:seguridad_social,EstadoNovedad:Activo"+
			"&fields=Concepto,Persona,FechaDesde,FechaHasta", &conceptoPorPersona)
		errores = ""
		fmt.Println(err, errConcepto)

		if errores != "" {
			fmt.Println("err en la peticion", errores)
			alertas = append(alertas, "error al traer detalle liquidacion")
			fmt.Println(alertas)
			c.Data["json"] = alertas
		} else {
			persona, _ := strconv.ParseInt(personaStr, 10, 64)
			for index := 0; index < len(detallePreliquidacion); index++ {
				if detallePreliquidacion[index].Concepto == conceptoPorPersona[index].Concepto {
					concepto = conceptoPorPersona[index]
					fmt.Println(concepto)
				}
				novedadesPersonaSs = append(novedadesPersonaSs, models.NovedadesPersonaSS{
					Persona:                  persona,
					Preliquidacion:           detallePreliquidacion[index],
					ConceptoNominaPorPersona: concepto})
			}

			c.Data["json"] = novedadesPersonaSs
		}
		c.ServeJSON()
	}
}

func (c *PagoController) CalcularSegSocial() {
	idStr := c.Ctx.Input.Param(":id")
	_, err := strconv.Atoi(idStr)
	var alertas []string
	var predicado []models.Predicado
	var buffer bytes.Buffer //objeto para concatenar strings a la variable errores

	var detallePreliquidacion []models.DetallePreliquidacion
	var pagosSeguridadSocial []models.PagosSeguridadSocial

	if err != nil {
		c.Data["json"] = err.Error()
		fmt.Println("Error en detalle liquidacion: ", err)
	} else {

		err := getJson("http://"+beego.AppConfig.String("titanServicio")+"/detalle_preliquidacion"+
			"?limit=0&query=Preliquidacion.Id:"+idStr+",Concepto.NombreConcepto:ibc_liquidado", &detallePreliquidacion)

		fmt.Println(buffer.String())

		if err != nil {
			fmt.Println("ERROR EN LA PETICION:\n", buffer.String())
			alertas = append(alertas, "error al traer detalle liquidacion")
			c.Data["json"] = alertas
		} else {

			for index := 0; index < len(detallePreliquidacion); index++ {
				predicado = append(predicado, models.Predicado{Nombre: "ibc(" + detallePreliquidacion[index].NumeroContrato + "," + strconv.Itoa(int(detallePreliquidacion[index].ValorCalculado)) + ", salud)."})
				predicado = append(predicado, models.Predicado{Nombre: "ibc(" + detallePreliquidacion[index].NumeroContrato + "," + strconv.Itoa(int(detallePreliquidacion[index].ValorCalculado)) + ", riesgos)."})
				predicado = append(predicado, models.Predicado{Nombre: "ibc(" + detallePreliquidacion[index].NumeroContrato + "," + strconv.Itoa(int(detallePreliquidacion[index].ValorCalculado)) + ", apf)."})
			}

			reglas := CargarReglasBase() + FormatoReglas(predicado) + CargarNovedades(idStr) +
				ValorSaludEmpleado(idStr) + ValorPensionEmpleado(idStr)

			fmt.Println("http://" + beego.AppConfig.String("titanServicio") + "/detalle_preliquidacion" +
				"?limit=0&query=Preliquidacion.Id:" + idStr + ",Concepto.NombreConcepto:ibc_liquidado")

			numContrato := golog.GetString(reglas, "v_salud_ud(I,Y).", "I")
			saludUd := golog.GetFloat(reglas, "v_salud_ud(I,Y).", "Y")
			saludTotal := golog.GetInt64(reglas, "v_total_salud(X,T).", "T")
			pensionUd := golog.GetFloat(reglas, "v_pen_ud(I,Y).", "Y")
			pensionTotal := golog.GetInt64(reglas, "v_total_pen(X,T).", "T")
			arl := golog.GetInt64(reglas, "v_arl(I,Y).", "Y")
			caja := golog.GetInt64(reglas, "v_caja(I,Y).", "Y")
			icbf := golog.GetInt64(reglas, "v_icbf(I,Y).", "Y")

			for index := 0; index < len(numContrato); index++ {
				aux := models.PagosSeguridadSocial{
					NumeroContrato: numContrato[index],
					SaludUd:        saludUd[index],
					SaludTotal:     saludTotal[index],
					PensionUd:      pensionUd[index],
					PensionTotal:   pensionTotal[index],
					Caja:           caja[index],
					Icbf:           icbf[index],
					IdDetallePreliquidacion: detallePreliquidacion[index].Id,
					Arl: arl[index]}

				pagosSeguridadSocial = append(pagosSeguridadSocial, aux)
			}

			c.Data["json"] = pagosSeguridadSocial
		}
		c.ServeJSON()
	}
}

// ValorSaludEmpleado ...
// @Title Valor Salud Empleado
// @Description Crea todos los hechos con la información del valor de salud
// @Param	idLiquidacion		id de la liquidacion correspondiente
func ValorSaludEmpleado(idLiquidacion string) (valorSaludEmpleado string) {
	var detalleLiquSalud []models.DetallePreliquidacion
	var predicado []models.Predicado

	errSalud := getJson("http://"+beego.AppConfig.String("titanServicio")+"/detalle_preliquidacion"+
		"?limit=0&query=Preliquidacion:"+idLiquidacion+",Concepto.NombreConcepto:salud", &detalleLiquSalud)

	if errSalud != nil {
		fmt.Println("Error en ValorSaludEmpleado:\n", errSalud)
	} else {
		for index := 0; index < len(detalleLiquSalud); index++ {
			predicado = append(predicado, models.Predicado{Nombre: "v_salud_func(" + detalleLiquSalud[index].NumeroContrato + ", " + strconv.Itoa(int(detalleLiquSalud[index].ValorCalculado)) + ")."})
			valorSaludEmpleado += predicado[index].Nombre + "\n"
		}
	}
	return
}

// ValorPensionEmpleado ...
// @Title Valor Pensión Empleado
// @Description Crea todos los hechos con la información del valor de la pensión
// @Param	idLiquidacion		id de la liquidacion correspondiente
func ValorPensionEmpleado(idLiquidacion string) (valorPensionEmpleado string) {
	var detalleLiquPension []models.DetallePreliquidacion
	var predicado []models.Predicado

	errPension := getJson("http://"+beego.AppConfig.String("titanServicio")+"/detalle_preliquidacion"+
		"?limit=0&query=Preliquidacion:"+idLiquidacion+",Concepto.NombreConcepto:pension", &detalleLiquPension)

	if errPension != nil {
		fmt.Println("Error en ValorPensionEmpleado:\n", errPension)
	} else {
		for index := 0; index < len(detalleLiquPension); index++ {
			predicado = append(predicado, models.Predicado{Nombre: "v_pen_func(" + detalleLiquPension[index].NumeroContrato + ", " + strconv.Itoa(int(detalleLiquPension[index].ValorCalculado)) + ")."})
			valorPensionEmpleado += predicado[index].Nombre + "\n"
		}
	}
	return
}

// CargarNovedades ...
// @Title CargarNovedades
// @Description obtiene todas los conceptos con naturaleza seguridad_social desde detalle_preliquidacion por el id
// @Param	id		path 	string	true		"The key for staticblock"
func CargarNovedades(id string) (novedades string) {
	var conceptosPreliquidacion []models.DetallePreliquidacion
	var predicado []models.Predicado

	errLincNo := getJson("http://"+beego.AppConfig.String("titanServicio")+"/detalle_preliquidacion"+
		"?limit=0&query=Preliquidacion:"+id+",Concepto.NaturalezaConcepto.Nombre:seguridad_social&fields=Concepto,NumeroContrato", &conceptosPreliquidacion)

	if errLincNo != nil {
		fmt.Println("error en CargarNovedades()", errLincNo)
	} else {
		for index := 0; index < len(conceptosPreliquidacion); index++ {
			predicado = append(predicado, models.Predicado{Nombre: "novedad_persona(" + conceptosPreliquidacion[index].Concepto.NombreConcepto + ", " + conceptosPreliquidacion[index].NumeroContrato + ")."})
			novedades += predicado[index].Nombre + "\n"
		}
	}
	return
}

// Post ...
// @Title Post
// @Description create Pago
// @Param	body		body 	models.Pago	true		"body for Pago content"
// @Success 201 {int} models.Pago
// @Failure 403 body is empty
// @router / [post]
func (c *PagoController) Post() {
	var v models.Pago
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if _, err := models.AddPago(&v); err == nil {
			c.Ctx.Output.SetStatus(201)
			c.Data["json"] = v
		} else {
			c.Data["json"] = err.Error()
		}
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// GetOne ...
// @Title Get One
// @Description get Pago by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Pago
// @Failure 403 :id is empty
// @router /:id [get]
func (c *PagoController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetPagoById(id)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = v
	}
	c.ServeJSON()
}

// GetAll ...
// @Title Get All
// @Description get Pago
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Pago
// @Failure 403
// @router / [get]
func (c *PagoController) GetAll() {
	var fields []string
	var sortby []string
	var order []string
	var query = make(map[string]string)
	var limit int64 = 10
	var offset int64

	// fields: col1,col2,entity.col3
	if v := c.GetString("fields"); v != "" {
		fields = strings.Split(v, ",")
	}
	// limit: 10 (default is 10)
	if v, err := c.GetInt64("limit"); err == nil {
		limit = v
	}
	// offset: 0 (default is 0)
	if v, err := c.GetInt64("offset"); err == nil {
		offset = v
	}
	// sortby: col1,col2
	if v := c.GetString("sortby"); v != "" {
		sortby = strings.Split(v, ",")
	}
	// order: desc,asc
	if v := c.GetString("order"); v != "" {
		order = strings.Split(v, ",")
	}
	// query: k:v,k:v
	if v := c.GetString("query"); v != "" {
		for _, cond := range strings.Split(v, ",") {
			kv := strings.SplitN(cond, ":", 2)
			if len(kv) != 2 {
				c.Data["json"] = errors.New("Error: invalid query key/value pair")
				c.ServeJSON()
				return
			}
			k, v := kv[0], kv[1]
			query[k] = v
		}
	}

	l, err := models.GetAllPago(query, fields, sortby, order, offset, limit)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = l
	}
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the Pago
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Pago	true		"body for Pago content"
// @Success 200 {object} models.Pago
// @Failure 403 :id is not int
// @router /:id [put]
func (c *PagoController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.Pago{Id: id}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if err := models.UpdatePagoById(&v); err == nil {
			c.Data["json"] = "OK"
		} else {
			c.Data["json"] = err.Error()
		}
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the Pago
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *PagoController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	if err := models.DeletePago(id); err == nil {
		c.Data["json"] = "OK"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}
