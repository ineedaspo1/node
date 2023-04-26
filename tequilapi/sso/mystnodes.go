package sso

import (
	"encoding/base64"
	"github.com/mysteriumnetwork/node/config"
	"github.com/mysteriumnetwork/node/eventbus"
	"github.com/mysteriumnetwork/node/identity"
	"github.com/mysteriumnetwork/node/tequilapi/pkce"
	"net/url"
)

type Mystnodes struct {
	baseUrl              string
	ssoPath              string
	signer               identity.SignerFactory
	lastUnlockedIdentity identity.Identity
}

func NewMystnodes(signer identity.SignerFactory) *Mystnodes {
	return &Mystnodes{
		baseUrl: config.GetString(config.FlagMMNAddress),
		ssoPath: "/login-sso",
		signer:  signer,
	}
}

func (m *Mystnodes) Subscribe(eventBus eventbus.EventBus) error {
	if err := eventBus.SubscribeAsync(identity.AppTopicIdentityUnlock, m.onIdentityUnlocked); err != nil {
		return err
	}
	return nil
}

func (m *Mystnodes) onIdentityUnlocked(ev identity.AppEventIdentityUnlock) {
	m.lastUnlockedIdentity = ev.ID
}

func (m *Mystnodes) message(info pkce.Info) MystnodesMessage {
	return MystnodesMessage{
		CodeChallenge: info.CodeChallenge,
		Identity:      m.lastUnlockedIdentity.Address,
		RedirectURL:   "placeholder",
	}
}

func (m *Mystnodes) sign(msg []byte) (identity.Signature, error) {
	return m.signer(m.lastUnlockedIdentity).Sign(msg)
}

func (m *Mystnodes) SSOLink() (*url.URL, error) {
	u, err := url.Parse(m.baseUrl)
	if err != nil {
		return &url.URL{}, err
	}

	u = u.JoinPath(m.ssoPath)

	info, err := pkce.New(128)
	if err != nil {
		return nil, err
	}

	messageJson, err := m.message(info).json()
	if err != nil {
		return &url.URL{}, err
	}

	signature, err := m.sign(messageJson)
	if err != nil {
		return &url.URL{}, err
	}

	q := u.Query()
	q.Set("message", base64.RawURLEncoding.EncodeToString(messageJson))
	q.Set("signature", base64.RawURLEncoding.EncodeToString(signature.Bytes()))
	u.RawQuery = q.Encode()

	return u, nil
}
