// Package routers paquete con la ruta de los controladores y servicios
// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact jcamilosarmientor@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/ss_mid_api/controllers"
)

func init() {

	// auditoria.InitMiddleware()

	ns := beego.NewNamespace("/v1",

		beego.NSNamespace("/periodo_pago",
			beego.NSInclude(
				&controllers.PeriodoPagoController{},
			),
		),

		beego.NSNamespace("/pago",
			beego.NSInclude(
				&controllers.PagoController{},
			),
		),

		beego.NSNamespace("/planillas",
			beego.NSInclude(
				&controllers.PlanillasController{},
			),
		),

		beego.NSNamespace("/aportante",
			beego.NSInclude(
				&controllers.PlanillasController{},
			),
		),
		beego.NSNamespace("/incapacidades",
			beego.NSInclude(
				&controllers.IncapacidadesController{},
			),
		),
		beego.NSNamespace("/utils",
			beego.NSInclude(
				&controllers.UtilsController{},
			),
		),
		beego.NSNamespace("/conceptos_ibc",
			beego.NSInclude(
				&controllers.ConceptosIbcController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
