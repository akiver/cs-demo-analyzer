package api

import (
	"github.com/akiver/cs-demo-analyzer/pkg/api/constants"
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/common"
)

func getEquipmentWeaponType(equipment common.Equipment) constants.WeaponType {
	switch equipment.Type {
	case common.EqCZ, common.EqDeagle, common.EqDualBerettas, common.EqFiveSeven, common.EqGlock, common.EqP2000, common.EqP250, common.EqRevolver, common.EqTec9, common.EqUSP:
		return constants.WeaponTypePistol
	case common.EqAK47, common.EqAUG, common.EqFamas, common.EqGalil, common.EqM4A1, common.EqM4A4, common.EqSG553:
		return constants.WeaponTypeRifle
	case common.EqAWP, common.EqG3SG1, common.EqScar20, common.EqScout:
		return constants.WeaponTypeSniper
	case common.EqMac10, common.EqMP5, common.EqMP7, common.EqMP9, common.EqP90, common.EqBizon, common.EqUMP:
		return constants.WeaponTypeSMG
	case common.EqSwag7, common.EqNova, common.EqSawedOff, common.EqXM1014:
		return constants.WeaponTypeShotgun
	case common.EqM249, common.EqNegev:
		return constants.WeaponTypeMachineGun
	case common.EqDecoy, common.EqFlash, common.EqHE, common.EqIncendiary, common.EqMolotov, common.EqSmoke:
		return constants.WeaponTypeGrenade
	case common.EqBomb, common.EqDefuseKit, common.EqKevlar, common.EqHelmet:
		return constants.WeaponTypeEquipment
	case common.EqKnife, common.EqZeus:
		return constants.WeaponTypeMelee
	case common.EqWorld:
		return constants.WeaponTypeWorld
	}

	return constants.WeaponTypeUnknown
}

var equipmentToWeaponName = map[common.EquipmentType]constants.WeaponName{
	common.EqAK47:         constants.WeaponAK47,
	common.EqAUG:          constants.WeaponAUG,
	common.EqAWP:          constants.WeaponAWP,
	common.EqBomb:         constants.WeaponBomb,
	common.EqCZ:           constants.WeaponCZ75,
	common.EqDecoy:        constants.WeaponDecoy,
	common.EqDefuseKit:    constants.WeaponDefuseKit,
	common.EqDeagle:       constants.WeaponDeagle,
	common.EqDualBerettas: constants.WeaponDualBerettas,
	common.EqFamas:        constants.WeaponFamas,
	common.EqFiveSeven:    constants.WeaponFiveSeven,
	common.EqFlash:        constants.WeaponFlashbang,
	common.EqG3SG1:        constants.WeaponG3SG1,
	common.EqGalil:        constants.WeaponGalilAR,
	common.EqGlock:        constants.WeaponGlock,
	common.EqHE:           constants.WeaponHEGrenade,
	common.EqKevlar:       constants.WeaponKevlar,
	common.EqHelmet:       constants.WeaponHelmet,
	common.EqKnife:        constants.WeaponKnife,
	common.EqIncendiary:   constants.WeaponIncendiary,
	common.EqM249:         constants.WeaponM249,
	common.EqM4A1:         constants.WeaponM4A1,
	common.EqM4A4:         constants.WeaponM4A4,
	common.EqMac10:        constants.WeaponMac10,
	common.EqSwag7:        constants.WeaponMAG7,
	common.EqMolotov:      constants.WeaponMolotov,
	common.EqMP5:          constants.WeaponMP5,
	common.EqMP7:          constants.WeaponMP7,
	common.EqMP9:          constants.WeaponMP9,
	common.EqNegev:        constants.WeaponNegev,
	common.EqNova:         constants.WeaponNova,
	common.EqP2000:        constants.WeaponP2000,
	common.EqP250:         constants.WeaponP250,
	common.EqP90:          constants.WeaponP90,
	common.EqBizon:        constants.WeaponPPBizon,
	common.EqRevolver:     constants.WeaponRevolver,
	common.EqSawedOff:     constants.WeaponSawedOff,
	common.EqScar20:       constants.WeaponScar20,
	common.EqSG553:        constants.WeaponSG553,
	common.EqSmoke:        constants.WeaponSmoke,
	common.EqScout:        constants.WeaponScout,
	common.EqTec9:         constants.WeaponTec9,
	common.EqUMP:          constants.WeaponUMP45,
	common.EqUnknown:      constants.WeaponUnknown,
	common.EqUSP:          constants.WeaponUSP,
	common.EqWorld:        constants.WeaponWorld,
	common.EqXM1014:       constants.WeaponXM1014,
	common.EqZeus:         constants.WeaponZeus,
}
