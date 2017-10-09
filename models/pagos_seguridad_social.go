package models

type PagosSeguridadSocial struct {
	NumeroContrato          string
	SaludUd                 float64
	SaludTotal              int64
	PensionUd               float64
	PensionTotal            int64
	Arl                     int64
	Caja                    int64 //Caja de compensación
	Icbf                    int64 //Instituto Colombiano de Bienestar Familiar
	IdDetallePreliquidacion int
	UpcAdicional            []UpcAdicional
}
