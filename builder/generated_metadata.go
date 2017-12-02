package builder

type GeneratedMetadata struct {
	Name             string
	StemcellCriteria StemcellCriteria
	Releases         []Release
	IconImage        string

	FormTypes      []Part
	JobTypes       []Part
	RuntimeConfigs []Part
	Variables      []Part

	Metadata Metadata
}

type Release struct {
	Name    string
	File    string
	Version string
}

type StemcellCriteria struct {
	Version     string
	OS          string
	RequiresCPI bool `yaml:"requires_cpi"`
}

func (gm GeneratedMetadata) MarshalYAML() (interface{}, error) {
	m := map[string]interface{}{}

	m["name"] = gm.Name
	m["stemcell_criteria"] = gm.StemcellCriteria
	m["releases"] = gm.Releases
	m["icon_image"] = gm.IconImage

	if len(gm.FormTypes) > 0 {
		m["form_types"] = gm.metadataOnly(gm.FormTypes)
	}
	if len(gm.JobTypes) > 0 {
		m["job_types"] = gm.metadataOnly(gm.JobTypes)
	}
	if len(gm.RuntimeConfigs) > 0 {
		m["runtime_configs"] = gm.metadataOnly(gm.RuntimeConfigs)
	}
	if len(gm.Variables) > 0 {
		m["variables"] = gm.metadataOnly(gm.Variables)
	}

	for k, v := range gm.Metadata {
		m[k] = v
	}

	return m, nil
}

func (gm GeneratedMetadata) metadataOnly(parts []Part) []interface{} {
	metadata := []interface{}{}
	for _, p := range parts {
		metadata = append(metadata, p.Metadata)
	}
	return metadata
}
