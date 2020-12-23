# Discosay

Build auto-responding Discord bots.

# Usage

All configuration is done via environment variables.

* `CONFIG_PATH`: Point this to a YAML config file for Discosay.
* `DEBUG`: Set to enable debug-level logging.
* `DEVELOPMENT`: Set to enable colorized human-readable output.

## Config Structure

The config file is structured as shown:

```yaml
templates:
  danger: |
    ####################
    DANGER DANGER DANGER
    ####################
    $MSG
    ####################
    DANGER DANGER DANGER
    ####################

bots:
  rolldice:
    - rolld6
  lostinspace:
    - dangerwill
    - dangersay

responders:
  rolld6:
    match: "^!rolld6$"
    responses:
      - You rolled a 1
      - You rolled a 2
      - You rolled a 3
      - You rolled a 4
      - You rolled a 5
      - You rolled a 6
  dangersay:
    match: "^!dangersay (.+)$"
    template: danger
  dangerwill:
    match: '\bdanger\b'
    probability: 0.1
    template: danger
    responses:
      - Danger, Will Robinson!
```

* `templates`: A `dict[str: str]` of named templates that responders can optionally use.
* `bots`: A `dict[str: array[str]]` of bots and the named responders they use.
* `responders`: A `dict[str, dict]` of responders that bots use.
  * `match`: A regex string. If this regex matches a message, the responder will reply in the channel. If this regex has a capture group, it will be used as the reply.
  * `responses`: Optional. An `array[str]` of possible responses. If provided, the responder selects one of these as the reply.
  * `template`: Optional. String. If provided, the reply message is injected into this template in place of the string `$MSG`.
  * `probability`: Optional. Float from 0.0 to 1.0. Defaults to 1.0 (100%). If provided, this is the probability this responder will send a reply for a message it matches.