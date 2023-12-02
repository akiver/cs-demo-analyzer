package constants

type DemoType string

func (demoType DemoType) String() string {
	return string(demoType)
}

const (
	DemoTypeGOTV DemoType = "GOTV"
	DemoTypePOV  DemoType = "POV"
)
