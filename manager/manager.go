package manager

import (
	"errors"
	"fmt"
	"os"

	"github.com/drone/drone-go/drone"
	"github.com/uphy/drone-util/model"
	"golang.org/x/oauth2"
)

type Manager struct {
	client drone.Client
}

func NewFromEnv() (*Manager, error) {
	host := os.Getenv("DRONE_SERVER")
	if host == "" {
		return nil, errors.New("DRONE_SERVER not set")
	}
	token := os.Getenv("DRONE_TOKEN")
	if token == "" {
		return nil, errors.New("DRONE_TOKEN not set")
	}
	return New(host, token), nil
}

func New(host string, token string) *Manager {
	// create an http client with oauth authentication.
	config := new(oauth2.Config)
	auther := config.Client(
		oauth2.NoContext,
		&oauth2.Token{
			AccessToken: token,
		},
	)
	// create the drone client with authenticator
	client := drone.NewClient(host, auther)
	return &Manager{client}
}

func (m *Manager) Export() (model.Repos, error) {
	repoList, err := m.client.RepoList()
	if err != nil {
		return nil, err
	}
	repos := model.Repos{}
	for _, r := range repoList {
		secretList, err := m.client.SecretList(r.Owner, r.Name)
		if err != nil {
			return nil, err
		}
		secrets := map[string]model.Secret{}
		for _, secret := range secretList {
			events := []model.Event{}
			for _, e := range secret.Events {
				events = append(events, model.Event(e))
			}
			secrets[secret.Name] = model.Secret{
				Name:   secret.Name,
				Value:  secret.Value,
				Events: events,
			}
		}
		repos[r.FullName] = model.Repo{
			Hooks: model.Hooks{
				Push:        &r.AllowPush,
				PullRequest: &r.AllowPull,
				Tag:         &r.AllowTag,
				Deployment:  &r.AllowDeploy,
			},
			Settings: model.Settings{
				Protected: &r.IsGated,
				Trusted:   &r.IsTrusted,
			},
			Secrets: secrets,
			Timeout: &r.Timeout,
		}
	}
	return repos, nil
}

func (m *Manager) Import(repos model.Repos) error {
	for name, repo := range repos {
		ownerName, repoName := model.SplitOwnerAndRepoName(name)
		p := drone.RepoPatch{}
		p.AllowDeploy = repo.Hooks.Deployment
		p.AllowPull = repo.Hooks.PullRequest
		p.AllowPush = repo.Hooks.Push
		p.AllowTag = repo.Hooks.Tag
		p.IsGated = repo.Settings.Protected
		p.IsTrusted = repo.Settings.Trusted
		p.Timeout = repo.Timeout
		if repo.Visibility != nil {
			s := string(*(repo.Visibility))
			p.Visibility = &s
		}
		_, err := m.client.RepoPatch(ownerName, repoName, &p)
		if err != nil {
			return fmt.Errorf("Failed to update repo. (name=%s, err=%v)", name, err)
		}
		secretList, err := m.client.SecretList(ownerName, repoName)
		if err != nil {
			return fmt.Errorf("Failed to get secret list: %v", err)
		}
		secretMap := map[string]drone.Secret{}
		for _, secret := range secretList {
			secretMap[secret.Name] = *secret
		}
		for _, secret := range repo.Secrets {
			events := []string{}
			for _, e := range secret.Events {
				events = append(events, string(e))
			}
			converted := &drone.Secret{
				Name:   secret.Name,
				Value:  secret.Value,
				Images: nil,
				Events: events,
			}
			if _, exist := secretMap[secret.Name]; exist {
				if _, err := m.client.SecretUpdate(ownerName, repoName, converted); err != nil {
					return fmt.Errorf("Failed to create secret. (name=%s, secretName=%s, err=%v)", name, converted.Name, err)
				}
			} else {
				if _, err := m.client.SecretCreate(ownerName, repoName, converted); err != nil {
					return fmt.Errorf("Failed to update secret. (name=%s, secretName=%s, err=%v)", name, converted.Name, err)
				}
			}
		}

	}
	return nil
}
