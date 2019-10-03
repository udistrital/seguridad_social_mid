package models

type InformacionProveedor struct {
	Id                      int                `orm:"column(id_proveedor);pk"`
	Tipopersona             string             `orm:"column(tipopersona)"`
	NumDocumento            string             `orm:"column(num_documento)"`
	IdCiudadContacto        float64            `orm:"column(id_ciudad_contacto)"`
	Direccion               string             `orm:"column(direccion)"`
	Correo                  string             `orm:"column(correo)"`
	Web                     string             `orm:"column(web);null"`
	NomAsesor               string             `orm:"column(nom_asesor);null"`
	TelAsesor               string             `orm:"column(tel_asesor);null"`
	Descripcion             string             `orm:"column(descripcion);null"`
	PuntajeEvaluacion       float64            `orm:"column(puntaje_evaluacion);null"`
	ClasificacionEvaluacion string             `orm:"column(clasificacion_evaluacion);null"`
	Estado                  *ParametroEstandar `orm:"column(estado);rel(fk)"`
	TipoCuentaBancaria      string             `orm:"column(tipo_cuenta_bancaria)"`
	NumCuentaBancaria       string             `orm:"column(num_cuenta_bancaria)"`
	IdEntidadBancaria       float64            `orm:"column(id_entidad_bancaria)"`
	FechaRegistro           string             `orm:"column(fecha_registro)"`
	FechaUltimaModificacion string             `orm:"column(fecha_ultima_modificacion)"`
	NomProveedor            string             `orm:"column(nom_proveedor);null"`
	Anexorut                string             `orm:"column(anexorut)"`
	Anexorup                string             `orm:"column(anexorup);null"`
	RegimenContributivo     string             `orm:"column(regimen_contributivo);null"`
}
