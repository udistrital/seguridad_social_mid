package controllers

import (
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
	c.Mapping("SumarPagosSalud", c.SumarPagosSalud)
	c.Mapping("RegistrarPagos", c.RegistrarPagos)
}

// SumarPagosSalud ...
// @Title Sumar pagos de salid
// @Description Suma el total de los pagos de salud y pensión de ud
// con el total de pagos de salud y pensión del empleado
// @Param	idPeriodoPago		id del periodo pago de seguridad social
// @router SumarPagosSalud/:idPeriodoPago [get]
func (c *PagoController) SumarPagosSalud() {
	idStr := c.Ctx.Input.Param(":idPeriodoPago")
	//id, _ := strconv.Atoi(idStr)
	var pagosUd []models.TotalPagosUd
	errTotalPagos := getJson("http://"+beego.AppConfig.String("segSocialService")+
		"/pago/GetPagos/"+idStr, &pagosUd)

	if errTotalPagos != nil {
		c.Data["json"] = errTotalPagos.Error()
	} else {
		/*for i := 0; i < len(errTotalPagos); i++ {

		}*/
		c.Data["json"] = pagosUd
	}
	c.ServeJSON()
}

func (c *PagoController) ConceptosIbc() {
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
			alertas = append(alertas, "error al traer detalle liquidacion")
			c.Data["json"] = alertas
		} else {
			persona, _ := strconv.ParseInt(personaStr, 10, 64)
			for index := 0; index < len(detallePreliquidacion); index++ {
				if detallePreliquidacion[index].Concepto == conceptoPorPersona[index].Concepto {
					concepto = conceptoPorPersona[index]
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
	var (
		alertas, contratos []string
		predicado          []models.Predicado
		//buffer                bytes.Buffer //objeto para concatenar strings a la variable errores
		detallePreliquidacion []models.DetallePreliquidacion
		pagosSeguridadSocial  []*models.PagosSeguridadSocial
	)

	if err != nil {
		c.Data["json"] = err.Error()
	} else {

		err := getJson("http://"+beego.AppConfig.String("titanServicio")+"/detalle_preliquidacion"+
			"?limit=0&query=Preliquidacion.Id:"+idStr+",Concepto.NombreConcepto:salarioBase", &detallePreliquidacion)

		if err != nil {
			alertas = append(alertas, "error al traer detalle liquidacion")
			c.Data["json"] = alertas
		} else {
			idDetallePreliquidacion := detallePreliquidacion[0].Preliquidacion.Id

			for index := 0; index < len(detallePreliquidacion); index++ {
				predicado = append(predicado, models.Predicado{Nombre: "ibc(" + detallePreliquidacion[index].NumeroContrato + "," + strconv.Itoa(int(detallePreliquidacion[index].ValorCalculado)) + ", salud)."})
				predicado = append(predicado, models.Predicado{Nombre: "ibc(" + detallePreliquidacion[index].NumeroContrato + "," + strconv.Itoa(int(detallePreliquidacion[index].ValorCalculado)) + ", riesgos)."})
				predicado = append(predicado, models.Predicado{Nombre: "ibc(" + detallePreliquidacion[index].NumeroContrato + "," + strconv.Itoa(int(detallePreliquidacion[index].ValorCalculado)) + ", apf)."})
			}

			reglas := CargarReglasBase() + FormatoReglas(predicado) + CargarNovedades(idStr) +
				valorSaludEmpleado(idStr) + ValorPensionEmpleado(idStr)

			numContrato := golog.GetString(reglas, "v_salud_ud(I,Y).", "I")
			saludUd := golog.GetFloat(reglas, "v_salud_ud(I,Y).", "Y")
			saludTotal := golog.GetInt64(reglas, "v_total_salud(X,T).", "T")
			pensionUd := golog.GetFloat(reglas, "v_pen_ud(I,Y).", "Y")
			pensionTotal := golog.GetInt64(reglas, "v_total_pen(X,T).", "T")
			arl := golog.GetInt64(reglas, "v_arl(I,Y).", "Y")
			caja := golog.GetInt64(reglas, "v_caja(I,Y).", "Y")
			icbf := golog.GetInt64(reglas, "v_icbf(I,Y).", "Y")

			// Acá se debe cambiar como se arma el modelo
			for index := 0; index < len(numContrato); index++ {
				contratos = append(contratos, numContrato[index])
				aux := &models.PagosSeguridadSocial{
					NombrePersona:           "",
					NumeroContrato:          numContrato[index],
					SaludUd:                 saludUd[index],
					SaludTotal:              saludTotal[index],
					PensionUd:               pensionUd[index],
					PensionTotal:            pensionTotal[index],
					Caja:                    caja[index],
					Icbf:                    icbf[index],
					IdPreliquidacion:        idDetallePreliquidacion,
					IdDetallePreliquidacion: detallePreliquidacion[index].Id,
					Arl: arl[index]}

				pagosSeguridadSocial = append(pagosSeguridadSocial, aux)
			}

			mapProveedores, _ := getInfoProveedor(contratos)

			for i, _ := range pagosSeguridadSocial {
				pagosSeguridadSocial[i].NombrePersona = mapProveedores[pagosSeguridadSocial[i].NumeroContrato].NomProveedor
			}

			c.Data["json"] = pagosSeguridadSocial
		}
		c.ServeJSON()
	}
}

// valorSaludEmpleado ...
// @Title Valor Salud Empleado
// @Description Crea todos los hechos con la información del valor de salud
// @Param	idLiquidacion		id de la liquidacion correspondiente
func valorSaludEmpleado(idLiquidacion string) (valorSaludEmpleado string) {
	var detalleLiquSalud []models.DetallePreliquidacion
	var predicado []models.Predicado

	errSalud := getJson("http://"+beego.AppConfig.String("titanServicio")+"/detalle_preliquidacion"+
		"?limit=0&query=Preliquidacion:"+idLiquidacion+",Concepto.NombreConcepto:salud", &detalleLiquSalud)

	if errSalud != nil {
		return ""
	} else {
		for index := 0; index < len(detalleLiquSalud); index++ {
			predicado = append(predicado, models.Predicado{Nombre: "v_salud_func(" + detalleLiquSalud[index].NumeroContrato + ", " + strconv.Itoa(int(detalleLiquSalud[index].ValorCalculado)) + ")."})
			valorSaludEmpleado += predicado[index].Nombre + "\n"
		}
	}
	return
}

// SaludHCHonorarios
//@Title Valor correspondiente a salud de hora catedra honorarios
//@Description consulta una preliqudacion correspondiente a hora catedra
//@Param idLiquidacion id de la preliquidacion correspondiente
func SaludHCHonorarios(idLiquidacion string) (valorSaludEmpleado string) {
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

// RegistrarPagos ...
// @Title Registrar pagos de seguridad social
// @Description Recibe los pagos para registrar seguridad social, les agrega un estado,
// y a cada uno de los pagos les asocia la entidad correspondiente
// @Param	body body models.TrPeriodoPago true	"body for TrPeriodoPago"
// @Success 201 {int} models.PeriodoPago
// @Failure 403 body is empty
// @router /RegistrarPagos [post]
func (c *PagoController) RegistrarPagos() {
	var (
		PeriodoPago models.TrPeriodoPago
		alerta      interface{}
	)

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &PeriodoPago); err == nil {
		mapProveedores, err := getInfoProveedor(PeriodoPago.Contratos)
		if err != nil {
			c.Data["json"] = err.Error()
		}

		mapPersonas, err := getInfoPersona(mapProveedores)
		if err != nil {
			c.Data["json"] = err.Error()
		}

		pagosSeg, _ := getPagosSeg()
		contContratista := 0
		for i := range PeriodoPago.Pagos {
			nombrePago := pagosSeg[PeriodoPago.Pagos[i].TipoPago]
			switch nombrePago {
			case "arl":
				PeriodoPago.Pagos[i].EntidadPago = mapPersonas[PeriodoPago.Contratos[contContratista]].IdArl
			case "pension_ud":
				PeriodoPago.Pagos[i].EntidadPago = mapPersonas[PeriodoPago.Contratos[contContratista]].IdFondoPension
			case "salud_ud":
				PeriodoPago.Pagos[i].EntidadPago = mapPersonas[PeriodoPago.Contratos[contContratista]].IdEps
			case "caja_compensacion":
				PeriodoPago.Pagos[i].EntidadPago = mapPersonas[PeriodoPago.Contratos[contContratista]].IdCajaCompensacion
			default: // ICBF
				PeriodoPago.Pagos[i].EntidadPago = 0
			}
			contContratista++
			if contContratista == 5 {
				contContratista = 0
			}
		}

		if err = sendJson("http://"+beego.AppConfig.String("segSocialService")+"/tr_periodo_pago", "POST", &alerta, PeriodoPago); err == nil {
			c.Ctx.Output.SetStatus(201)
		}
		c.Data["json"] = alerta

	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// getInfoProveedor Recibe un arreglo de strings con los contratos y devuelve un map con la información del proveedor
func getInfoProveedor(contratos []string) (map[string]models.InformacionProveedor, error) {
	personas := make(map[string]models.InformacionProveedor)
	var (
		proveedor models.InformacionProveedor
		contrato  models.ContratoGeneral
	)
	for i := range contratos {
		if err := getJson("http://"+beego.AppConfig.String("argoServicio")+"/contrato_general/"+contratos[i], &contrato); err == nil {
			if err = getJson("http://"+beego.AppConfig.String("agoraServicio")+"/informacion_proveedor/"+strconv.Itoa(contrato.Contratista), &proveedor); err == nil {
				personas[contratos[i]] = proveedor
			} else {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	return personas, nil
}

// getInfoPersona recibe un map de proveedores para consultar el número de contrato y devuelve un mapa con la inforamción de la persona, cuya llave es también el número de contrato
func getInfoPersona(proveedores map[string]models.InformacionProveedor) (map[string]models.InformacionPersonaNatural, error) {
	personas := make(map[string]models.InformacionPersonaNatural)
	var persona models.InformacionPersonaNatural
	for key, value := range proveedores {
		if err := getJson("http://"+beego.AppConfig.String("agoraServicio")+"/informacion_persona_natural/"+value.NumDocumento, &persona); err == nil {
			personas[key] = persona
		} else {
			return nil, err
		}
	}
	return personas, nil
}

/* getPagosSeg busca todos los pagos correspondientes a seguridad social en los
conceptos de titan y devuelve un mapa cuya llave es el nombre del pago y el valor es el id del pago */
func getPagosSeg() (map[int]string, error) {
	var f interface{}
	pagos := make(map[int]string)

	if err := getJson("http://"+beego.AppConfig.String("titanServicio")+"/concepto_nomina"+
		"?limit=0&"+
		"fields=Id,NombreConcepto&"+
		"query=NaturalezaConcepto.Nombre:seguridad_social", &f); err != nil {
		return nil, err
	}
	interfaceArr := f.([]interface{})
	for i := range interfaceArr {
		pagos[int(interfaceArr[i].(map[string]interface{})["Id"].(float64))] = interfaceArr[i].(map[string]interface{})["NombreConcepto"].(string)
	}
	return pagos, nil
}

/*
// Get ...
// @Title ObtenerSegSocial_X_Facultad
// @Param idPeriodoPago "corresponde al al id del periodo del pago de seguridad social"
// @Succes 201 {object} models.Pago
// @Failure 403 :id is empty
// @router /:id [get]
func (c *PagoController) ObtenerSegSocial_X_Facultad() {
	idPeriodoPago := c.Ctx.Input.Param(":id")
	id, _ : := strconv.Atoi(idStr)
	v, err := models
}

*/
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
