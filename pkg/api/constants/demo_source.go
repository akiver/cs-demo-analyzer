package constants

type DemoSource string

func (source DemoSource) String() string {
	return string(source)
}

const (
	DemoSourceCEVO           DemoSource = "cevo"
	DemoSourceChallengermode DemoSource = "challengermode"
	DemoSourceEbot           DemoSource = "ebot"
	DemoSourceESEA           DemoSource = "esea"
	DemoSourceESL            DemoSource = "esl"
	DemoSourceEsportal       DemoSource = "esportal"
	DemoSourceFaceIt         DemoSource = "faceit"
	DemoSourceFastcup        DemoSource = "fastcup"
	DemoSourceFiveEPlay      DemoSource = "5eplay"
	DemoSourceGamersclub     DemoSource = "gamersclub"
	// "Perfect World" (完美世界) is a Chinese company that Valve partnered with to release CS:GO in China.
	DemoSourcePerfectWorld DemoSource = "perfectworld"
	DemoSourcePopFlash     DemoSource = "popflash"
	DemoSourceUnknown      DemoSource = "unknown"
	DemoSourceValve        DemoSource = "valve"
)

var SupportedDemoSources = []DemoSource{
	DemoSourceChallengermode,
	DemoSourceEbot,
	DemoSourceESEA,
	DemoSourceESL,
	DemoSourceFaceIt,
	DemoSourceFastcup,
	DemoSourceFiveEPlay,
	DemoSourcePerfectWorld,
	DemoSourcePopFlash,
	DemoSourceValve,
	DemoSourceEsportal,
}
