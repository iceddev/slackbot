package robots

import (
  "fmt"
  "bytes"
  "regexp"
  "strings"
  "net/http"
  "io/ioutil"
  "encoding/json"
)

type PollBot struct {
}

func init() {
  r := &PollBot{}
  RegisterRobot("poll", r)
}

func (poll PollBot) Run(p *Payload) (slashCommandImmediateReturn string) {
  go poll.DeferredAction(p)
  return ""
}

type Poll struct {
  Title string `json:"title"`
  Options []string `json:"options"`
}

type CreatedPoll struct {
  Id int `json:"id"`
}

func (poll PollBot) DeferredAction(p *Payload) {
  re := regexp.MustCompile(`^(create)\s*([^:]*):(.*)$`)
  matches := re.FindAllStringSubmatch(p.Text, -1)

  if matches != nil {
    cmd := matches[0][1]
    title := matches[0][2]
    options := strings.Split(matches[0][3], ",")

    for idx, val := range options{
      options[idx] = strings.TrimSpace(val);
    }

    newPoll := Poll{
      Title: title,
      Options: options,
    }

    b, _ := json.Marshal(newPoll)

    if cmd == "create" {

      resp, _ := http.Post(
        "http://strawpoll.me/api/v2/polls",
        "application/json",
        bytes.NewBuffer(b))

      results := CreatedPoll{}

      r, _ := ioutil.ReadAll(resp.Body)

      _ = json.Unmarshal(r, &results)

      response := &IncomingWebhook{
        Channel:     p.ChannelID,
        Username:    "Poll Bot",
        IconEmoji:   ":notebook:",
        Text:        fmt.Sprintf("@group @%s created a poll! http://strawpoll.me/%v", p.UserName, results.Id),
        UnfurlLinks: true,
        Parse:       ParseStyleFull,
      }
      response.Send()
    }
  }
}

func (poll PollBot) Description() (description string) {
  return "Create or vote on polls!\n\tUsage: /poll {something}\n\tExpected Result: @user created a poll!"
}
