package constants

type WeaponType string

func (name WeaponType) String() string {
	return string(name)
}

const (
	WeaponTypeUnknown    WeaponType = "unknown"
	WeaponTypePistol     WeaponType = "pistol"
	WeaponTypeSMG        WeaponType = "smg"
	WeaponTypeShotgun    WeaponType = "shotgun"
	WeaponTypeRifle      WeaponType = "rifle"
	WeaponTypeSniper     WeaponType = "sniper"
	WeaponTypeMachineGun WeaponType = "machine_gun"
	WeaponTypeGrenade    WeaponType = "grenade"
	WeaponTypeEquipment  WeaponType = "equipment"
	WeaponTypeMelee      WeaponType = "melee"
	WeaponTypeWorld      WeaponType = "world"
)

type WeaponName string

func (name WeaponName) String() string {
	return string(name)
}

const (
	WeaponAK47         WeaponName = "AK-47"
	WeaponAUG          WeaponName = "AUG"
	WeaponAWP          WeaponName = "AWP"
	WeaponBomb         WeaponName = "C4"
	WeaponCZ75         WeaponName = "CZ75 Auto"
	WeaponDecoy        WeaponName = "Decoy Grenade"
	WeaponDeagle       WeaponName = "Desert Eagle"
	WeaponDefuseKit    WeaponName = "Defuse Kit"
	WeaponDualBerettas WeaponName = "Dual Berettas"
	WeaponFamas        WeaponName = "FAMAS"
	WeaponFiveSeven    WeaponName = "Five-SeveN"
	WeaponFlashbang    WeaponName = "Flashbang"
	WeaponG3SG1        WeaponName = "G3SG1"
	WeaponGalilAR      WeaponName = "Galil AR"
	WeaponGlock        WeaponName = "Glock-18"
	WeaponHEGrenade    WeaponName = "HE Grenade"
	WeaponHelmet       WeaponName = "Kevlar + Helmet"
	WeaponKevlar       WeaponName = "Kevlar Vest"
	WeaponIncendiary   WeaponName = "Incendiary Grenade"
	WeaponKnife        WeaponName = "Knife"
	WeaponM249         WeaponName = "M249"
	WeaponM4A1         WeaponName = "M4A1"
	WeaponM4A4         WeaponName = "M4A4"
	WeaponMac10        WeaponName = "MAC-10"
	WeaponMAG7         WeaponName = "MAG-7"
	WeaponMolotov      WeaponName = "Molotov"
	WeaponMP5          WeaponName = "MP5-SD"
	WeaponMP7          WeaponName = "MP7"
	WeaponMP9          WeaponName = "MP9"
	WeaponNegev        WeaponName = "Negev"
	WeaponNova         WeaponName = "Nova"
	WeaponP2000        WeaponName = "P2000"
	WeaponP250         WeaponName = "P250"
	WeaponP90          WeaponName = "P90"
	WeaponPPBizon      WeaponName = "PP-Bizon"
	WeaponRevolver     WeaponName = "R8 Revolver"
	WeaponSawedOff     WeaponName = "Sawed-Off"
	WeaponScar20       WeaponName = "SCAR-20"
	WeaponScout        WeaponName = "SSG 08"
	WeaponSG553        WeaponName = "SG 553"
	WeaponSmoke        WeaponName = "Smoke Grenade"
	WeaponTec9         WeaponName = "Tec-9"
	WeaponUMP45        WeaponName = "UMP-45"
	WeaponUnknown      WeaponName = "Unknown"
	WeaponUSP          WeaponName = "USP-S"
	WeaponWorld        WeaponName = "World"
	WeaponXM1014       WeaponName = "XM1014"
	WeaponZeus         WeaponName = "Zeus x27"
)
