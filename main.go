package main

import (
	"encoding/json"
	"fmt"
	"math"
	_ "math"
	"net/http"
)

// Input struct for JSON decoding
type Input struct {
	Values []float64 `json:"values"`
}

type ElectricPowerUnit struct {
	EfficiencyFactor         float64
	LoadFactor               float64
	NominalPower             float64
	UnitCount                float64
	ActualPower              float64
	VariationCoefficient     float64
	ReactivePowerCoefficient float64
	Rs1rivn                  float64
}

type EPStorage struct {
	NomPotEp    float64
	KoefVik     float64
	KoefReakPot float64

	GroupKVik    float64
	EffCount     float64
	RozKAkP      float64
	RozAkNav     float64
	SumKvNP      float64
	SumKvNPTg    float64
	RozReakNav   float64
	PovnPotu     float64
	RozGroupSt   float64
	ForgroupKVik float64
	ForeffCount  float64

	Eps map[string]*ElectricPowerUnit
}

func NewEPStorage(nomPotEp, koefVik, koefReakPot float64) *EPStorage {
	eps := map[string]*ElectricPowerUnit{
		"grinding":  {0.92, 0.9, 0.38, 4.0, 20.0, 0.15, 1.33, 0.0},
		"drilling":  {0.92, 0.9, 0.38, 2.0, 14.0, 0.12, 1.0, 0.0},
		"fuguval":   {0.92, 0.9, 0.38, 4.0, 42.0, 0.15, 1.33, 0.0},
		"circular":  {0.92, 0.9, 0.38, 1.0, 36.0, 0.3, 1.52, 0.0},
		"press":     {0.92, 0.9, 0.38, 1.0, 20.0, 0.5, 0.75, 0.0},
		"polishing": {0.92, 0.9, 0.38, 1.0, 40.0, 0.2, 1.0, 0.0},
		"milling":   {0.92, 0.9, 0.38, 2.0, 32.0, 0.2, 1.0, 0.0},
		"fan":       {0.92, 0.9, 0.38, 1.0, 20.0, 0.65, 0.75, 0.0},
	}

	return &EPStorage{
		NomPotEp:    nomPotEp,
		KoefVik:     koefVik,
		KoefReakPot: koefReakPot,
		RozKAkP:     1.25,
		Eps:         eps,
	}
}

func (e *EPStorage) CalculateGroup() {
	e.GroupKVik, e.SumKvNP, e.SumKvNPTg, e.EffCount, e.ForgroupKVik, e.ForeffCount = 0, 0, 0, 0, 0, 0

	for _, value := range e.Eps {
		e.SumKvNP += (value.VariationCoefficient * value.ActualPower) * value.UnitCount
		e.SumKvNPTg += (value.VariationCoefficient * value.ActualPower) * value.ReactivePowerCoefficient * value.UnitCount
		e.ForgroupKVik += (value.UnitCount * value.ActualPower)
		e.ForeffCount += (value.UnitCount * value.ActualPower * value.ActualPower)
	}

	e.GroupKVik = e.SumKvNP / e.ForgroupKVik
	e.EffCount = e.ForgroupKVik * e.ForgroupKVik / e.ForeffCount
	e.RozAkNav = e.RozKAkP * e.SumKvNP
	e.RozReakNav = e.RozKAkP * e.SumKvNPTg
	e.PovnPotu = math.Sqrt(e.RozAkNav*e.RozAkNav + e.RozReakNav*e.RozReakNav)
	e.RozGroupSt = e.RozAkNav / e.NomPotEp
}

func calculateTask1(nomPotEp, koefVik, koefReakPot float64) string {
	ePStorage := NewEPStorage(nomPotEp, koefVik, koefReakPot)
	ePStorage.CalculateGroup()

	cehPn := 2330
	cehPn2 := 96399
	cehnPnKv := 752
	cehRKAP := 0.7
	cehKvPntg := 657

	cehGroupKVik := float64(cehnPnKv) / float64(cehPn)
	cehNe := float64(cehPn*cehPn) / float64(cehPn2)
	cehRAN := cehRKAP * float64(cehnPnKv)
	cehRRN := cehRKAP * float64(cehKvPntg)
	cehPP := math.Sqrt(cehRAN*cehRAN + cehRRN*cehRRN)
	cehRGS := cehRAN / 0.38

	output := fmt.Sprintf(
		"Груповий коефіцієнт використання для ШР1=ШР2=ШР3: %.4f\n"+
			"Ефективна кількість ЕП для ШР1=ШР2=ШР3: %.4f\n"+
			"Розрахунковий коефіцієнт активної потужності для ШР1=ШР2=ШР3: %.4f\n"+
			"Розрахункове активне навантаження для ШР1=ШР2=ШР3: %.4f\n"+
			"Розрахункове реактивне навантаження для ШР1=ШР2=ШР3: %.4f\n"+
			"Повна потужність для ШР1=ШР2=ШР3: %.4f\n"+
			"Розрахунковий груповий струм для ШР1=ШР2=ШР3: %.4f\n"+
			"Коефіцієнти використання цеху в цілому: %.4f\n"+
			"Ефективна кількість ЕП цеху в цілому: %.4f\n"+
			"Розрахунковий коефіцієнт активної потужності цеху в цілому: %.4f\n"+
			"Розрахункове активне навантаження на шинах 0,38 кВ ТП: %.4f\n"+
			"Розрахункове реактивне навантаження на шинах 0,38 кВ ТП: %.4f\n"+
			"Повна потужність на шинах 0,38 кВ ТП: %.4f\n"+
			"Розрахунковий груповий струм на шинах 0,38 кВ ТП: %.4f\n",
		ePStorage.GroupKVik, ePStorage.EffCount, ePStorage.RozKAkP, ePStorage.RozAkNav, ePStorage.RozReakNav, ePStorage.PovnPotu, ePStorage.RozGroupSt,
		cehGroupKVik, cehNe, cehRKAP, cehRAN, cehRRN, cehPP, cehRGS)

	return output
}

func calculator1Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var input Input
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if len(input.Values) != 3 {
		http.Error(w, "Invalid number of inputs", http.StatusBadRequest)
		return
	}
	result := calculateTask1(input.Values[0], input.Values[1], input.Values[2])

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"result": result})
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/api/calculator1", calculator1Handler)

	fmt.Println("Server running at http://localhost:8086")
	http.ListenAndServe(":8086", nil)
}
