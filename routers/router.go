// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"github.com/udistrital/ss_mid_api/controllers"

	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("/v1",

		beego.NSNamespace("/upc_adicional",
			beego.NSInclude(
				&controllers.UpcAdicionalController{},
			),
		),

		beego.NSNamespace("/desc_seguridad_social",
			beego.NSInclude(
				&controllers.DescSeguridadSocialController{},
			),
		),

		beego.NSNamespace("/edad_upc",
			beego.NSInclude(
				&controllers.EdadUpcController{},
			),
		),

		beego.NSNamespace("/tipo_zona_upc",
			beego.NSInclude(
				&controllers.TipoZonaUpcController{},
			),
		),

		beego.NSNamespace("/tipo_pago_seguridad_social",
			beego.NSInclude(
				&controllers.TipoPagoSeguridadSocialController{},
			),
		),

		beego.NSNamespace("/tipo_upc",
			beego.NSInclude(
				&controllers.TipoUpcController{},
			),
		),

		beego.NSNamespace("/desc_seguridad_social_detalle",
			beego.NSInclude(
				&controllers.DescSeguridadSocialDetalleController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
