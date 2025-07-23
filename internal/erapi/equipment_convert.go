package erapi

// Equipment Code = 0 weapon, 1 clothes, 2 hat, 3 arm, 4 leg

var ItemTypeKeyMapping = map[string]string{
	"0": "weapon",
	"1": "cloth",
	"2": "hat",
	"3": "arm",
	"4": "leg",
}

// Grade Code = 0 common, 1 uncommon, 2 rare, 3 epic, 4 legendary, 5 mythic
var ItemGradeKeyMapping = map[string]string{
	"0": "common",
	"1": "uncommon",
	"2": "rare",
	"3": "epic",
	"4": "legendary",
	"5": "mythic",
}

func convertKeys(input UserGame) UserGame {
	converted := input

	// Equipment 키 변환
	converted.Equipment = map[string]int{}
	for oldKey, value := range input.Equipment {
		if newKey, exists := ItemTypeKeyMapping[oldKey]; exists {
			converted.Equipment[newKey] = value
		} else {
			converted.Equipment[oldKey] = value // 매핑 없는 키는 원본 유지
		}
	}

	// EquipmentGrade 키 변환
	converted.EquipmentGrade = map[string]int{}
	for oldKey, value := range input.EquipmentGrade {
		if newKey, exists := ItemTypeKeyMapping[oldKey]; exists {
			converted.EquipmentGrade[newKey] = value
		} else {
			converted.EquipmentGrade[oldKey] = value // 매핑 없는 키는 원본 유지
		}
	}

	converted.EquipFirstItemForLog = map[string][]int{}
	for oldKey, value := range input.EquipFirstItemForLog {
		if newKey, exists := ItemTypeKeyMapping[oldKey]; exists {
			converted.EquipFirstItemForLog[newKey] = value
		} else {
			converted.EquipFirstItemForLog[oldKey] = value // 매핑 없는 키는 원본 유지
		}
	}

	return converted
}
