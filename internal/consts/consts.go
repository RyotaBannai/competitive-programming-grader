package consts

type CPG_ANNOTATIONS struct {
	DIRSPEC     string
	TEST_IGNORE string
	TEST_ALLOW  string
}

type CPG_RUN_TYPES struct {
	COMPILE []string
	SCRIPT  []string
	RUST    []string
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

	RUN_TYPES = CPG_RUN_TYPES{
		COMPILE: []string{"c++", "java", "scala", "haskell", "ocaml", "d", "c", "c#", "f#", "elixir", "nim", "typescript", "swift", "kotlin", "dart", "julia", "groovy", "clojure"},
		SCRIPT:  []string{"python", "ruby", "perl", "php", "javascript", "r"},
		RUST:    []string{"rust"},
	}
)
