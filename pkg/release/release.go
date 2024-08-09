package release

type Release struct {
	APIVersion     float64    `yaml:"apiVersion"`
	ReleaseVersion string     `yaml:"releaseVersion"`
	Components     Components `yaml:"components,omitempty"`
}

type Components struct {
	Kubernetes      Kubernetes      `yaml:"kubernetes"`
	OperatingSystem OperatingSystem `yaml:"operatingSystem"`
	Workloads       []HelmChart     `yaml:"workloads"`
}

type HelmChart struct {
	ReleaseName      string      `yaml:"releaseName"`
	Name             string      `yaml:"chart"`
	Repository       string      `yaml:"repository,omitempty"`
	Version          string      `yaml:"version"`
	PrettyName       string      `yaml:"prettyName,omitempty"`
	DependencyCharts []HelmChart `yaml:"dependencyCharts,omitempty"`
	AddonCharts      []HelmChart `yaml:"addonCharts,omitempty"`
}

type Kubernetes struct {
	K3S  KubernetesDistribution `yaml:"k3s"`
	RKE2 KubernetesDistribution `yaml:"rke2"`
}

type KubernetesDistribution struct {
	Version string `yaml:"version"`
}

type OperatingSystem struct {
	Version        string   `yaml:"version"`
	ZypperID       string   `yaml:"zypperID"`
	CPEScheme      string   `yaml:"cpeScheme"`
	RepoGPGPath    string   `yaml:"repoGPGPath"`
	SupportedArchs []string `yaml:"supportedArchs"`
	PrettyName     string   `yaml:"prettyName"`
}
