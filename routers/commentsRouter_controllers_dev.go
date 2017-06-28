package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:PagoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:PagoController"],
		beego.ControllerComments{
			Method:           "CalcularSegSocial",
			Router:           `CalcularSegSocial/:id`,
			AllowHTTPMethods: []string{"get"},
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:PagoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:PagoController"],
		beego.ControllerComments{
			Method:           "NovedadesPorPersona",
			Router:           `NovedadesPorPersona/:persona`,
			AllowHTTPMethods: []string{"get"},
			Params:           nil})
	/*
		beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:FuncionarioProveedorController"] = append(beego.GlobalControllerRouter["github.com/miguelramirez93/titan_api_crud/controllers:FuncionarioProveedorController"],
			beego.ControllerComments{
				Method:           "ConsultarIDProveedor",
				Router:           `/`,
				AllowHTTPMethods: []string{"post"},
				Params:           nil})*/

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:PagoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:PagoController"],
		beego.ControllerComments{
			Method:           "ConceptosIbc",
			Router:           `ConceptosIbc`,
			AllowHTTPMethods: []string{"get"},
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:PlanillasController"] = append(beego.GlobalControllerRouter["github.com/udistrital/ss_mid_api/controllers:PlanillasController"],
		beego.ControllerComments{
			Method:           "GenerarPlanillaActivos",
			Router:           `GenerarPlanillaActivos/:id`,
			AllowHTTPMethods: []string{"get"},
			Params:           nil})
}
