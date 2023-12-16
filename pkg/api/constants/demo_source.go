package constants

type DemoSource string

func (source DemoSource) String() string {
	return string(source)
}

const (
	DemoSourceUnknown        DemoSource = "unknown"
	DemoSourceValve          DemoSource = "valve"
	DemoSourceESEA           DemoSource = "esea"
	DemoSourceFaceIt         DemoSource = "faceit"
	DemoSourceEbot           DemoSource = "ebot"
	DemoSourceCEVO           DemoSource = "cevo"
	DemoSourceChallengermode DemoSource = "challengermode"
	DemoSourceESL            DemoSource = "esl"
	DemoSourcePopFlash       DemoSource = "popflash"
	DemoSourceEsportal       DemoSource = "esportal"
	DemoSourceFastcup        DemoSource = "fastcup"
	DemoSourceGamersclub     DemoSource = "gamersclub"
	// "Perfect World" (完美世界) is a Chinese company that Valve partnered with to release CS:GO in China.
	DemoSourcePerfectWorld DemoSource = "perfectworld"
)

var SupportedDemoSources = []DemoSource{
	DemoSourceValve,
	DemoSourceESEA,
	DemoSourceFaceIt,
	DemoSourceEbot,
	DemoSourceESL,
	DemoSourcePopFlash,
	DemoSourceChallengermode,
	DemoSourcePerfectWorld,
}
