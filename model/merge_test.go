package model

import (
	"bytes"
	"testing"
)

var (
	True             = true
	False            = false
	I100       int64 = 100
	I200       int64 = 200
	VisPublic        = VisibilityPublic
	VisPrivate       = VisibilityPrivate
)

func TestMerge(t *testing.T) {
	testMerge("deployment", t, Repo{
		Hooks: Hooks{
			Deployment: &True,
		},
	}, Repo{
		Hooks: Hooks{},
	}, Repo{
		Hooks: Hooks{
			Deployment: &True,
		},
	})
	testMerge("deployment override", t, Repo{
		Hooks: Hooks{
			Deployment: &True,
		},
	}, Repo{
		Hooks: Hooks{
			Deployment: &False,
		},
	}, Repo{
		Hooks: Hooks{
			Deployment: &False,
		},
	})
	testMerge("push", t, Repo{
		Hooks: Hooks{
			Push: &True,
		},
	}, Repo{
		Hooks: Hooks{
			Push: &False,
		},
	}, Repo{
		Hooks: Hooks{
			Push: &False,
		},
	})
	testMerge("pullrequest", t, Repo{
		Hooks: Hooks{
			PullRequest: &True,
		},
	}, Repo{
		Hooks: Hooks{
			PullRequest: &False,
		},
	}, Repo{
		Hooks: Hooks{
			PullRequest: &False,
		},
	})
	testMerge("tag", t, Repo{
		Hooks: Hooks{
			Tag: &True,
		},
	}, Repo{
		Hooks: Hooks{
			Tag: &False,
		},
	}, Repo{
		Hooks: Hooks{
			Tag: &False,
		},
	})
	testMerge("secret", t, Repo{
		Secrets: map[string]Secret{
			"a": Secret{
				Name:  "a",
				Value: "A",
			},
			"c": Secret{
				Name:  "c",
				Value: "C",
			},
		},
	}, Repo{
		Secrets: map[string]Secret{
			"a": Secret{
				Name:  "a",
				Value: "A2",
			},
			"b": Secret{
				Name:  "b",
				Value: "B",
			},
		},
	}, Repo{
		Secrets: map[string]Secret{
			"a": Secret{
				Name:  "a",
				Value: "A2",
			},
			"b": Secret{
				Name:  "b",
				Value: "B",
			},
			"c": Secret{
				Name:  "c",
				Value: "C",
			},
		},
	})
	testMerge("protected", t, Repo{
		Settings: Settings{
			Protected: &True,
		},
	}, Repo{
		Settings: Settings{
			Protected: &False,
		},
	}, Repo{
		Settings: Settings{
			Protected: &False,
		},
	})
	testMerge("trusted", t, Repo{
		Settings: Settings{
			Trusted: &True,
		},
	}, Repo{
		Settings: Settings{
			Trusted: &False,
		},
	}, Repo{
		Settings: Settings{
			Trusted: &False,
		},
	})
	testMerge("timeout", t, Repo{
		Timeout: &I100,
	}, Repo{
		Timeout: &I200,
	}, Repo{
		Timeout: &I200,
	})
	testMerge("visibility", t, Repo{
		Visibility: &VisPrivate,
	}, Repo{
		Visibility: &VisPublic,
	}, Repo{
		Visibility: &VisPublic,
	})
}

func testMerge(name string, t *testing.T, base, target, want Repo) {
	m := Merge(base, target)

	s1 := repoToString(want)
	s2 := repoToString(*m)
	if s1 != s2 {
		t.Errorf("[%s] want %s but %s", name, s1, s2)
	}
}

func repoToString(r Repo) string {
	buf := new(bytes.Buffer)
	write(buf, r)
	return buf.String()
}
