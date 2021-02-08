package main

type Config struct {
	TenantId   string `yaml:"tenantId"`
	FilePath   string `yaml:"path"`
	CatalogId  string `yaml:"catalogId"`
	TemplateId string `yaml:"templateId"`
	TypeId     string `yaml:"typeId"`
	NestId     string `yaml:"nestId"`
	ParamGroup string `yaml:"paramGroupId"`
}
