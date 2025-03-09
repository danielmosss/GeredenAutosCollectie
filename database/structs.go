package database

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// RDWCarResponse represents the structure of the RDW API response
type RDWCarResponse struct {
	Kenteken                                      string  `json:"kenteken"`
	Voertuigsoort                                 string  `json:"voertuigsoort"`
	Merk                                          string  `json:"merk"`
	Handelsbenaming                               string  `json:"handelsbenaming"`
	VervaldatumAPK                                string  `json:"vervaldatum_apk"`
	DatumTenaamstelling                           string  `json:"datum_tenaamstelling"`
	BrutoBPM                                      string  `json:"bruto_bpm"`
	Inrichting                                    string  `json:"inrichting"`
	AantalZitplaatsen                             string  `json:"aantal_zitplaatsen"`
	EersteKleur                                   string  `json:"eerste_kleur"`
	TweedeKleur                                   string  `json:"tweede_kleur"`
	AantalCilinders                               string  `json:"aantal_cilinders"`
	Cilinderinhoud                                string  `json:"cilinderinhoud"`
	MassaLedigVoertuig                            string  `json:"massa_ledig_voertuig"`
	ToegestaneMaximumMassaVoertuig                string  `json:"toegestane_maximum_massa_voertuig"`
	MassaRijklaar                                 string  `json:"massa_rijklaar"`
	MaximumMassaTrekkenOngeremd                   string  `json:"maximum_massa_trekken_ongeremd"`
	MaximumTrekkenMassaGeremd                     string  `json:"maximum_trekken_massa_geremd"`
	DatumEersteToelating                          string  `json:"datum_eerste_toelating"`
	DatumEersteTenaamstellingInNederland          string  `json:"datum_eerste_tenaamstelling_in_nederland"`
	WachtOpKeuren                                 string  `json:"wacht_op_keuren"`
	Catalogusprijs                                string  `json:"catalogusprijs"`
	WAMVerzekerd                                  string  `json:"wam_verzekerd"`
	AantalDeuren                                  string  `json:"aantal_deuren"`
	AantalWielen                                  string  `json:"aantal_wielen"`
	Lengte                                        string  `json:"lengte"`
	EuropeseVoertuigcategorie                     string  `json:"europese_voertuigcategorie"`
	PlaatsChassisnummer                           string  `json:"plaats_chassisnummer"`
	TechnischeMaxMassaVoertuig                    string  `json:"technische_max_massa_voertuig"`
	Typegoedkeuringsnummer                        string  `json:"typegoedkeuringsnummer"`
	Variant                                       string  `json:"variant"`
	Uitvoering                                    string  `json:"uitvoering"`
	VolgnummerWijzigingEUTypegoedkeuring          string  `json:"volgnummer_wijziging_eu_typegoedkeuring"`
	VermogenMassarijklaar                         float64 `json:"vermogen_massarijklaar,string"`
	Wielbasis                                     string  `json:"wielbasis"`
	ExportIndicator                               string  `json:"export_indicator"`
	OpenstaandeTerugroepactieIndicator            string  `json:"openstaande_terugroepactie_indicator"`
	TaxiIndicator                                 string  `json:"taxi_indicator"`
	MaximumMassaSamenstelling                     string  `json:"maximum_massa_samenstelling"`
	JaarLaatsteRegistratieTellerstand             string  `json:"jaar_laatste_registratie_tellerstand"`
	Tellerstandoordeel                            string  `json:"tellerstandoordeel"`
	CodeToelichtingTellerstandoordeel             string  `json:"code_toelichting_tellerstandoordeel"`
	TenaamstellenMogelijk                         string  `json:"tenaamstellen_mogelijk"`
	VervaldatumAPKDT                              string  `json:"vervaldatum_apk_dt"`
	DatumTenaamstellingDT                         string  `json:"datum_tenaamstelling_dt"`
	DatumEersteToelatingDT                        string  `json:"datum_eerste_toelating_dt"`
	DatumEersteTenaamstellingInNederlandDT        string  `json:"datum_eerste_tenaamstelling_in_nederland_dt"`
	Zuinigheidsclassificatie                      string  `json:"zuinigheidsclassificatie"`
	APIGekentekendeVoertuigenAssen                string  `json:"api_gekentekende_voertuigen_assen"`
	APIGekentekendeVoertuigenBrandstof            string  `json:"api_gekentekende_voertuigen_brandstof"`
	APIGekentekendeVoertuigenCarrosserie          string  `json:"api_gekentekende_voertuigen_carrosserie"`
	APIGekentekendeVoertuigenCarrosserieSpecifiek string  `json:"api_gekentekende_voertuigen_carrosserie_specifiek"`
	APIGekentekendeVoertuigenVoertuigklasse       string  `json:"api_gekentekende_voertuigen_voertuigklasse"`
}

type Car struct {
	ID                       primitive.ObjectID `bson:"_id,omitempty"`
	Kenteken                 string             `bson:"kenteken"`
	Merk                     string             `bson:"merk"`
	Handelsbenaming          string             `bson:"handelsbenaming"`
	Variant                  string             `bson:"variant"`
	Uitvoering               string             `bson:"uitvoering"`
	Inrichting               string             `bson:"inrichting"`
	EersteKleur              string             `bson:"eerste_kleur"`
	TweedeKleur              string             `bson:"tweede_kleur"`
	AantalZitplaatsen        int                `bson:"aantal_zitplaatsen"`
	AantalDeuren             int                `bson:"aantal_deuren"`
	AantalWielen             int                `bson:"aantal_wielen"`
	AantalCilinders          int                `bson:"aantal_cilinders"`
	Cilinderinhoud           int                `bson:"cilinderinhoud"`
	Catalogusprijs           int                `bson:"catalogusprijs"`
	Lengte                   int                `bson:"lengte"`
	Wielbasis                int                `bson:"wielbasis"`
	DatumEersteToelating     string             `bson:"datum_eerste_toelating"`
	Zuinigheidsclassificatie string             `bson:"zuinigheidsclassificatie"`
	Picture                  string             `bson:"picture,omitempty"`
	DrivenDate               time.Time          `bson:"driven_date"`
}
