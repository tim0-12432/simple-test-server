package dtos

type Status int

const (
	Running Status = iota
	Discarded
)

func (s Status) String() string {
	return [...]string{"running", "discarded"}[s]
}

func ToStatus(v any) Status {
	if v == nil {
		return Discarded
	}

	switch t := v.(type) {
	case string:
		switch t {
		case "running":
			return Running
		case "discarded":
			return Discarded
		default:
			return Discarded
		}
	case Status:
		return t
	default:
		return Discarded
	}
}

type Container struct {
	ID          string            `json:"container_id"`
	Name        string            `json:"name"`
	Image       string            `json:"image"`
	CreatedAt   int64             `json:"created_at"`
	Environment map[string]string `json:"environment"`
	Ports       map[int]int       `json:"ports"`
	Volumes     map[string]string `json:"volumes"`
	Networks    []string          `json:"networks"`
	Status      Status            `json:"status"`
	Type        string            `json:"type"`
}

func (c *Container) GetID() string {
	return c.ID
}

func (c *Container) GetName() string {
	return c.Name
}

func (c *Container) GetImage() string {
	return c.Image
}

func (c *Container) GetCreatedAt() int64 {
	return c.CreatedAt
}

func (c *Container) GetEnvironment() map[string]string {
	return c.Environment
}

func (c *Container) GetPorts() map[int]int {
	return c.Ports
}

func (c *Container) GetVolumes() map[string]string {
	return c.Volumes
}

func (c *Container) GetNetworks() []string {
	return c.Networks
}

func (c *Container) GetStatus() Status {
	return c.Status
}

func (c *Container) GetType() string {
	return c.Type
}
