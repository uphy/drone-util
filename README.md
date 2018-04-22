# Drone Utilities

Drone Utilities(`drone-util`) provides tools for [Drone](https://drone.io/). 

## Installation

Get the binary from the [GitHub releases page](https://github.com/uphy/drone-util/releases) and extract it.  
You can find `drone-util`(or `drone-util.exe`).

Also you can use Docker.

```bash
$ docker run --rm \
    -e DRONE_SERVER=http://<HOSTNAME>:8000/ \
    -e DRONE_TOKEN=<DRONETOKEN> \
    uphy/drone-util \
    <SUBCOMMAND>
```

## General

`drone-util` is a command with subcommands.
Currentry there are the following subcommands.

* import
* export

There's common environment variables `drone-util` use.  
You must set them before the execution.

* `DRONE_HOST`: the Drone server URL
* `DRONE_TOKEN`: the Drone server token. (See http://docs.drone.io/api-authentication/)

## Import/Export

Import/Export repository configuration. (e.g., secrets, trusted/protected, visibility...)  

### Export

Export all of the repository configuration.

```bash
$ drone-util export
repos:
  user1/test:
    secrets:
      foo:
        value: ""
        events:
        - push
...
```

Note: The `value` of the `secrets` will not be displayed because we cannot get from API.

### Import

Import the repository configuration.

Create a configuration file,

```bash
$ cat << EOF > config.yml
repos:
  user1/test:
    secrets:
      http_proxy:
        value: "http://proxy.example.com:8080/"
        events:
        - push
        - tag
        - deployment
    settings:
      protected: true
      trusted: false
    hooks:
      push: true
      pullrequest: true
      tag: false
      deployment: false
    timeout: 60
EOF
```

and import it.

```bash
$ drone-util import config.yml
```

You don't need to write all configuration.  You can write some of the configuration you want to update.
Same as above example, you can setup multiple repository configuration at once.

#### Secrets

`secrets` has 2 styles.

```yaml
secrets:
  a: AAA
```

equals to

```yaml
secrets:
  a:
    value: AAA
    events:
    - push
    - tag
    - deployment     
```

#### Scopes

`drone-util`'s configuratin file has 3 scopes, `global`/`owners`/`repos`.  
`repos` inherits `owners`, and `owners` inherits `global`.

For example:

```yaml
global:
  secrets:
    A: 1
    B: "aaa"
owners:
  # inherits `global`
  user1:
    secrets:
      B: "bbb"
repos:
  # inherits `owners/user1`
  user1/repo1:
    secrets:
      C: "ccc"
```

The actual settings can be checked with `--dry-run`(`-d`) option.

```bash
$ drone-util import --dry-run settings.yml
user1/repo1:
  secrets:
    A:
      value: "1"
      events:
      - push
      - tag
    B:
      value: bbb
      events:
      - push
      - tag
    C:
      value: ccc
      events:
      - push
      - tag
```

#### Template

In the configuration file, you can use the environment variables.

```yaml
global:
  secrets:
    HTTP_PROXY: "{{ env "HTTP_PROXY" }}"
```

Template format is [here](https://golang.org/pkg/text/template/)

