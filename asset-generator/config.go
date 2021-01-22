package main

type Config struct {
	Url          string `yaml:"url"`
	TenantId     string `yaml:"tenantId"`
	NestId       string `yaml:"nestId"`
	TemplateId   string `yaml:"templateId"`
	CatalogId    string `yaml:"catalogId"`
	ParamGroupId string `yaml:"paramGroupId"`
	DeviceId     string `yaml:"deviceId"`
}
