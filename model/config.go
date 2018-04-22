package model

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

type (
	Event string

	Visibility string

	Secret struct {
		Name   string  `yaml:"-"`
		Value  string  `yaml:"value"`
		Events []Event `yaml:"events,omitempty"`
	}

	Hooks struct {
		Push        *bool `yaml:"push,omitempty"`
		PullRequest *bool `yaml:"pullrequest,omitempty"`
		Tag         *bool `yaml:"tag,omitempty"`
		Deployment  *bool `yaml:"deployment,omitempty"`
	}

	Settings struct {
		Protected *bool `yaml:"protected,omitempty"`
		Trusted   *bool `yaml:"trusted,omitempty"`
	}

	Repos map[string]Repo

	Repo struct {
		Secrets    map[string]Secret `yaml:"secrets,omitempty"`
		Settings   Settings          `yaml:"settings,omitempty"`
		Visibility *Visibility       `yaml:"visibility,omitempty"`
		Hooks      Hooks             `yaml:"hooks,omitempty"`
		Timeout    *int64            `yaml:"timeout,omitempty"`
	}

	Config struct {
		Global Repo  `yaml:"global,omitempty"`
		Owners Repos `yaml:"owners,omitempty"`
		Repos  Repos `yaml:"repos,omitempty"`
	}
)

const (
	EventPush          Event      = "push"
	EventTag           Event      = "tag"
	EventDeployment    Event      = "deployment"
	EventPullRequest   Event      = "pull_request"
	VisibilityPublic   Visibility = "public"
	VisibilityPrivate  Visibility = "private"
	VisibilityInternal Visibility = "internal"
)

func ParseFile(file string) (*Config, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	// apply template
	s, err := ApplyTemplate(string(b))
	if err != nil {
		return nil, fmt.Errorf("failed to apply the template: %v", err)
	}
	// parse yaml
	var c Config
	if err := yaml.Unmarshal([]byte(s), &c); err != nil {
		return nil, err
	}
	// completion
	if err := c.fix(); err != nil {
		return nil, err
	}
	return &c, nil
}

func (c *Config) fix() error {
	c.Global.fix()
	for _, r := range c.Owners {
		r.fix()
	}
	for _, r := range c.Repos {
		r.fix()
	}
	return nil
}

func (r *Repo) fix() error {
	for name, s := range r.Secrets {
		s.Name = name
		r.Secrets[name] = s
	}
	return nil
}

func (c *Config) Resolve() Repos {
	owners := map[string]Repo{}
	for name, repo := range c.Owners {
		owners[name] = *Merge(c.Global, repo)
	}
	repos := Repos{}
	for name, repo := range c.Repos {
		ownerName, _ := SplitOwnerAndRepoName(name)
		owner, exist := owners[ownerName]
		if exist {
			repos[name] = *Merge(owner, repo)
		} else {
			repos[name] = *Merge(c.Global, repo)
		}
	}
	return repos
}

func SplitOwnerAndRepoName(name string) (owner, repo string) {
	i := strings.Index(name, "/")
	if i < 0 {
		return "", name
	}
	return name[0:i], name[i+1:]
}

func write(w io.Writer, v interface{}) error {
	b, err := yaml.Marshal(v)
	if err != nil {
		return err
	}
	if _, err := io.WriteString(w, string(b)); err != nil {
		return err
	}
	return nil
}

func (r Repos) Write(w io.Writer) error {
	return write(w, r)
}

func (c *Config) Write(w io.Writer) error {
	return write(w, c)
}

func (s *Secret) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var value string
	if err := unmarshal(&value); err == nil {
		s.Value = value
		s.Events = []Event{EventPush, EventTag}
		return nil
	}
	var m map[string]interface{}
	if err := unmarshal(&m); err == nil {
		value, exist := m["value"]
		if !exist {
			return fmt.Errorf("'value' is required")
		}
		s.Value = value.(string)

		events, exist := m["events"]
		if exist {
			e1, ok := events.([]interface{})
			if !ok {
				return fmt.Errorf("invalid 'events' value. events must be a string array")
			}
			e2 := []Event{}
			for _, e := range e1 {
				e2 = append(e2, Event(fmt.Sprint(e)))
			}
			s.Events = e2
		}
	}
	return nil
}

func (s *Secret) MarshalYAML() (interface{}, error) {
	return s, nil
}
