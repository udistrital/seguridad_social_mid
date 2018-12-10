package models

type PagoSeguridadSocial struct {
	NombrePersona           string
	IdProveedor             int64
	SaludUd                 float64
	SaludTotal              int64
	PensionUd               float64
	PensionTotal            int64
	FondoSolidaridad        float64
	Arl                     int64
	Caja                    int64 //Caja de compensación
	Icbf                    int64 //Instituto Colombiano de Bienestar Familiar
	IdPreliquidacion        int
	IdDetallePreliquidacion int
}
