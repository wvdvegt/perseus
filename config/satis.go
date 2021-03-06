package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

// Satis (https://github.com/composer/satis) is a simple and static Composer repository generator.
// Once medusa is finished it writes a satis file itself with a (kind of) dynamic content.
// The content of the satis file is driven by the medusa config file + the results of the operations.
// Operations are downloading packages or resolving dependency trees of those packages.
// Once they are successful, the package will land into a satis file.
//
// The satis file itself is a typical JSON file.
// It has various settings and is (nearly) complete documented via a schema
// available at https://github.com/composer/satis/blob/master/res/satis-schema.json.
// One proper solution would be to reflect this schema as a struct with JSON tags and read/write
// those values via the struct.
// In this file we do it differently. The reason is simple: We don't want to change perseus
// once satis added a new feature into the JSON schema.
// Perseus/Medusa is only interested to write the `repositories` section.
// So we only modify this.
// For implementation details checkout the WriteFile() function.

// Satis reflects the a Satis configuration file.
type Satis struct {
	// config is the configuration provider object that has read the satis configuration file
	config Provider
	// List of repositories
	repositories map[string]SatisRepository
}

// SatisRepository reflects a single repository entry in satis `repositories` section
type SatisRepository struct {
	// Type is the repository type, like `git` or `svn`
	Type string `json:"type"`
	// URL is the URL of the repository that contains packages
	URL string `json:"url"`
}

// NewSatis will create a new satis configuration object.
// If no configuration is given, an error will be returned.
func NewSatis(c Provider) (*Satis, error) {
	if c == nil {
		return nil, errors.New("No configurations provider applied")
	}

	// Read initial repositories
	repositories := []SatisRepository{}
	if v := c.Get("repositories").(*json.RawMessage); v != nil {
		err := json.Unmarshal(*v, &repositories)
		if err != nil {
			return nil, err
		}
	}

	// Deduplicate the repositories
	m := map[string]SatisRepository{}
	for _, v := range repositories {
		m[v.URL] = v
	}

	s := &Satis{
		config:       c,
		repositories: m,
	}
	return s, nil
}

// AddRepository will add repository u to the current satis configuration
func (s *Satis) AddRepository(u string) {
	r := SatisRepository{
		Type: "git",
		URL:  u,
	}

	s.repositories[r.URL] = r
}

// AddRepositories will add a list of repositories u to the current satis configuration
func (s *Satis) AddRepositories(u ...string) {
	for _, r := range u {
		s.AddRepository(r)
	}
}

// WriteFile will write the satis configuration to file filename with permissions perm
func (s *Satis) WriteFile(filename string, perm os.FileMode) error {
	// We maintain the Satis configuration file on our own.
	// This is not managed by viper.
	// Maybe it make sense to switch this in feature.
	// Viper is not able to write configuration files (yet).
	// A PR is available for this. See https://github.com/spf13/viper/pull/287

	contentMap := s.config.GetContentMap()
	m := make(map[string]*json.RawMessage, len(contentMap))
	for k, v := range contentMap {
		t := v.(json.RawMessage)
		m[k] = &t
	}

	repositories := s.GetRepositoriesAsSlice()
	rawRepositories, err := json.MarshalIndent(&repositories, "", "    ")
	if err != nil {
		return err
	}

	jsonRepositories := json.RawMessage(rawRepositories)
	m["repositories"] = &jsonRepositories

	b, err := json.MarshalIndent(&m, "", "    ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, b, perm)
}

// GetRepositoriesAsSlice returns all configured repositories
// from the configuration as a list.
func (s *Satis) GetRepositoriesAsSlice() []SatisRepository {
	c := make([]SatisRepository, 0, len(s.repositories))
	for _, value := range s.repositories {
		c = append(c, value)
	}

	return c
}
