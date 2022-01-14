package rules

type RuleType string

const (
	None            RuleType = "None"
	Time            RuleType = "Time"
	Location        RuleType = "Location"
	FaceRecognize   RuleType = "FaceRecognize"
	ObjectRecognize RuleType = "ObjectRecognize"
)
