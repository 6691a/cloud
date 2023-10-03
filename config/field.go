package config

type Server struct {
	Debug     bool   `yaml:"debug"`
	SentryDsn string `yaml:"sentry_dsn"`
}

type DNS struct {
	Domain     string `yaml:"domain"`
	Service    string `yaml:"service"`
	WorkerSize uint8  `yaml:"worker_size"`
	BufferSize uint8  `yaml:"buffer_size"`
	GCP        GCP    `yaml:"gcp"`
}

type GCP struct {
	ProjectId      string `yaml:"project_id"`
	ManagedZone    string `yaml:"managed_zone"`
	CredentialPath string `yaml:"credential_path"`
}

type VM struct {
	Service    string  `yaml:"service"`
	WorkerSize uint8   `yaml:"worker_size"`
	bufferSize uint8   `yaml:"buffer_size"`
	Proxmox    Proxmox `yaml:"proxmox"`
}

type Proxmox struct {
	Url string `yaml:"url"`
}

type Router struct {
	Service    string   `yaml:"service"`
	WorkerSize uint8    `yaml:"worker_size"`
	BufferSize uint8    `yaml:"buffer_size"`
	RouterOS   RouterOS `yaml:"routeros"`
}

type RouterOS struct {
	Url      string `yaml:"url"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}
