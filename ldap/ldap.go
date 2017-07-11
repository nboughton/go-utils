package ldap

import (
	"crypto/tls"
	"fmt"

	"gopkg.in/ldap.v2"
)

// Config defines the minimum fields required to connect to and search LDAP
type Config struct {
	Host   string `json:"host"`
	Port   int    `json:"port"`
	User   string `json:"user"`
	Pass   string `json:"pass"`
	BaseDN string `json:"baseDN"`
}

// SangerLdap wraps ldap.Conn so that methods can be attached to it for convenience
type SangerLdap struct {
	*ldap.Conn
	Conf Config
}

var (
	// DefaultAttr contains the most likely LDAP Attributes to search for
	DefaultAttr = []string{"uid", "mail", "cn"}
)

// ConnectTLS requests a secure binding and requires a Host, Port, User and Pass, this is required if you're
// performing any kind of admin action.
func ConnectTLS(c Config) (*SangerLdap, error) {
	l, err := ldap.DialTLS("tcp", fmt.Sprintf("%s:%d", c.Host, c.Port), &tls.Config{InsecureSkipVerify: true})
	if err != nil {
		return &SangerLdap{}, err
	}

	if err := l.Bind(fmt.Sprintf("cn=%s,%s", c.User, c.BaseDN), c.Pass); err != nil {
		return &SangerLdap{}, err
	}

	return &SangerLdap{l, c}, nil
}

// Connect requests and anonymous binding that requires only a Host and Port. This can only be used to run queries
func Connect(c Config) (*SangerLdap, error) {
	// Anonymous LDAP connection
	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", c.Host, c.Port))
	if err != nil {
		return &SangerLdap{}, err
	}
	return &SangerLdap{l, c}, nil
}

// GetEntry attempts to retrieve a users LDAP record based on their numerical UID or string UID,
// and returns the desired attributes in the attr arg.
// For example: entry, _ := l.GetEntry(0, "nb5", []string{"cn", "mail"})
// OR         : entry, _ := l.GetEntry(1001, "", []string{"cn", "mail"})
// GetEntry prioritises UIDStr as the preferred seach term because it is considered to be more
// likely to be known.
func (l *SangerLdap) GetEntry(UIDnum uint32, UIDStr string, attr []string) (*ldap.Entry, error) {
	// Constuct search request
	searchRequest := ldap.NewSearchRequest(
		fmt.Sprintf("ou=people,%s", l.Conf.BaseDN), // DN
		ldap.ScopeWholeSubtree,                     // Scope
		ldap.NeverDerefAliases,                     // Deref Aliases
		0,     // Size limit
		0,     // Time limit
		false, // Types Only?
		selectUIDFilter(UIDnum, UIDStr), // Filter
		attr, // Attributes
		nil,  // Controls
	)

	// Run search and process results
	res, err := l.Search(searchRequest)
	if err != nil {
		return &ldap.Entry{}, fmt.Errorf("LDAP search error: %s", err.Error())
	} else if len(res.Entries) != 1 {
		return &ldap.Entry{}, fmt.Errorf("Invalid UID (%d/%s), %d results matched", UIDnum, UIDStr, len(res.Entries))
	}

	return res.Entries[0], nil
}

func selectUIDFilter(n uint32, s string) string {
	if n == 0 && s != "" {
		return fmt.Sprintf("(uid=%s)", s)
	}
	return fmt.Sprintf("(uidNumber=%d)", n)
}
