package Configs

type Parames struct {
	HostName       string `yaml:"domain"`
	Port           uint16 `yaml:"port"`
	AdminPort      uint16 `yaml:"admin_port"`
	Pathes         Pathes `yaml:"folders"`
	Mail           Mail   `yaml:"mail"`
	TLS            TLS    `yaml:"TLS"`
	normalshutdown bool   `yaml:"-"`
	istls          bool   `yaml:"-"`
}

type Pathes struct {
	UserDir   string `yaml:"user"`
	DataBases string `yaml:"databases"`
	Media     string `yaml:"media"`
	DialogDir string `yaml:"dialogs"`
}

type Mail struct {
	PathRSA      string `yaml:"namersa"`
	TemplateHTML string `yaml:"namehtmltemplate"`
}

type TLS struct {
	Certs string `yaml:"pathcerts"`
	Key   string `yaml:"key"`
}

type DefaultValues struct {
	TrustCleint bool `yaml:"trustclient"`
}
