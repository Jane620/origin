package internalversion

import (
	"errors"
	"net/url"

	"k8s.io/client-go/rest"

	buildapi "github.com/openshift/origin/pkg/build/apis/build"
	buildinternalversion "github.com/openshift/origin/pkg/build/generated/internalclientset/typed/build/internalversion"
)

var ErrTriggerIsNotAWebHook = errors.New("the specified trigger is not a webhook")

type WebHookURLInterface interface {
	WebHookURL(name string, trigger *buildapi.BuildTriggerPolicy) (*url.URL, error)
}

func NewWebhookURLClient(c buildinternalversion.BuildInterface, ns string) WebHookURLInterface {
	return &webhooks{client: c.RESTClient(), ns: ns}
}

type webhooks struct {
	client rest.Interface
	ns     string
}

func (c *webhooks) WebHookURL(name string, trigger *buildapi.BuildTriggerPolicy) (*url.URL, error) {
	hooks := c.client.Get().Namespace(c.ns).Resource("buildConfigs").Name(name).SubResource("webhooks")
	switch {
	case trigger.GenericWebHook != nil:
		return hooks.Suffix(trigger.GenericWebHook.Secret, "generic").URL(), nil
	case trigger.GitHubWebHook != nil:
		return hooks.Suffix(trigger.GitHubWebHook.Secret, "github").URL(), nil
	case trigger.GitLabWebHook != nil:
		return hooks.Suffix(trigger.GitLabWebHook.Secret, "gitlab").URL(), nil
	case trigger.BitbucketWebHook != nil:
		return hooks.Suffix(trigger.BitbucketWebHook.Secret, "bitbucket").URL(), nil
	default:
		return nil, ErrTriggerIsNotAWebHook
	}
}
