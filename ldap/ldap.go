// Package ldap is for general use case functions when dealing with ldap
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

var (
	// DefaultAttr contains the most likely LDAP Attributes to search for
	DefaultAttr = []string{"uid", "mail", "cn"}
)

// ConnectTLS requests a secure binding and requires a Host, Port, User and Pass, this is required if you're
// performing any kind of admin action.
func ConnectTLS(c Config) (*Conn, error) {
	l, err := ldap.DialTLS("tcp", fmt.Sprintf("%s:%d", c.Host, c.Port), &tls.Config{InsecureSkipVerify: true})
	if err != nil {
		return &Conn{}, fmt.Errorf("Dial error: %s", err)
	}

	if err := l.Bind(fmt.Sprintf("cn=%s,%s", c.User, c.BaseDN), c.Pass); err != nil {
		return &Conn{}, fmt.Errorf("Bind error: %s", err)
	}

	return &Conn{l, c}, nil
}

// Connect will attempt to create an unencrypted dial connection to the server. You may need to manually bind
// afterwards in order to
func Connect(c Config) (*Conn, error) {
	// Anonymous LDAP connection
	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", c.Host, c.Port))
	if err != nil {
		return &Conn{}, fmt.Errorf("Dial error: %s", err)
	}
	return &Conn{l, c}, nil
}

// Conn wraps ldap.Conn so that methods can be attached to it for convenience
type Conn struct {
	*ldap.Conn
	Conf Config
}

// GetEntry attempts to retrieve a users LDAP record based on their numerical UID or string UID,
// and returns the desired attributes in the attr arg.
// For example: entry, _ := l.GetEntry(0, "nb5", []string{"cn", "mail"})
// OR         : entry, _ := l.GetEntry(1001, "", []string{"cn", "mail"})
// GetEntry prioritises UIDStr as the preferred seach term because it is considered to be more
// likely to be known.
func (l *Conn) GetEntry(UIDnum uint32, UIDStr string, attr []string) (*Entry, error) {
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
		return &Entry{&ldap.Entry{}, l}, fmt.Errorf("LDAP search error: %s", err.Error())
	} else if len(res.Entries) != 1 {
		uids := []string{}
		for _, e := range res.Entries {
			uids = append(uids, e.GetAttributeValue("uid"))
		}
		return &Entry{&ldap.Entry{}, l}, fmt.Errorf("%d results matched for (%d/%s): %v", len(uids), UIDnum, UIDStr, uids)
	}

	return &Entry{res.Entries[0], l}, nil
}

func selectUIDFilter(n uint32, s string) string {
	if n == 0 && s != "" {
		return fmt.Sprintf("(uid=%s)", s)
	}
	return fmt.Sprintf("(uidNumber=%d)", n)
}

// Entry wraps ldap.Entry so that it can be extended and provides the access to the current ldap Conn
// struct in order to provide convenience in wrapping functions.
type Entry struct {
	*ldap.Entry
	C *Conn
}

// Update updates an ldap entry
func (e *Entry) Update(attr string, data []string, overwrite bool) error {
	// Create modify request object
	m := ldap.NewModifyRequest(e.DN)

	exists := false
	if len(e.GetAttributeValue(attr)) != 0 {
		exists = true
	}

	// Either the attribute isn't set or we wish to overwrite it either way
	if overwrite || !exists {
		if exists {
			// Modify replace contents of attribute
			m.Replace(attr, data)
		} else {
			// Add new data
			m.Add(attr, data)
		}

		if err := e.C.Modify(m); err != nil {
			return fmt.Errorf("Could not modify LDAP record: %s", err)
		}
	} else if !overwrite && exists {
		return fmt.Errorf("data exists and overwrite not set. No change made")
	}

	return nil
}
