package consts

type CPG_ANNOTATIONS struct {
	DIRSPEC string
}

const (
	APP_CONF_EV                  = "CPG_CONF_PATH"
	NIL                          = "<nil>"
	DEFAULT_CONF_FILENAME        = "cpg_conf.toml"
	DEFAULT_CONF_DIR             = "."
	COMMAND_FILENAME_PLACEHOLDER = "{FILENAME}"
)

var (
	ANNOTATIONS = CPG_ANNOTATIONS{DIRSPEC: "@cpg_dirspec"}
)
