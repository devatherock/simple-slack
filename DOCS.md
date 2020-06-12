## Config

The following parameters/secrets can be set to configure the plugin.

### Parameters
* **color** - Color in which the message block will be highlighted.
* **text** - The message content

### Secrets

The following secret values can be set to configure the plugin.

* **SLACK_WEBHOOK** - The slack webhook to post the message to

## Examples

**Drone:**

```yaml
pipeline:
  notify_slack:
    when:
      event: [ push ]
      status: [ success, failure, error ]
    image: devatherock/simple-slack:0.3.0
    secrets: [ slack_webhook ]
    color: "#33ad7f"
    text: |-
       {{.DroneBuildStatus}} {{.DroneBuildLink}} ({{.DroneCommitRef}}) by {{DroneCommitAuthor}}
       {{.DroneCommitMessage}}
```

**Vela:**

```yaml
steps:
  - name: notify_slack
    ruleset:
      event: [ push ]
      status: [ success ]
    image: devatherock/simple-slack:0.3.0
    secrets: [ slack_webhook ]
    parameters:
      color: "#33ad7f"
      text: |-
        Success: {{.BuildLink}} ({{.BuildRef}}) by {{.BuildAuthor}}
        {{.BuildMessage}}
```  