// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/ss_mid_api/controllers"
)

func init() {
	ns := beego.NewNamespace("/v1",

		beego.NSNamespace("/tipo_pago",
			beego.NSInclude(
				&controllers.TipoPagoController{},
			),
		),

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

		beego.NSNamespace("/upc_adicional",
			beego.NSInclude(
				&controllers.UpcAdicionalController{},
			),
		),

		beego.NSNamespace("/zona_upc",
			beego.NSInclude(
				&controllers.ZonaUpcController{},
			),
		),

		beego.NSNamespace("/rango_edad_upc",
			beego.NSInclude(
				&controllers.RangoEdadUpcController{},
			),
		),

		beego.NSNamespace("/tipo_upc",
			beego.NSInclude(
				&controllers.TipoUpcController{},
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
	)
	beego.AddNamespace(ns)
}
