package consts

type CPG_ANNOTATIONS struct {
	DIRSPEC     string
	TEST_IGNORE string
	TEST_ALLOW  string
}

const (
	NIL                          = "<nil>"
	APP_NAME                     = "cpg"           // unchangeable
	APP_CONF_EV                  = "CPG_CONF_PATH" // unchangeable
	APP_CONF_FILENAME            = "cpg_conf.toml" // unchangeable
	DEFAULT_CONF_DIR             = "."             // changeable
	COMMAND_FILENAME_PLACEHOLDER = "{FILENAME}"    // unchangeable
	EXECUTABLE_NAME              = "a.out"         // unchangeable
)

var (
	ANNOTATIONS = CPG_ANNOTATIONS{
		DIRSPEC:     "@cpg_dirspec",
		TEST_IGNORE: "@cpg_test_ignore",
		TEST_ALLOW:  "@cpg_test_allow"}
)
