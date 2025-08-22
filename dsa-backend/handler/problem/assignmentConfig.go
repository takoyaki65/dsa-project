package problem

type AssignmentConfig struct {
	SubID         int        `json:"sub_id"`
	Title         string     `json:"title"`
	MDfile        string     `json:"md_file"`
	TimeMS        *int       `json:"time_ms,omitempty"`
	MemoryMB      *int       `json:"memory_mb,omitempty"`
	TestFiles     []string   `json:"test_files"`
	RequiredFiles []string   `json:"required_files"`
	Build         []TestCase `json:"build"`
	Judge         []TestCase `json:"judge"`
}

type TestCase struct {
	EvalOnly      *bool  `json:"eval_only,omitempty"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	MessageOnFail string `json:"message_on_fail,omitempty"`
	Command       string `json:"command"`
	Stdin         string `json:"stdin,omitempty"`
	Stdout        string `json:"stdout,omitempty"`
	Stderr        string `json:"stderr,omitempty"`
	ExitCode      *int   `json:"exit,omitempty"`
}

func (conf *AssignmentConfig) SetDefaults() {
	if conf.TimeMS == nil {
		defaultTime := 1000 // Default time in milliseconds
		conf.TimeMS = &defaultTime
	}
	if conf.MemoryMB == nil {
		defaultMemory := 256 // Default memory in MB
		conf.MemoryMB = &defaultMemory
	}

	for i := range conf.Build {
		conf.Build[i].SetDefaults()
	}

	for i := range conf.Judge {
		conf.Judge[i].SetDefaults()
	}
}

func (t *TestCase) SetDefaults() {
	if t.EvalOnly == nil {
		defaultEvalOnly := false
		t.EvalOnly = &defaultEvalOnly
	}
	if t.ExitCode == nil {
		defaultExitCode := 0
		t.ExitCode = &defaultExitCode
	}
	if t.MessageOnFail == "" {
		t.MessageOnFail = "failed to execute " + t.Title
	}
}
