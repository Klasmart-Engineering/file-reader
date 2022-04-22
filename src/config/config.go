package config

type Config struct {
	Logger Logger
	Kafka  Kafka
	Server Server
}

// Server config
type Server struct {
	Port        string
	Development bool
}

// Logger config
type Logger struct {
	DisableCaller     bool
	DisableStacktrace bool
	Encoding          string
	Level             string
}

type Kafka struct {
	Brokers                []string
	MinBytes               int // 10e3 (10KB)
	MaxBytes               int // 10e6 (10MB)
	QueueCapacity          int // 100
	HeartbeatInterval      int //3 * time.Second
	CommitInterval         int // 0
	PartitionWatchInterval int // 5 * time.Second
	MaxAttempts            int // 3
	DialTimeout            int //3 * time.Minute

	WriterReadTimeout  int // 10 * time.Second
	WriterWriteTimeout int // 10 * time.Second
	WriterRequiredAcks int // -1
	WriterMaxAttempts  int // 3

	CreateOrganizationTopic   string // "organization"
	CreateOrganizationWorkers int    // 3

	DeadLetterQueueTopic string // "dead-letter-queue"

	OrganizationsGroupID string // "organizations_group"
}
