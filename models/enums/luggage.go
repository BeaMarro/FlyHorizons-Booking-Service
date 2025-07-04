package enums

import "encoding/json"

type Luggage string

const (
	SmallBag        Luggage = "SmallBag"
	CabinBag        Luggage = "CabinBag"
	Cargo20kg       Luggage = "Cargo20kg"
	Cargo30kg       Luggage = "Cargo30kg"
	SportsEquipment Luggage = "SportsEquipment"
	BabyCarrier     Luggage = "BabyCarrier"
)

func LuggageClassesFromJSONString(jsonInput string) []Luggage {
	var luggageStrings []string
	if err := json.Unmarshal([]byte(jsonInput), &luggageStrings); err != nil {
		return []Luggage{}
	}

	var result []Luggage
	for _, item := range luggageStrings {
		switch Luggage(item) {
		case SmallBag, CabinBag, Cargo20kg, Cargo30kg, SportsEquipment, BabyCarrier:
			result = append(result, Luggage(item))
		}
	}

	return result
}

func JSONStringToLuggageClasses(luggage []Luggage) string {
	jsonData, err := json.Marshal(luggage)
	if err != nil {
		return "[]"
	}
	return string(jsonData)
}
