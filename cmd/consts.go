package cmd

type CPG_ANNOTATIONS struct {
	DIRSPEC string
}

const (
	APP_CONF_EV           = "CPG_CONF_PATH"
	NIL                   = "<nil>"
	DEFAULT_CONF_FILENAME = "cpg_conf"
	DEFAULT_CONF_DIR      = "."
)

var (
	ANNOTATIONS = CPG_ANNOTATIONS{DIRSPEC: "@cpg_dirspec"}
)
