package models


type PagosSeguridadSocial struct {
  Id                int64
  Nombre            string
  SaludPersona      float64
  SaludUd           float64
  SaludTotal        float64
  PensionPersona    float64
  PensionUd         float64
  PensionTotal      float64
  Arl               float64
}
