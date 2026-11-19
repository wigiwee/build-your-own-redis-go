package config

type ServerConfig struct {
	Port        int
	Version     string
	Mode        string
	BindAddress string

	ReplicationConfig ReplicationConfig
}

type ReplicationConfig struct {
	Role                           string
	Connected_slaves               string
	Master_replid                  string
	Master_repl_offset             int
	Second_repl_offset             int
	Repl_backlog_active            string
	Repl_backlog_size              int
	Repl_backlog_first_byte_offset int
	Repl_backlog_histlen           int
	Master_host                    string
	Master_port                    string
}
