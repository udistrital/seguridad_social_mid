package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/seguridad_social_mid/golog"
	"github.com/udistrital/seguridad_social_mid/models"
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
	c.Mapping("GetInfoCabecera", c.GetInfoCabecera)
	c.Mapping("ConceptosIbcParafiscales", c.ConceptosIbcParafiscales)
}

// GetInfoCabecera ...
// @Title GetInfoCabecera
// @Description Obtiene información adicional para la cabecera
// con el total de pagos de salud y pensión del empleado
// @Param	idPeriodoPago		id del periodo pago de seguridad social
// @router /GetInfoCabecera/:idPreliquidacion/:tipoPlanilla [get]
func (c *PagoController) GetInfoCabecera() {
	idStr := c.Ctx.Input.Param(":idPreliquidacion")
	tipoPlanilla := c.Ctx.Input.Param(":tipoPlanilla")
	var detallesPreliquidacion []models.DetallePreliquidacion

	err := getJson("http://"+beego.AppConfig.String("titanServicio")+
		"/detalle_preliquidacion"+
		"?limit=-1"+
		"&query=Concepto.NombreConcepto:ibc_liquidado"+
		",Preliquidacion.Id:"+idStr, &detallesPreliquidacion)

	if err != nil {
		ImprimirError("error en GetTotalIbc()", err)
	}

	acumuladorPreliquidacion := 0.0
	for _, value := range detallesPreliquidacion {
		acumuladorPreliquidacion += value.ValorCalculado
	}
	anioActual, mesActual, _ := time.Now().Date()

	mesString := strconv.Itoa(int(mesActual))
	anioString := strconv.Itoa(anioActual)

	totalPreliquidacion := AproximarPesoSuperior(acumuladorPreliquidacion, 100)

	if int(mesActual) < 10 {
		mesString = "0" + mesString
	}

	cabecera := models.CabeceraPlanilla{
		Codigo:           models.Columna{Valor: "0100000", Longitud: 7},
		NombreProveedor:  models.Columna{Valor: "UNIVERSIDAD DISTRITAL FRANCISCO JOSÉ DE CALDAS", Longitud: 200},
		NitProveedor:     models.Columna{Valor: "NI899999230", Longitud: 18},
		CodigoArl:        models.Columna{Valor: "14-23", Longitud: 6},
		PeriodoPension:   models.Columna{Valor: anioString + "-" + mesString, Longitud: 27},
		CantidadPersonas: models.Columna{Valor: len(detallesPreliquidacion), Longitud: 5},
		TotalNomina:      models.Columna{Valor: totalPreliquidacion, Longitud: 12},
		CodigoProveedor:  models.Columna{Valor: 1, Longitud: 2},
		CodigoOperador:   models.Columna{Valor: 83, Longitud: 2},
	}

	switch tipoPlanilla {
	case "CT":
		cabecera.TipoPlanilla = models.Columna{Valor: "7Y", Longitud: 22}
		cabecera.Sucursal = models.Columna{Valor: "S59", Longitud: 51}
		cabecera.PeriodoSalud = models.Columna{Valor: anioString + "-" + mesString, Longitud: 7}
	default:
		if int(mesActual) == 12 {
			mesString = "01"
		} else {
			if int(mesActual) < 10 {
				mesString = "0" + strconv.Itoa(int(mesActual-1))
			}
		}
		cabecera.TipoPlanilla = models.Columna{Valor: "7E", Longitud: 22}
		cabecera.Sucursal = models.Columna{Valor: "S01", Longitud: 51}
		cabecera.PeriodoSalud = models.Columna{Valor: anioString + "-" + mesString, Longitud: 7}
	}

	c.Data["json"] = cabecera
	c.ServeJSON()
}

// SumarPagosSalud ...
// @Title Sumar pagos de salid
// @Description Suma el total de los pagos de salud y pensión de ud
// con el total de pagos de salud y pensión del empleado
// @Param	idPeriodoPago		id del periodo pago de seguridad social
// @router /SumarPagosSalud/:idPeriodoPago [get]
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
// @router /ConceptosIbc [get]
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
		fmt.Println(nombres, len(nombres))
		for i := 0; i < len(nombres); i++ {
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

// ConceptosIbcParafiscales ...
// @Title ConceptosIbcParafiscales
// @Description Obtiene todos los conceptos IBC Parafiscales del ruler y los cruza con los conceptos de nómina
// @router /conceptos_ibc_parafiscales [get]
func (c *PagoController) ConceptosIbcParafiscales() {
	var predicados []models.Predicado
	var conceptos []models.Concepto
	var conceptosIbc []models.ConceptosIbc
	err := getJson("http://"+beego.AppConfig.String("rulerServicio")+
		"/predicado?limit=-1&query=Nombre__startswith:concepto_ibc_parafiscales,Dominio.Id:19", &predicados)

	errConceptoTitan := getJson("http://"+beego.AppConfig.String("titanServicio")+
		"/concepto_nomina?limit=-1", &conceptos)

	if err != nil && errConceptoTitan != nil {
		c.Data["json"] = err.Error() + errConceptoTitan.Error()
	} else {
		fmt.Println(predicados)
		nombres := golog.GetString(FormatoReglas(predicados), "concepto_ibc_parafiscales(X,Y).", "X")
		estados := golog.GetString(FormatoReglas(predicados), "concepto_ibc_parafiscales(X,Y).", "Y")
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

// NovedadesPorPersona ...
// @Title NovedadesPorPersona
// @Description Obtiene todos los conceptos IBC del ruler y los cruza con los conceptos de nómina
// @Success 200 {object} []models.NovedadesPersonaSS
// @router /NovedadesPorPersona/:persona [get]
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

// CalcularSegSocial ...
// @Title CalcularSegSocial
// @Description Cálcula la seguridad social para una preliquidación correspondiente
// @Param	id		id de la preliquidación
// @Success 200 {object} []*models.PagoSeguridadSocial
// @router /CalcularSegSocial/:id [get]
func (c *PagoController) CalcularSegSocial() {
	idStr := c.Ctx.Input.Param(":id")
	_, err := strconv.Atoi(idStr)

	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		var (
			alertas               []string
			detallePreliquidacion []models.DetallePreliquidacion
		)

		err := getJson("http://"+beego.AppConfig.String("titanServicio")+"/detalle_preliquidacion"+
			"?limit=-1&query=Preliquidacion.Id:"+idStr+",Concepto.NombreConcepto:ibc_liquidado", &detallePreliquidacion)
		if err != nil {
			ImprimirError("error en CalcularSegSocial()", err)
			alertas = append(alertas, "error al traer detalle liquidacion")
			c.Data["json"] = alertas
		} else {
			var (
				predicado            []models.Predicado
				pagosSeguridadSocial []models.PagoSeguridadSocial
			)

			idDetallePreliquidacion := detallePreliquidacion[0].Preliquidacion.Id
			idNomina := strconv.Itoa(detallePreliquidacion[0].Preliquidacion.Nomina.Id)

			totalSaludLiquidacion := valorConceptoTotal(idStr, "salud")
			totalPensionLiquidacion := valorConceptoTotal(idStr, "pension")
			totalFondoSolidaridad := valorFondoTotal(idNomina)

			for i := range detallePreliquidacion {
				persona := strconv.FormatInt(detallePreliquidacion[i].Persona, 10)
				valorCalculado := strconv.Itoa(int(detallePreliquidacion[i].ValorCalculado))

				predicado = append(predicado, []models.Predicado{
					models.Predicado{Nombre: "ibc(" + persona + "," + valorCalculado + ", salud)."},
					models.Predicado{Nombre: "ibc(" + persona + "," + valorCalculado + ", riesgos)."},
					models.Predicado{Nombre: "ibc(" + persona + "," + valorCalculado + ", apf)."},
					models.Predicado{Nombre: "v_salud_func(" + persona + ", " + strconv.Itoa(totalSaludLiquidacion[persona]) + ")."},
					models.Predicado{Nombre: "v_pen_func(" + persona + ", " + strconv.Itoa(totalPensionLiquidacion[persona]) + ")."},
				}...)
			}

			reglas := CargarReglasBase() + FormatoReglas(predicado) + cargarNovedades()

			idProveedores := golog.GetInt64(reglas, "v_salud_ud(I,Y).", "I")
			saludUd := golog.GetFloat(reglas, "v_salud_ud(I,Y).", "Y")
			saludTotal := golog.GetInt64(reglas, "v_total_salud(X,T).", "T")
			pensionUd := golog.GetFloat(reglas, "v_pen_ud(I,Y).", "Y")
			pensionTotal := golog.GetInt64(reglas, "v_total_pen(X,T).", "T")
			arl := golog.GetInt64(reglas, "v_arl(I,Y).", "Y")
			caja := golog.GetInt64(reglas, "v_caja(I,Y).", "Y")
			icbf := golog.GetInt64(reglas, "v_icbf(I,Y).", "Y")

			for i := range idProveedores {
				aux := models.PagoSeguridadSocial{
					NombrePersona:           "",
					IdProveedor:             idProveedores[i],
					SaludUd:                 saludUd[i],
					SaludTotal:              saludTotal[i],
					PensionUd:               pensionUd[i],
					PensionTotal:            pensionTotal[i],
					FondoSolidaridad:        totalFondoSolidaridad[int(idProveedores[i])],
					IdPreliquidacion:        idDetallePreliquidacion,
					IdDetallePreliquidacion: detallePreliquidacion[i].Id,
					Arl:                     arl[i],
					Caja:                    caja[i],
					Icbf:                    icbf[i]}

				pagosSeguridadSocial = append(pagosSeguridadSocial, aux)
			}

			c.Data["json"] = pagosSeguridadSocial
		}
		c.ServeJSON()
	}
}

func valorFondoTotal(idNomina string) (valores map[int]float64) {
	var conceptoNominaPorPersona []models.ConceptoNominaPorPersona

	err := getJson("http://"+beego.AppConfig.String("titanServicio")+"/concepto_nomina_por_persona"+
		"?limit=0&query=Concepto.NombreConcepto:fondoSolidaridad,Activo:true,Nomina.Id:"+idNomina,
		&conceptoNominaPorPersona)

	if err != nil {
		ImprimirError("error en valorPagoFondoSolidaridad() ", err)
	}

	if len(conceptoNominaPorPersona) > 0 {
		valores = make(map[int]float64)
		for _, value := range conceptoNominaPorPersona {
			valores[value.Persona] = value.ValorNovedad
		}
	}

	return
}

// valorConceptoTotal obtiene todos los datos del concepto de una preliquidación y devuelve un mapa
// el mapa valores tiene la forma [idPersona] = totalValorCalculado
func valorConceptoTotal(idLiquidacion, concepto string) (valores map[string]int) {
	var detalleLiquSalud []models.DetallePreliquidacion

	err := getJson("http://"+beego.AppConfig.String("titanServicio")+"/detalle_preliquidacion"+
		"?limit=0&query=Preliquidacion:"+idLiquidacion+",Concepto.NombreConcepto:"+concepto, &detalleLiquSalud)

	if err != nil {
		ImprimirError("error en valorSaludTotal()", err)
	} else {
		valores = make(map[string]int)

		for _, value := range detalleLiquSalud {
			valores[strconv.FormatInt(value.Persona, 10)] += int(value.ValorCalculado)
		}
	}

	return
}

// valorSaludEmpleado ...
// @Title Valor Salud Empleado
// @Description Crea todos los hechos con la información del valor de salud
// @Param	idLiquidacion		id de la liquidacion correspondiente
// @Param	persona				id correspondiente a la columna persona
func valorSaludEmpleado(idLiquidacion, persona string) (predicado models.Predicado) {
	var detalleLiquSalud []models.DetallePreliquidacion
	var totalSalud int
	errSalud := getJson("http://"+beego.AppConfig.String("titanServicio")+"/detalle_preliquidacion"+
		"?limit=0&query=Preliquidacion:"+idLiquidacion+",Concepto.NombreConcepto:salud,Persona:"+persona, &detalleLiquSalud)

	if errSalud != nil {
		// logs.Error("Error en ValorSaludEmpleado:", errSalud)
	} else {
		for index := 0; index < len(detalleLiquSalud); index++ {
			totalSalud += int(detalleLiquSalud[index].ValorCalculado)
		}
		predicado = models.Predicado{Nombre: "v_salud_func(" + persona + ", " + strconv.Itoa(totalSalud) + ")."}
	}

	return
}

// ValorPensionEmpleado ...
// @Title Valor Pensión Empleado
// @Description Crea todos los hechos con la información del valor de la pensión
// @Param	idLiquidacion		id de la liquidacion correspondiente
func ValorPensionEmpleado(idLiquidacion, persona string) (predicado models.Predicado) {
	var detalleLiquPension []models.DetallePreliquidacion
	var totalPension int

	errPension := getJson("http://"+beego.AppConfig.String("titanServicio")+"/detalle_preliquidacion"+
		"?limit=0&query=Preliquidacion:"+idLiquidacion+",Concepto.NombreConcepto:pension,Persona:"+persona, &detalleLiquPension)

	if errPension != nil {
		logs.Error("Error en ValorPensionEmpleado:", errPension)
		fmt.Println("http://" + beego.AppConfig.String("titanServicio") + "/detalle_preliquidacion" +
			"?limit=0&query=Preliquidacion:" + idLiquidacion + ",Concepto.NombreConcepto:pension,Persona:" + persona)
	} else {
		fmt.Println(detalleLiquPension)
		for index := 0; index < len(detalleLiquPension); index++ {
			// fmt.Println(detalleLiquPension[index])
			totalPension += int(detalleLiquPension[index].ValorCalculado)
		}
		predicado = models.Predicado{Nombre: "v_pen_func(" + persona + ", " + strconv.Itoa(totalPension) + ")."}
	}
	return
}

// SaludHCHonorarios ...
//@Title Valor correspondiente a salud de hora catedra honorarios
//@Description consulta una preliqudacion correspondiente a hora catedra
//@Param idLiquidacion id de la preliquidacion correspondiente
func SaludHCHonorarios(idLiquidacion string) (valorSaludEmpleado string) {
	var detalleLiquSalud []models.DetallePreliquidacion
	var predicado []models.Predicado

	errSalud := getJson("http://"+beego.AppConfig.String("titanServicio")+"/detalle_preliquidacion"+
		"?limit=0&query=Preliquidacion:"+idLiquidacion+",Concepto.NombreConcepto:salud", &detalleLiquSalud)

	if errSalud != nil {
		logs.Error("Error en SaludHCHonorarios:", errSalud)
	} else {
		for index := 0; index < len(detalleLiquSalud); index++ {
			predicado = append(predicado, models.Predicado{Nombre: "v_salud_func(" + detalleLiquSalud[index].NumeroContrato + ", " + strconv.Itoa(int(detalleLiquSalud[index].ValorCalculado)) + ")."})
			valorSaludEmpleado += predicado[index].Nombre + "\n"
		}
	}
	return
}

// cargarNovedades ...
// @Title cargarNovedades
// @Description obtiene todas los conceptos con naturaleza seguridad_social desde detalle_preliquidacion por el id
func cargarNovedades() (novedades string) {
	var conceptoNominaPorPersona []models.ConceptoNominaPorPersona
	var predicado models.Predicado

	errLincNo := getJson("http://"+beego.AppConfig.String("titanServicio")+"/concepto_nomina_por_persona?"+
		"limit=0&query=Concepto.EstadoConceptoNomina.Nombre:activo", &conceptoNominaPorPersona)
	if errLincNo != nil {
		fmt.Println("error en cargarNovedades()", errLincNo)
		//beego.Error("error en cargarNovedades()", errLincNo)
	} else {
		for index := 0; index < len(conceptoNominaPorPersona); index++ {
			predicado = models.Predicado{Nombre: "novedad_persona(" + conceptoNominaPorPersona[index].Concepto.NombreConcepto + ", " + strconv.Itoa(conceptoNominaPorPersona[index].Persona) + ")."}
			novedades += predicado.Nombre + "\n"
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
				c.Ctx.Output.SetStatus(500)
				c.ServeJSON()
				return
			}
		}

		mapProveedores, err := GetInfoProveedores(PeriodoPago.Personas)
		if err != nil {
			c.Data["json"] = err.Error()
			c.Ctx.Output.SetStatus(500)
			c.ServeJSON()
			return
		}

		mapPersonas, err := InfoPersona(mapProveedores)
		if err != nil {
			fmt.Println("aquí fue el error...")
			c.Data["json"] = err.Error()
			c.Ctx.Output.SetStatus(500)
			c.ServeJSON()
			return
		}

		pagosSeg, err := GetPagosSeguridadSocial()
		if err != nil {
			c.Data["json"] = err.Error()
			c.Ctx.Output.SetStatus(500)
			c.ServeJSON()
			return
		}
		contPagos, contContratista := 0, 0 // conPagos sirve para que cuente los 5 pagos de seguridad social, contContratista es para que recorrer los contratistas
		for i := range PeriodoPago.Pagos {
			nombrePago := pagosSeg[PeriodoPago.Pagos[i].TipoPago]
			aux := PeriodoPago.Personas[contContratista]
			switch nombrePago {
			case "arl":
				PeriodoPago.Pagos[i].EntidadPago = mapPersonas[aux].IdArl
			case "pension_ud":
				PeriodoPago.Pagos[i].EntidadPago = mapPersonas[aux].IdFondoPension
			case "salud_ud":
				PeriodoPago.Pagos[i].EntidadPago = mapPersonas[aux].IdEps
			case "caja_compensacion":
				PeriodoPago.Pagos[i].EntidadPago = mapPersonas[aux].IdCajaCompensacion
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
		} else {
			c.Data["json"] = err.Error()
			c.Ctx.Output.SetStatus(500)
			c.ServeJSON()
			return
		}
		c.Data["json"] = alerta

	} else {
		c.Data["json"] = err.Error()
		c.Abort("500")
	}
	c.ServeJSON()
}

// GetInfoProveedores Recibe un arreglo de strings con las personas de una nómina y devuelve un map con la información del proveedor
func GetInfoProveedores(personasNomina []string) (map[string]models.InformacionProveedor, error) {
	var proveedor models.InformacionProveedor
	personas := make(map[string]models.InformacionProveedor)

	for i := range personasNomina {
		fmt.Println("http://" + beego.AppConfig.String("agoraServicio") + "/informacion_proveedor/" + personasNomina[i])
		if err := getJson("http://"+beego.AppConfig.String("agoraServicio")+"/informacion_proveedor/"+personasNomina[i], &proveedor); err == nil {
			personas[personasNomina[i]] = proveedor
		} else {
			return nil, err
		}
	}
	return personas, nil
}

// InfoPersona recibe un map de proveedores para consultar el número de contrato y devuelve un mapa con la inforamción de la persona, cuya llave es también el número de contrato
func InfoPersona(proveedores map[string]models.InformacionProveedor) (map[string]models.InformacionPersonaNatural, error) {
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

// GetInfoProveedor devuelve un map con la información de todos los proveedores de tipo persona NATURAL
func GetInfoProveedor() (map[int64]models.InformacionProveedor, error) {
	var infoProveedores []models.InformacionProveedor
	proveedores := make(map[int64]models.InformacionProveedor)

	if err := getJson("http://"+beego.AppConfig.String("agoraServicio")+"informacion_proveedor?"+
		"limit=-1&query=Tipopersona:NATURAL", &infoProveedores); err == nil {
		for _, proveedor := range infoProveedores {
			proveedores[int64(proveedor.Id)] = proveedor
		}
	} else {
		ImprimirError("Error en GetInfoProveedor: ", err)
		return nil, err
	}

	return proveedores, nil
}

/*
GetInfoPersonas Recibe un arreglo de strings con los contratos, cruza cada uno de los elementos del arreglo con un valor de proveedores y retonar un map
que tenga la información del proveedor y cuya llave sea el id del proveedor
*/
func GetInfoPersonas(detallesPreliquidacion []models.DetallePreliquidacion) (map[string]models.InformacionPersonaNatural, error) {
	var infoProveedores []models.InformacionProveedor
	var personasNaturales []models.InformacionPersonaNatural

	proveedores := make(map[string]models.InformacionProveedor)
	personas := make(map[string]models.InformacionPersonaNatural)

	if err := getJson("http://"+beego.AppConfig.String("agoraServicio")+"/informacion_proveedor?limit=-1", &infoProveedores); err != nil {
		ImprimirError("Error en GetInfoPersonas: ", err)
		return nil, err
	}

	if err := getJson("http://"+beego.AppConfig.String("agoraServicio")+"/informacion_persona_natural?limit=-1", &personasNaturales); err != nil {
		ImprimirError("Error en GetInfoPersonas: ", err)
		return nil, err
	}

	for _, detallePreliquidacion := range detallesPreliquidacion {
		for _, infoProveedor := range infoProveedores {
			if detallePreliquidacion.Persona == int64(infoProveedor.Id) {
				proveedores[strconv.FormatInt(detallePreliquidacion.Persona, 10)] = infoProveedor
				break
			}
		}

	}

	for key, proveedor := range proveedores {
		for _, personaNatural := range personasNaturales {
			if proveedor.NumDocumento == personaNatural.Id {
				personas[key] = personaNatural
			}
		}
	}

	// for i := range idProveedores {
	// }
	return personas, nil
}

/*
GetInfoPersona recibe un map de proveedores para consultar el número de contrato y
devuelve un map con la inforamción de la persona, cuya llave es también el número de contrato
*/
func GetInfoPersona(proveedores map[int64]models.InformacionProveedor) (map[string]models.InformacionPersonaNatural, error) {
	personas := make(map[string]models.InformacionPersonaNatural)
	var persona models.InformacionPersonaNatural
	for key, value := range proveedores {
		fmt.Println("http://" + beego.AppConfig.String("agoraServicio") + "/informacion_persona_natural/" + value.NumDocumento)
		if err := getJson("http://"+beego.AppConfig.String("agoraServicio")+"/informacion_persona_natural/"+value.NumDocumento, &persona); err == nil {
			fmt.Println(key)
			personas[strconv.FormatInt(key, 10)] = persona
		} else {
			return nil, err
		}
	}
	return personas, nil
}

// ComprovarCajaProveedor revisa si el proveedor tiene una caja de compensación asociada
func ComprovarCajaProveedor(cedulaProveedor string) (tieneCaja bool, err error) {
	var persona models.InformacionPersonaNatural
	err = getJson("http://"+beego.AppConfig.String("agoraServicio")+"/informacion_persona_natural/"+cedulaProveedor, &persona)
	if err != nil {
		return
	}

	if persona.IdCajaCompensacion != 0 {
		tieneCaja = true
	}
	return
}

/*
GetPagosSeguridadSocial busca todos los pagos correspondientes a seguridad social en los
conceptos de titan y devuelve un mapa cuya llave es el nombre del pago y el valor es el id del pago
*/
func GetPagosSeguridadSocial() (map[int]string, error) {
	var f interface{}
	pagos := make(map[int]string)

	if err := getJson("http://"+beego.AppConfig.String("titanServicio")+"/concepto_nomina"+
		"?limit=0&"+
		"fields=Id,NombreConcepto&"+
		"query=NaturalezaConcepto.Nombre:seguridad_social"+
		",TipoConcepto.Nombre:pago_seguridad_social", &f); err != nil {
		return nil, err
	}

	interfaceArr := f.([]interface{})
	for i := range interfaceArr {
		pagos[int(interfaceArr[i].(map[string]interface{})["Id"].(float64))] = interfaceArr[i].(map[string]interface{})["NombreConcepto"].(string)
	}
	return pagos, nil
}
