package model

// Merge merges the `target` to `base` repo.
func Merge(base Repo, target Repo) *Repo {
	r := base
	if target.Hooks.Deployment != nil {
		r.Hooks.Deployment = target.Hooks.Deployment
	}
	if target.Hooks.PullRequest != nil {
		r.Hooks.PullRequest = target.Hooks.PullRequest
	}
	if target.Hooks.Push != nil {
		r.Hooks.Push = target.Hooks.Push
	}
	if target.Hooks.Tag != nil {
		r.Hooks.Tag = target.Hooks.Tag
	}
	r.Secrets = mergeSecret(base.Secrets, target.Secrets)
	if target.Settings.Protected != nil {
		r.Settings.Protected = target.Settings.Protected
	}
	if target.Settings.Trusted != nil {
		r.Settings.Trusted = target.Settings.Trusted
	}
	if target.Timeout != nil {
		r.Timeout = target.Timeout
	}
	if target.Visibility != nil {
		r.Visibility = target.Visibility
	}
	return &r
}

func mergeSecret(base map[string]Secret, target map[string]Secret) map[string]Secret {
	s := map[string]Secret{}
	for k, v := range base {
		s[k] = v
	}
	for k, v := range target {
		s[k] = v
	}
	return s
}
