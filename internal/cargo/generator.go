package cargo

import (
	"github.com/pivotal-cf/kiln/internal/cargo/opsman"
	"github.com/pivotal-cf/kiln/internal/proofing"
)

type OpsManagerConfig struct {
	DeploymentName    string
	AvailabilityZones []string
	Stemcells         []opsman.Stemcell
	ResourceConfigs   []opsman.ResourceConfig
}

type Generator struct{}

func NewGenerator() Generator {
	return Generator{}
}

func (g Generator) Execute(template proofing.ProductTemplate, config OpsManagerConfig) Manifest {
	var releases []Release
	for _, release := range template.Releases {
		releases = append(releases, Release{
			Name:    release.Name,
			Version: release.Version,
		})
	}

	var stemcell Stemcell
	for _, boshStemcell := range config.Stemcells {
		if boshStemcell.OS == template.StemcellCriteria.OS {
			if boshStemcell.Version == template.StemcellCriteria.Version {
				stemcell = Stemcell{
					Alias:   boshStemcell.Name,
					OS:      boshStemcell.OS,
					Version: boshStemcell.Version,
				}
			}
		}
	}

	update := Update{
		Canaries:        1,
		CanaryWatchTime: "30000-300000",
		UpdateWatchTime: "30000-300000",
		MaxInFlight:     1,
		MaxErrors:       2,
		Serial:          template.Serial,
	}

	var instanceGroups []InstanceGroup
	for _, jobType := range template.JobTypes {
		lifecycle := "service"
		if jobType.Errand {
			lifecycle = "errand"
		}

		instances := jobType.InstanceDefinition.Default
		for _, resourceConfig := range config.ResourceConfigs {
			if resourceConfig.Name == jobType.Name {
				if !resourceConfig.Instances.IsAutomatic() {
					instances = resourceConfig.Instances.Value
				}
			}
		}

		instanceGroups = append(instanceGroups, InstanceGroup{
			Name:      jobType.Name,
			AZs:       config.AvailabilityZones,
			Lifecycle: lifecycle,
			Stemcell:  stemcell.Alias,
			Instances: instances,
		})
	}

	var variables []Variable
	for _, variable := range template.Variables {
		variables = append(variables, Variable{
			Name:    variable.Name,
			Options: variable.Options,
			Type:    variable.Type,
		})
	}

	return Manifest{
		Name:           config.DeploymentName,
		Releases:       releases,
		Stemcells:      []Stemcell{stemcell},
		Update:         update,
		Variables:      variables,
		InstanceGroups: instanceGroups,
	}
}
