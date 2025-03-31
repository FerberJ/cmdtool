package models

type Config struct {
	ConfigFile   ConfigFile
	RunIndexList []uint
}

type ConfigFile struct {
	RemoteUser string            `yaml:"remoteUser"`
	RemoteHost string            `yaml:"remoteHost"`
	IdRsaPath  string            `yaml:"idRsaPath"`
	Variables  map[string]string `yaml:"variables"`
	RunCmds    []RunCmd          `yaml:"runCmds"`
	Cmds       map[string]Cmd    `yaml:"cmds"`
}

type RunCmd struct {
	Cmd           string   `yaml:"cmd"`
	RunIndex      uint     `yaml:"runIndex"`
	Params        []string `yaml:"params"`
	Description   string   `yaml:"description"`
	StopAfterFail bool     `yaml:"stopAfterFail"`
}

type Cmd struct {
	Cmd    string   `yaml:"cmd"`    // example "scp"
	Type   string   `yaml:"type"`   // exec or ssh. example: "exec"
	Args   []string `yaml:"args"`   // example: ["-r"]
	Params uint     `yaml:"params"` // example: 2 (runCmd will have ["{{localfile}}", "{{remoteUser}}@{{remoteHost}}:{{remotePath}}"])
}
