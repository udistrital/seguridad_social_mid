package controllers

import (
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/manucorporat/try"
	"github.com/udistrital/ss_mid_api/models"
)

// PagoHonorariosController operations for Pago_honorarios
type PagoHonorariosController struct {
	beego.Controller
}

// URLMapping ...
func (c *PagoHonorariosController) URLMapping() {
}

// CalcularSegSocialHonorarios ...
// @Title CalcularSegSocialHonorarios
// @Description Cálcula la seguridad social para una preliquidación de honorarios o contratistas correspondiente
// @Param	id		id de la preliquidación
// @Success 200 {object} []*models.PagoSeguridadSocial
// @router CalcularSegSocialHonorarios/:idPreliquidacion [get]
func (c *PagoController) CalcularSegSocialHonorarios() {
	var alerta models.Alert

	idPreliquidacion := c.Ctx.Input.Param(":idPreliquidacion")
	preliquidacion, err := strconv.Atoi(idPreliquidacion)

	try.This(func() {
		if err != nil {
			panic(err)
		}
		var (
			detalleLiquSalud []models.DetallePreliquidacion
			proveedores      []string
		)

		pagosSeguridadSocial := make(map[int]*models.PagoSeguridadSocial)

		err = getJson("http://"+beego.AppConfig.String("titanServicio")+"/detalle_preliquidacion"+
			"?limit=-1&query=Preliquidacion.Id:"+idPreliquidacion+",Concepto.NombreConcepto:salud", &detalleLiquSalud)
		fmt.Println("http://" + beego.AppConfig.String("titanServicio") + "/detalle_preliquidacion" +
			"?limit=-1&query=Preliquidacion.Id:" + idPreliquidacion + ",Concepto.NombreConcepto:salud")
		if err != nil {
			panic(err)
		}

		for _, value := range detalleLiquSalud {
			proveedores = append(proveedores, strconv.Itoa(value.Persona))

			aux := &models.PagoSeguridadSocial{
				NombrePersona:           "",
				IdProveedor:             int64(value.Persona),
				SaludUd:                 0,
				SaludTotal:              int64(value.ValorCalculado),
				PensionUd:               0,
				PensionTotal:            0,
				Caja:                    0,
				Icbf:                    0,
				IdPreliquidacion:        preliquidacion,
				IdDetallePreliquidacion: value.Id,
				Arl:                     0}

			if pagosSeguridadSocial[value.Persona] == nil {
				pagosSeguridadSocial[value.Persona] = aux
			} else {
				fmt.Println("ya existe: ", value.Persona, "  ", pagosSeguridadSocial[value.Persona])
				if err = GetValoresSaludHonorarios(aux); err != nil {
					panic(err)
				}
			}

		}

		mapProveedores, _ := GetInfoProveedor(proveedores)
		for i := range pagosSeguridadSocial {
			pagosSeguridadSocial[i].NombrePersona = mapProveedores[fmt.Sprint(pagosSeguridadSocial[i].IdProveedor)].NomProveedor
		}

		c.Data["json"] = pagosSeguridadSocial
	}).Finally(func() {
		c.ServeJSON()
	}).Catch(func(e try.E) {
		alerta = models.Alert{
			Type: "error",
			Code: "404",
			Body: e,
		}
		c.Data["json"] = alerta
	})
}

// GetValoresSaludHonorarios asgina todos los valores de salud al apuntador pagoSeguridadSocial
//@Param pagoSeguridadSocial apuntador que tiene los pagos, la persona y la preliquidación
func GetValoresSaludHonorarios(pagoSeguridadSocial *models.PagoSeguridadSocial) (err error) {
	var detalleLiquidacion []models.DetallePreliquidacion

	err = getJson("http://"+beego.AppConfig.String("titanServicio")+"/detalle_preliquidacion"+
		"?limit=-1&query=Preliquidacion.Id:"+strconv.Itoa(pagoSeguridadSocial.IdPreliquidacion)+
		",Persona:"+fmt.Sprint(pagoSeguridadSocial.IdProveedor), &detalleLiquidacion)
	if err != nil {
		return
	}

	for _, value := range detalleLiquidacion {
		valorConcepto := int64(value.ValorCalculado)
		switch value.Concepto.NombreConcepto {
		case "pension":
			pagoSeguridadSocial.PensionTotal += valorConcepto
		case "arl":
			pagoSeguridadSocial.Arl += valorConcepto
		case "icbf":
			pagoSeguridadSocial.Icbf += valorConcepto
		case "caja_compensacion":
			pagoSeguridadSocial.Caja += valorConcepto
		}
	}
	return
}
