package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"

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

// ConceptosIbc ...
// @Title ConceptosIbc
// @Description Obtiene todos los conceptos IBC del ruler y los cruza con los conceptos de nómina
// @router ConceptosIbc/ [get]
func (c *PagoController) ConceptosIbc() {
	var predicados []models.Predicado
	var conceptos []models.Concepto
	var conceptosIbc []models.ConceptosIbc
	err := getJson("http://"+beego.AppConfig.String("rulerServicio")+
		"/predicado?limit=-1&query=Nombre__startswith:concepto_ibc,Dominio.Id:19", &predicados)

	errConceptoTitan := getJson("http://"+beego.AppConfig.String("titanServicio")+
		"/concepto_nomina?limit=-1", &conceptos)

	if err != nil && errConceptoTitan != nil {
		c.Data["json"] = err.Error() + errConceptoTitan.Error()
	} else {
		nombres := golog.GetString(FormatoReglas(predicados), "concepto_ibc(X,Y).", "X")
		estados := golog.GetString(FormatoReglas(predicados), "concepto_ibc(X,Y).", "Y")
		for i := 0; i < len(predicados); i++ {
			for j := 0; j < len(conceptos); j++ {
				if nombres[i] == conceptos[j].NombreConcepto {
					aux := models.ConceptosIbc{
						Id:               predicados[i].Id,
						Nombre:           nombres[i],
						Descripcion:      conceptos[j].AliasConcepto,
						DescripcionHecho: predicados[i].Descripcion,
						Estado:           true,
						Dominio:          predicados[i].Dominio,
						TipoPredicado:    predicados[i].TipoPredicado,
					}
					if estados[i] == "inactivo" {
						aux.Estado = false
					}
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

	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		var (
			alertas               []string
			proveedores           []string
			predicado             []models.Predicado
			detallePreliquidacion []models.DetallePreliquidacion
			pagosSeguridadSocial  []*models.PagoSeguridadSocial
		)

		err := getJson("http://"+beego.AppConfig.String("titanServicio")+"/detalle_preliquidacion"+
			"?limit=-1&query=Preliquidacion.Id:"+idStr+",Concepto.NombreConcepto:ibc_liquidado", &detallePreliquidacion)

		if err != nil {
			beego.Error(err)
			alertas = append(alertas, "error al traer detalle liquidacion")
			c.Data["json"] = alertas
		} else {
			idDetallePreliquidacion := detallePreliquidacion[0].Preliquidacion.Id

			for index := 0; index < len(detallePreliquidacion); index++ {
				persona := strconv.Itoa(detallePreliquidacion[index].Persona)
				valorCalculado := strconv.Itoa(int(detallePreliquidacion[index].ValorCalculado))
				predicado = append(predicado, models.Predicado{Nombre: "ibc(" + persona + "," + valorCalculado + ", salud)."})
				predicado = append(predicado, models.Predicado{Nombre: "ibc(" + persona + "," + valorCalculado + ", riesgos)."})
				predicado = append(predicado, models.Predicado{Nombre: "ibc(" + persona + "," + valorCalculado + ", apf)."})
				predicado = append(predicado, valorSaludEmpleado(idStr, persona), ValorPensionEmpleado(idStr, persona))
			}

			reglas := CargarReglasBase() + FormatoReglas(predicado) + cargarNovedades(idStr)

			idProveedores := golog.GetInt64(reglas, "v_salud_ud(I,Y).", "I")
			saludUd := golog.GetFloat(reglas, "v_salud_ud(I,Y).", "Y")
			saludTotal := golog.GetInt64(reglas, "v_total_salud(X,T).", "T")
			pensionUd := golog.GetFloat(reglas, "v_pen_ud(I,Y).", "Y")
			pensionTotal := golog.GetInt64(reglas, "v_total_pen(X,T).", "T")
			arl := golog.GetInt64(reglas, "v_arl(I,Y).", "Y")
			caja := golog.GetInt64(reglas, "v_caja(I,Y).", "Y")
			icbf := golog.GetInt64(reglas, "v_icbf(I,Y).", "Y")

			// Acá se debe cambiar como se arma el modelo
			for index := 0; index < len(idProveedores); index++ {
				proveedores = append(proveedores, fmt.Sprint(idProveedores[index]))
				aux := &models.PagoSeguridadSocial{
					NombrePersona:           "",
					IdProveedor:             idProveedores[index],
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

			mapProveedores, _ := GetInfoProveedor(proveedores)

			for i := range pagosSeguridadSocial {
				pagosSeguridadSocial[i].NombrePersona = mapProveedores[fmt.Sprint(pagosSeguridadSocial[i].IdProveedor)].NomProveedor
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
// @Param	persona				id correspondiente a la columna persona
func valorSaludEmpleado(idLiquidacion, persona string) (predicado models.Predicado) {
	var detalleLiquSalud []models.DetallePreliquidacion

	errSalud := getJson("http://"+beego.AppConfig.String("titanServicio")+"/detalle_preliquidacion"+
		"?limit=0&query=Preliquidacion:"+idLiquidacion+",Concepto.NombreConcepto:salud,Persona:"+persona, &detalleLiquSalud)

	if errSalud == nil {
		for index := 0; index < len(detalleLiquSalud); index++ {
			predicado = models.Predicado{Nombre: "v_salud_func(" + strconv.Itoa(detalleLiquSalud[index].Persona) + ", " + strconv.Itoa(int(detalleLiquSalud[index].ValorCalculado)) + ")."}
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
		beego.Error("Error en ValorSaludEmpleado:", errSalud)
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
func ValorPensionEmpleado(idLiquidacion, persona string) (predicado models.Predicado) {
	var detalleLiquPension []models.DetallePreliquidacion

	errPension := getJson("http://"+beego.AppConfig.String("titanServicio")+"/detalle_preliquidacion"+
		"?limit=0&query=Preliquidacion:"+idLiquidacion+",Concepto.NombreConcepto:pension,Persona:"+persona, &detalleLiquPension)

	if errPension != nil {
		beego.Error("Error en ValorPensionEmpleado:", errPension)
	} else {
		for index := 0; index < len(detalleLiquPension); index++ {
			predicado = models.Predicado{Nombre: "v_pen_func(" + strconv.Itoa(detalleLiquPension[index].Persona) + ", " + strconv.Itoa(int(detalleLiquPension[index].ValorCalculado)) + ")."}
		}
	}
	return
}

// cargarNovedades ...
// @Title cargarNovedades
// @Description obtiene todas los conceptos con naturaleza seguridad_social desde detalle_preliquidacion por el id
// @Param	id		path 	string	true		"The key for staticblock"
func cargarNovedades(id string) (novedades string) {
	var conceptosPreliquidacion []models.DetallePreliquidacion
	var predicado []models.Predicado

	errLincNo := getJson("http://"+beego.AppConfig.String("titanServicio")+"/detalle_preliquidacion"+
		"?limit=0&query=Preliquidacion:"+id+",Concepto.NaturalezaConcepto.Nombre:seguridad_social&fields=Concepto,NumeroContrato", &conceptosPreliquidacion)

	if errLincNo != nil {
		fmt.Println("error en cargarNovedades()", errLincNo)
	} else {
		for index := 0; index < len(conceptosPreliquidacion); index++ {
			predicado = append(predicado, models.Predicado{Nombre: "novedad_persona(" + conceptosPreliquidacion[index].Concepto.NombreConcepto + ", " + strconv.Itoa(conceptosPreliquidacion[index].Persona) + ")."})
			novedades += predicado[index].Nombre + "\n"
		}
	}
	return
}

// Devuelve true si hay un periodo ya registrado, de lo contrario devuelve false
func validarPeriodo(PeriodoPago models.TrPeriodoPago) (bool, *models.PeriodoPago) {
	var (
		periodoAntiguo    []*models.PeriodoPago
		periodoModificado *models.PeriodoPago
		estadoSegSocial   []*models.EstadoSeguridadSocial
	)

	if err := getJson("http://"+beego.AppConfig.String("segSocialService")+"/periodo_pago?"+
		"query=Mes:"+strconv.Itoa(int(PeriodoPago.PeriodoPago.Mes))+
		",Anio:"+strconv.Itoa(int(PeriodoPago.PeriodoPago.Anio))+
		",TipoLiquidacion:"+PeriodoPago.PeriodoPago.TipoLiquidacion+
		"&EstadoSeguridadSocial.Nombre:Abierta", &periodoAntiguo); err == nil {

		if err := getJson("http://"+beego.AppConfig.String("segSocialService")+"/estado_seguridad_social?query=Nombre:reemplazada", &estadoSegSocial); err != nil {
			return false, nil
		}

		if len(periodoAntiguo) < 1 || periodoAntiguo[0].Id == 0 {
			return false, nil
		}

		periodoModificado = periodoAntiguo[0]
		periodoModificado.EstadoSeguridadSocial = estadoSegSocial[0]

	} else {
		return false, nil
	}

	return true, periodoModificado
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

		if validar, periodoModificado := validarPeriodo(PeriodoPago); validar {
			if err = sendJson("http://"+beego.AppConfig.String("segSocialService")+"/periodo_pago/"+strconv.Itoa(periodoModificado.Id), "PUT", &alerta, periodoModificado); err != nil {
				c.Data["json"] = err.Error()
			}
		}

		mapProveedores, err := GetInfoProveedor(PeriodoPago.Contratos)
		if err != nil {
			c.Data["json"] = err.Error()
		}

		mapPersonas, err := GetInfoPersona(mapProveedores)
		if err != nil {
			c.Data["json"] = err.Error()
		}

		pagosSeg, _ := getPagosSeg()
		contPagos, contContratista := 0, 0 // conPagos sirve para que cuente los 5 pagos de seguridad social, contContratista es para que recorrer los contratistas
		fmt.Println(PeriodoPago)
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
				fmt.Println("Aqui es el error......")
				fmt.Println(PeriodoPago.Pagos[i])
				PeriodoPago.Pagos[i].EntidadPago = mapPersonas[PeriodoPago.Contratos[contContratista]].IdCajaCompensacion
			default: // ICBF
				PeriodoPago.Pagos[i].EntidadPago = 0
			}
			contPagos++
			// Como se tienen 5 pagos, cada vez que se asigne uno nuevo, el contratista debe pasar al siguiente
			if contPagos == 5 {
				contPagos = 0
				contContratista++
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

// GetInfoProveedor Recibe un arreglo de strings con los contratos y devuelve un map con la información del proveedor
func GetInfoProveedor(idProveedores []string) (map[string]models.InformacionProveedor, error) {
	proveedores := make(map[string]models.InformacionProveedor)
	var proveedor models.InformacionProveedor

	for i := range idProveedores {
		if err := getJson("http://"+beego.AppConfig.String("agoraServicio")+"/informacion_proveedor/"+idProveedores[i], &proveedor); err == nil {
			proveedores[idProveedores[i]] = proveedor
		} else {
			beego.Error("error en GetInfoProveedor: ", err.Error())
			return nil, err
		}
	}
	return proveedores, nil
}

// GetInfoPersona recibe un map de proveedores para consultar el número de contrato y devuelve un mapa con la inforamción de la persona, cuya llave es también el número de contrato
func GetInfoPersona(proveedores map[string]models.InformacionProveedor) (map[string]models.InformacionPersonaNatural, error) {
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
