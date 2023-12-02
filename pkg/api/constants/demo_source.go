package constants

type DemoSource string

func (source DemoSource) String() string {
	return string(source)
}

const (
	DemoSourceUnknown    DemoSource = "unknown"
	DemoSourceValve      DemoSource = "valve"
	DemoSourceESEA       DemoSource = "esea"
	DemoSourceFaceIt     DemoSource = "faceit"
	DemoSourceEbot       DemoSource = "ebot"
	DemoSourceCEVO       DemoSource = "cevo"
	DemoSourceESL        DemoSource = "esl"
	DemoSourcePopFlash   DemoSource = "popflash"
	DemoSourceEsportal   DemoSource = "esportal"
	DemoSourceFastcup    DemoSource = "fastcup"
	DemoSourceGamersclub DemoSource = "gamersclub"
)

var SupportedDemoSources = []DemoSource{
	DemoSourceValve,
	DemoSourceESEA,
	DemoSourceFaceIt,
	DemoSourceEbot,
	DemoSourceESL,
	DemoSourcePopFlash,
}
