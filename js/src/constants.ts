export const Game = {
  CSGO: 'CSGO',
  CS2: 'CS2',
  CS2LT: 'CS2 LT',
} as const;
export type Game = (typeof Game)[keyof typeof Game];

export const DemoSource = {
  Unknown: 'unknown',
  Valve: 'valve',
  Ebot: 'ebot',
  Popflash: 'popflash',
  FaceIt: 'faceit',
  Cevo: 'cevo',
  Challengermode: 'challengermode',
  Esl: 'esl',
  Esea: 'esea',
  Esportal: 'esportal',
  Fastcup: 'fastcup',
  Gamersclub: 'gamersclub',
  PerfectWorld: 'perfectworld',
} as const;
export type DemoSource = (typeof DemoSource)[keyof typeof DemoSource];

export const SupportedDemoSources: DemoSource[] = [
  DemoSource.Valve,
  DemoSource.Esea,
  DemoSource.FaceIt,
  DemoSource.Ebot,
  DemoSource.Esl,
  DemoSource.Popflash,
  DemoSource.Challengermode,
  DemoSource.PerfectWorld,
];

export const DemoType = {
  POV: 'POV',
  GOTV: 'GOTV',
} as const;
export type DemoType = (typeof DemoType)[keyof typeof DemoType];

export const ExportFormat = {
  CSV: 'csv',
  JSON: 'json',
  CSDM: 'csdm', // Special CSV export dedicated to the application CS Demo Manager
} as const;
export type ExportFormat = (typeof ExportFormat)[keyof typeof ExportFormat];

export const TeamNumber = {
  UNASSIGNED: 0,
  SPECTATOR: 1,
  T: 2,
  CT: 3,
} as const;
export type TeamNumber = (typeof TeamNumber)[keyof typeof TeamNumber];

export const TeamLetter = {
  A: 'A',
  B: 'B',
} as const;
export type TeamLetter = (typeof TeamLetter)[keyof typeof TeamLetter];

export const WeaponType = {
  Unknown: 'unknown',
  Pistol: 'pistol',
  SMG: 'smg',
  Shotgun: 'shotgun',
  Rifle: 'rifle',
  Sniper: 'sniper',
  MachineGun: 'machine_gun',
  Grenade: 'grenade',
  Equipment: 'equipment',
  Melee: 'melee',
  World: 'world',
} as const;
export type WeaponType = (typeof WeaponType)[keyof typeof WeaponType];

export const WeaponName = {
  AK47: 'AK-47',
  AUG: 'AUG',
  AWP: 'AWP',
  Bomb: 'C4',
  CZ75: 'CZ75 Auto',
  Decoy: 'Decoy Grenade',
  Deagle: 'Desert Eagle',
  DefuseKit: 'Defuse Kit',
  DualBerettas: 'Dual Berettas',
  Famas: 'FAMAS',
  FiveSeven: 'Five-SeveN',
  Flashbang: 'Flashbang',
  G3SG1: 'G3SG1',
  GalilAR: 'Galil AR',
  Glock: 'Glock-18',
  HEGrenade: 'HE Grenade',
  Helmet: 'Kevlar + Helmet',
  Kevlar: 'Kevlar Vest',
  Incendiary: 'Incendiary Grenade',
  Knife: 'Knife',
  M249: 'M249',
  M4A1: 'M4A1',
  M4A4: 'M4A4',
  Mac10: 'MAC-10',
  MAG7: 'MAG-7',
  Molotov: 'Molotov',
  MP5: 'MP5-SD',
  MP7: 'MP7',
  MP9: 'MP9',
  Negev: 'Negev',
  Nova: 'Nova',
  P2000: 'P2000',
  P250: 'P250',
  P90: 'P90',
  PPBizon: 'PP-Bizon',
  Revolver: 'R8 Revolver',
  SawedOff: 'Sawed-Off',
  Scar20: 'SCAR-20',
  Scout: 'SSG 08',
  SG553: 'SG 553',
  Smoke: 'Smoke Grenade',
  Tec9: 'Tec-9',
  UMP45: 'UMP-45',
  Unknown: 'Unknown',
  USP: 'USP-S',
  World: 'World',
  XM1014: 'XM1014',
  Zeus: 'Zeus x27',
} as const;
export type WeaponName = (typeof WeaponName)[keyof typeof WeaponName];

/**
 * Values come from https://github.com/alliedmodders/sourcemod/blob/master/plugins/include/cstrike.inc#L53
 */
export const RoundEndReason = {
  StillInProgress: 0 /* Round not over */,
  TargetBombed: 1 /* Target Successfully Bombed! */,
  VipEscaped: 2 /* The VIP has escaped! - Doesn't exist on CS:GO */,
  VipKilled: 3 /* VIP has been assassinated! - Doesn't exist on CS:GO */,
  TerroristsEscaped: 4 /* The terrorists have escaped! */,
  CounterTerroristsStoppedEscape: 5 /* The CTs have prevented most of the terrorists from escaping! */,
  TerroristsStopped: 6 /* Escaping terrorists have all been neutralized! */,
  BombDefused: 7 /* The bomb has been defused! */,
  CtWin: 8 /* Counter-Terrorists Win! */,
  TerroristWin: 9 /* Terrorists Win! */,
  Draw: 10 /* Round Draw! */,
  HostagesRescued: 11 /* All Hostages have been rescued! */,
  TargetSaved: 12 /* Target has been saved! */,
  HostagesNotRescued: 13 /* Hostages have not been rescued! */,
  TerroristsNotEscaped: 14 /* Terrorists have not escaped! */,
  VipNotEscaped: 15 /* VIP has not escaped! - Doesn't exist on CS:GO */,
  GameStart: 16 /* Game Commencing! */,

  // The below only exist on CS:GO
  TerroristsSurrender: 17 /* Terrorists Surrender */,
  CounterTerroristsSurrender: 18 /* CTs Surrender */,
  TerroristsPlanted: 19 /* Terrorists Planted the bomb */,
  CounterTerroristsReachedHostage: 20 /* CTs Reached the hostage */,
} as const;
export type RoundEndReason = (typeof RoundEndReason)[keyof typeof RoundEndReason];

export const PlayerColor = {
  Grey: -1,
  Yellow: 0,
  Purple: 1,
  Green: 2,
  Blue: 3,
  Orange: 4,
} as const;
export type PlayerColor = (typeof PlayerColor)[keyof typeof PlayerColor];

export const HitGroup = {
  Generic: 0,
  Head: 1,
  Chest: 2,
  Stomach: 3,
  LeftArm: 4,
  RightArm: 5,
  LeftLeg: 6,
  RightLeg: 7,
  Gear: 10,
} as const;
export type HitGroup = (typeof HitGroup)[keyof typeof HitGroup];

export const GrenadeName = {
  Smoke: WeaponName.Smoke,
  Flashbang: WeaponName.Flashbang,
  HE: WeaponName.HEGrenade,
  Decoy: WeaponName.Decoy,
  Molotov: WeaponName.Molotov,
  Incendiary: WeaponName.Incendiary,
} as const;
export type GrenadeName = (typeof GrenadeName)[keyof typeof GrenadeName];

export const EconomyType = {
  Pistol: 'pistol',
  Eco: 'eco',
  Semi: 'semi',
  ForceBuy: 'force-buy',
  Full: 'full',
} as const;
export type EconomyType = (typeof EconomyType)[keyof typeof EconomyType];

/**
 * The game type and game mode are related.
 * Depending of the game type, the game mode will have different values.
 * Values come from the file "gamemodes.txt" located inside pak01_dir.vpk of CS2.
 */
export const GameType = {
  Classic: 0,
  GunGame: 1,
  Training: 2,
  Custom: 3,
  CoOperative: 4,
  Skirmish: 5,
  FFA: 6,
} as const;
export type GameType = (typeof GameType)[keyof typeof GameType];

/**
 * When game type is Classic.
 */
export const GameModeClassic = {
  Casual: 0,
  Competitive: 1,
  Scrimmage2V2: 2,
  Scrimmage5V5: 3,
} as const;
export type GameModeClassic = (typeof GameModeClassic)[keyof typeof GameModeClassic];

/**
 * When game type is GunGame.
 */
export const GameModeGunGame = {
  Progressive: 0,
  Bomb: 1,
  Deathmatch: 2,
} as const;
export type GameModeGunGame = (typeof GameModeGunGame)[keyof typeof GameModeGunGame];

/**
 * When game type is Training.
 */
export const GameModeTraining = {
  Training: 0,
} as const;
export type GameModeTraining = (typeof GameModeTraining)[keyof typeof GameModeTraining];

/**
 * When game type is Custom.
 */
export const GameModeCustom = {
  Custom: 0,
} as const;
export type GameModeCustom = (typeof GameModeCustom)[keyof typeof GameModeCustom];

/**
 * When game type is CoOperative.
 */
export const GameModeCoOperative = {
  CoOperative: 0,
  CoOperativeMission: 1,
} as const;
export type GameModeCoOperative = (typeof GameModeCoOperative)[keyof typeof GameModeCoOperative];

/**
 * When game type is Skirmish.
 */
export const GameModeSkirmish = {
  Skirmish: 0,
} as const;
export type GameModeSkirmish = (typeof GameModeSkirmish)[keyof typeof GameModeSkirmish];

/**
 * When game type is FFA.
 */
export const GameModeSurvival = {
  Survival: 0,
} as const;
export type GameModeSurvival = (typeof GameModeSurvival)[keyof typeof GameModeSurvival];

export const HostageState = {
  Idle: 0,
  BeingUntied: 1,
  GettingPickedUp: 2,
  BeingCarried: 3,
  FollowingPlayer: 4,
  GettingDropped: 5,
  Rescued: 6,
  Dead: 7,
} as const;
export type HostageState = (typeof HostageState)[keyof typeof HostageState];

export const CsgoRankType = {
  Unknown: -1,
  None: 0,
  ClassicCompetitive: 6,
  Wingman2v2: 7,
  DangerZone: 10,
} as const;
export type CsgoRankType = (typeof CsgoRankType)[keyof typeof CsgoRankType];

export const Cs2RankType = {
  Unknown: -1,
  None: 0,
  Wingman2v2: 7,
  PremierMode: 11,
  ClassicCompetitive: 12,
} as const;
export type Cs2RankType = (typeof Cs2RankType)[keyof typeof Cs2RankType];

export type RankType = CsgoRankType | Cs2RankType;

export const CompetitiveRank = {
  Unknown: 0,
  SilverI: 1,
  SilverII: 2,
  SilverIII: 3,
  SilverIV: 4,
  SilverElite: 5,
  SilverEliteMaster: 6,
  GoldNovaI: 7,
  GoldNovaII: 8,
  GoldNovaIII: 9,
  GoldNovaMaster: 10,
  MasterGuardianI: 11,
  MasterGuardianII: 12,
  MasterGuardianElite: 13,
  DistinguishedMasterGuardian: 14,
  LegendaryEagle: 15,
  LegendaryEagleMaster: 16,
  SupremeMasterFirstClass: 17,
  GlobalElite: 18,
} as const;
export type CompetitiveRank = (typeof CompetitiveRank)[keyof typeof CompetitiveRank];
export type PremierRank = number;

export type Rank = CompetitiveRank | PremierRank;

// Game mode as a string reported in CSVCMsg_ServerInfo messages.
export const GameMode = {
  Casual: 'casual',
  Premier: 'premier',
  Competitive: 'competitive',
  Scrimmage2V2: 'scrimcomp2v2',
  Scrimmage5V5: 'scrimcomp5v5',
  Deathmatch: 'deathmatch',
  GunGameProgressive: 'gungameprogressive',
  GunGameBomb: 'gungametrbomb',
  Custom: 'custom',
  CoOperative: 'cooperative',
  CoOperativeMission: 'coopmission',
  Skirmish: 'skirmish',
  Survival: 'survival',
} as const;
export type GameMode = (typeof GameMode)[keyof typeof GameMode];
