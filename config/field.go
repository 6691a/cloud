package config

type Server struct {
	Debug     bool                     `yaml:"debug"`
	SentryDsn string                   `yaml:"sentry_dsn"`
	Logging   map[string]LoggingConfig `yaml:"logging"`
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

type Hypervisor struct {
	Service    string  `yaml:"service"`
	WorkerSize uint8   `yaml:"worker_size"`
	BufferSize uint8   `yaml:"buffer_size"`
	Proxmox    Proxmox `yaml:"proxmox"`
}

type Proxmox struct {
	Url   string `yaml:"url"`
	Node  string `yaml:"node"`
	User  string `yaml:"user"`
	Token string `yaml:"token"`
}

type LoggingConfig struct {
	Level        string `yaml:"level"`
	Path         string `yaml:"path"`
	MaxAge       int    `yaml:"max_age"`
	RotationTime int    `yaml:"rotation_time"`
}
