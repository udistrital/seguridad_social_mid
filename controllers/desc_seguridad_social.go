package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/udistrital/ss_mid_api/golog"
	"github.com/udistrital/ss_mid_api/models"

	"github.com/astaxie/beego"
)

// DescSeguridadSocialController oprations for DescSeguridadSocial
type DescSeguridadSocialController struct {
	beego.Controller
}

// URLMapping ...
func (c *DescSeguridadSocialController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
	c.Mapping("CalcularSegSocial", c.CalcularSegSocial)
	c.Mapping("GetConceptosIbc", c.GetConceptosIbc)
	c.Mapping("GenerarPlanillaActivos", c.GenerarPlanillaActivos)
	c.Mapping("GenerarPlanillaPensionados", c.GenerarPlanillaPensionados)
}

// Post ...
// @Title Post
// @Description create DescSeguridadSocial
// @Param	body		body 	models.DescSeguridadSocial	true		"body for DescSeguridadSocial content"
// @Success 201 {int} models.DescSeguridadSocial
// @Failure 403 body is empty
// @router / [post]
func (c *DescSeguridadSocialController) Post() {
	var v models.DescSeguridadSocial
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if _, err := models.AddDescSeguridadSocial(&v); err == nil {
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

func (c *DescSeguridadSocialController) GetConceptosIbc() {
	fmt.Println("GetConceptosIbc")
	var predicados []models.Predicado
	var conceptos []models.Concepto
	var conceptosIbc []models.ConceptosIbc
	err := getJson("http://"+beego.AppConfig.String("rulerServicio")+
		"/predicado?limit=-1&query=Nombre__startswith:conceptos_ibc,Dominio.Id:4", &predicados)
	errConceptoTitan := getJson("http://"+beego.AppConfig.String("titanServicio")+
		"/concepto?limit=-1", &conceptos)
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

func (c *DescSeguridadSocialController) CalcularSegSocial() {
	idStr := c.Ctx.Input.Param(":id")
	_, err := strconv.Atoi(idStr)
	var alertas []string
	var predicado []models.Predicado
	var buffer bytes.Buffer //objeto para concatenar strings a la variable errores

	var detalleLiquidacion []models.DetalleLiquidacion
	var pagosSeguridadSocial []models.PagosSeguridadSocial

	if err != nil {
		c.Data["json"] = err.Error()
		fmt.Println("Error en detalle liquidacion: ", err)
	} else {

		err := getJson("http://"+beego.AppConfig.String("titanServicio")+"/detalle_liquidacion"+
			"?limit=0&query=Liquidacion.Id:"+idStr+",Concepto.NombreConcepto:ibc_liquidado", &detalleLiquidacion)

		fmt.Println(buffer.String())

		if err != nil {
			fmt.Println("ERROR EN LA PETICION:\n", buffer.String())
			alertas = append(alertas, "error al traer detalle liquidacion")
			c.Data["json"] = alertas
		} else {

			for index := 0; index < len(detalleLiquidacion); index++ {
				predicado = append(predicado, models.Predicado{Nombre: "ibc(" + strconv.Itoa(detalleLiquidacion[index].Persona) + "," + strconv.Itoa(int(detalleLiquidacion[index].ValorCalculado)) + ", salud)."})
				predicado = append(predicado, models.Predicado{Nombre: "ibc(" + strconv.Itoa(detalleLiquidacion[index].Persona) + "," + strconv.Itoa(int(detalleLiquidacion[index].ValorCalculado)) + ", riesgos)."})
				predicado = append(predicado, models.Predicado{Nombre: "ibc(" + strconv.Itoa(detalleLiquidacion[index].Persona) + "," + strconv.Itoa(int(detalleLiquidacion[index].ValorCalculado)) + ", apf)."})
			}

			reglas := CargarReglasBase() + FormatoReglas(predicado) + CargarNovedades(idStr) +
				ValorSaludEmpleado(idStr) + ValorPensionEmpleado(idStr)

			ids := golog.GetInt64(reglas, "v_salud_ud(I,Y).", "I")
			saludUd := golog.GetFloat(reglas, "v_salud_ud(I,Y).", "Y")
			saludTotal := golog.GetInt64(reglas, "v_total_salud(X,T).", "T")
			pensionUd := golog.GetFloat(reglas, "v_pen_ud(I,Y).", "Y")
			pensionTotal := golog.GetInt64(reglas, "v_total_pen(X,T).", "T")
			arl := golog.GetInt64(reglas, "v_arl(I,Y).", "Y")

			for index := 0; index < len(ids); index++ {
				aux := models.PagosSeguridadSocial{
					Persona:              ids[index],
					SaludUd:              saludUd[index],
					SaludTotal:           saludTotal[index],
					PensionUd:            pensionUd[index],
					PensionTotal:         pensionTotal[index],
					IdDetalleLiquidacion: detalleLiquidacion[index].Id,
					Arl:                  arl[index]}

				pagosSeguridadSocial = append(pagosSeguridadSocial, aux)
			}

			c.Data["json"] = pagosSeguridadSocial
		}
		c.ServeJSON()
	}
}

// AproximarValor ...
// @Title Aproximar Valor
// @Description Aproxima un valor al número de aproximacion más cercano
// @Param valorInicial valor a cual quiere aproximar
// @Param	aproximacion	multiplo sobre al cual debe aproximarse el valor inicial
// @return valorAproximado el valor aproximado de acuerdo a la aproximacion
func AproximarValor(valorInicial int64, aproximacion int64) (valorAproximado int64) {
	x := float64(valorInicial) / float64(aproximacion)
	y := math.Trunc(float64(valorInicial / aproximacion))
	if (x - y) > 0 {
		valorAproximado = AproximarValor(valorInicial+1, aproximacion)
	} else {
		valorAproximado = valorInicial
	}
	return
}

// ValorSaludEmpleado ...
// @Title Valor Salud Empleado
// @Description Crea todos los hechos con la información del valor de salud
// @Param	idLiquidacion		id de la liquidacion correspondiente
func ValorSaludEmpleado(idLiquidacion string) (valorSaludEmpleado string) {
	var detalleLiquSalud []models.DetalleLiquidacion
	var predicado []models.Predicado

	errSalud := getJson("http://"+beego.AppConfig.String("titanServicio")+"/detalle_liquidacion"+
		"?limit=0&query=Liquidacion.Id:"+idLiquidacion+",Concepto.NombreConcepto:salud", &detalleLiquSalud)

	if errSalud != nil {
		fmt.Println("Error en ValorSaludEmpleado:\n", errSalud)
	} else {
		for index := 0; index < len(detalleLiquSalud); index++ {
			predicado = append(predicado, models.Predicado{Nombre: "v_salud_func(" + strconv.Itoa(detalleLiquSalud[index].Persona) + ", " + strconv.Itoa(int(detalleLiquSalud[index].ValorCalculado)) + ")."})
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
	var detalleLiquPension []models.DetalleLiquidacion
	var predicado []models.Predicado

	errPension := getJson("http://"+beego.AppConfig.String("titanServicio")+"/detalle_liquidacion"+
		"?limit=0&query=Liquidacion.Id:"+idLiquidacion+",Concepto.NombreConcepto:pension", &detalleLiquPension)

	if errPension != nil {
		fmt.Println("Error en ValorPensionEmpleado:\n", errPension)
	} else {
		for index := 0; index < len(detalleLiquPension); index++ {
			predicado = append(predicado, models.Predicado{Nombre: "v_pen_func(" + strconv.Itoa(detalleLiquPension[index].Persona) + ", " + strconv.Itoa(int(detalleLiquPension[index].ValorCalculado)) + ")."})
			valorPensionEmpleado += predicado[index].Nombre + "\n"
		}
	}
	return
}

// CargarNovedades ...
// @Title CargarNovedades
// @Description obtiene todas los conceptos con naturaleza seguridad_social desde detalle_liquidacion por el id
// @Param	id		path 	string	true		"The key for staticblock"
func CargarNovedades(id string) (novedades string) {
	var conceptos []models.DetalleLiquidacion
	var predicado []models.Predicado

	errLincNo := getJson("http://"+beego.AppConfig.String("titanServicio")+"/detalle_liquidacion?limit=0&query=Liquidacion.Id:"+id+",Concepto.naturaleza:seguridad_social&fields=Concepto,Persona", &conceptos)

	if errLincNo != nil {
		fmt.Println("error en CargarNovedades()", errLincNo)
	} else {
		for index := 0; index < len(conceptos); index++ {
			predicado = append(predicado, models.Predicado{Nombre: "novedad_persona(" + conceptos[index].Concepto.NombreConcepto + ", " + strconv.Itoa(int(conceptos[index].Persona)) + ")."})
			novedades += predicado[index].Nombre + "\n"
		}
	}
	return
}

/*func (c *DescSeguridadSocialController) GenerarArchivoPlano() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	GetA

}*/

// GetOne ...
// @Title Get One
// @Description get DescSeguridadSocial by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.DescSeguridadSocial
// @Failure 403 :id is empty
// @router /:id [get]
func (c *DescSeguridadSocialController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetDescSeguridadSocialById(id)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = v
	}
	c.ServeJSON()
}

// GetAll ...
// @Title Get All
// @Description get DescSeguridadSocial
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.DescSeguridadSocial
// @Failure 403
// @router / [get]
func (c *DescSeguridadSocialController) GetAll() {
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

	l, err := models.GetAllDescSeguridadSocial(query, fields, sortby, order, offset, limit)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = l
	}
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the DescSeguridadSocial
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.DescSeguridadSocial	true		"body for DescSeguridadSocial content"
// @Success 200 {object} models.DescSeguridadSocial
// @Failure 403 :id is not int
// @router /:id [put]
func (c *DescSeguridadSocialController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.DescSeguridadSocial{Id: id}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if err := models.UpdateDescSeguridadSocialById(&v); err == nil {
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
// @Description delete the DescSeguridadSocial
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *DescSeguridadSocialController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	if err := models.DeleteDescSeguridadSocial(id); err == nil {
		c.Data["json"] = "OK"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

func (c *DescSeguridadSocialController) GenerarPlanillaActivos() {
	var proveedores []models.InformacionProveedor
	//var pagosSalud []models.DescSeguridadSocialDetalle
	//var pagoSalud []models.DescSeguridadSocial
	var detalleLiquidacion []models.DetalleLiquidacion
	var errStrings []string
	tipoRegistro := "02"
	fila := ""

	errLiquidacion := getJson("http://"+beego.AppConfig.String("titanServicio")+"/detalle_liquidacion"+
		"?limit=-1", &detalleLiquidacion)
	if errLiquidacion != nil {
		errStrings = append(errStrings, errLiquidacion.Error())
	}

	errProveedores := getJson("http://"+beego.AppConfig.String("agoraServicio")+"/informacion_proveedor"+
		"?limit=0", &proveedores)
	if errProveedores != nil {
		errStrings = append(errStrings, errProveedores.Error())
	}

	fmt.Println("errStrings: ", errStrings)
	if errStrings == nil {
		secuencia := 1
		x := 1
		for i := 0; i < len(proveedores); i++ {
			for j := 0; j < len(detalleLiquidacion); j++ {
				if proveedores[i].Id == detalleLiquidacion[j].Persona {
					fila += formatoDato(tipoRegistro, 2)
					fila += formatoDato(completarSecuencia(secuencia), 5)
					fila += formatoDato("CC", 2)
					fila += formatoDato(strconv.Itoa(int(proveedores[j].NumDocumento)), 16)
					fila += formatoDato("1", 2)
					fila += formatoDato("1", 2)
					fila += formatoDato(" ", 1)
					fila += formatoDato(detalleLiquidacion[j].Concepto, 1)
					if x == 1 {
						fmt.Printf("Tamaño fila : %d\n", len(fila))
						x = 2
					}

					fila += "\n"
					secuencia++
				}
			}
		}
		fmt.Println("Filas:\n", fila)
	}
}

func (c *DescSeguridadSocialController) GenerarPlanillaPensionados() {
	var proveedores []models.InformacionProveedor
	var personasNatural []models.InformacionPersonaNatural
	//var pagosSalud []models.DescSeguridadSocialDetalle
	//var pagoSalud []models.DescSeguridadSocial
	var detalleLiquidacion []models.DetalleLiquidacion
	var errStrings []string
	tipoRegistro := "02"
	fila := ""

	errLiquidacion := getJson("http://"+beego.AppConfig.String("titanServicio")+"/detalle_liquidacion"+
		"?limit=-1", &detalleLiquidacion)
	if errLiquidacion != nil {
		errStrings = append(errStrings, errLiquidacion.Error())
	}

	errProveedores := getJson("http://"+beego.AppConfig.String("titanServicio")+"/informacion_proveedor"+
		"?limit=-1", &proveedores)
	if errProveedores != nil {
		errStrings = append(errStrings, errProveedores.Error())
	}

	errPersonaNatural := getJson("http://"+beego.AppConfig.String("titanServicio")+"/informacion_persona_natural"+
		"?limit=0", &personasNatural)
	if errPersonaNatural != nil {
		errStrings = append(errStrings, errPersonaNatural.Error())
	}

	fmt.Println("**errStrings: ", errStrings)
	if errStrings == nil {
		secuencia := 1
		x := 1
		for i := 0; i < len(proveedores); i++ {
			for j := 0; j < len(detalleLiquidacion); j++ {
				for k := 0; k < len(personasNatural); k++ {
					if proveedores[i].Id == detalleLiquidacion[j].Persona {
						if int(proveedores[i].NumDocumento) == personasNatural[k].Id {
							fila += formatoDato(tipoRegistro, 2)
							fila += formatoDato(completarSecuencia(secuencia), 5)
							fila += formatoDato(personasNatural[k].PrimerApellido, 20)
							fila += formatoDato(personasNatural[k].SegundoApellido, 30)
							fila += formatoDato(personasNatural[k].PrimerNombre, 20)
							fila += formatoDato(personasNatural[k].SegundoApellido, 30)
							fila += formatoDato("CC", 2)
							fila += formatoDato(strconv.Itoa(int(personasNatural[k].Id)), 16)
							if x == 1 {
								fmt.Printf("Tamaño fila : %d\n", len(fila))
								x = 2
							}
							fila += "\n"
							secuencia++
						}
					}
				}
			}
		}
		fmt.Println("Filas:\n", fila)
	}
}

func completarSecuencia(num int) (secuencia string) {
	tamanioNum := len(strconv.Itoa(num))
	for i := 0; i < 5-tamanioNum; i++ {
		secuencia += "0"
	}
	secuencia += strconv.Itoa(num)
	return
}

func formatoDato(texto string, longitud int) (textoEscribir string) {
	for _, r := range texto {
		textoEscribir += string(r)
	}
	for i := 0; i < longitud-len(texto); i++ {
		textoEscribir += " "
	}
	return
}
