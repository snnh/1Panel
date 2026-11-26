package dto

// GPUEnvironmentCheck reports the result of a GPU runtime check. The
// community edition only keeps the Docker NVIDIA runtime probe.
type GPUEnvironmentCheck struct {
	Vendor  string `json:"vendor"`
	Item    string `json:"item"`
	Status  string `json:"status"`
	Code    string `json:"code,omitempty"`
	Version string `json:"version,omitempty"`
	Reason  string `json:"reason,omitempty"`
}
