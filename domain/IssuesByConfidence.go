package domain

type IssuesByConfidence struct {
	Undefined float64 `json:"undefined"`
	Low       float64 `json:"low"`
	Medium    float64 `json:"medium"`
	High      float64 `json:"high"`
}
