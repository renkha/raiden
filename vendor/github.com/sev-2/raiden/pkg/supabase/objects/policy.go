package objects

import "strings"

type PolicyCommand string

type Policy struct {
	ID         int           `json:"id"`
	Schema     string        `json:"schema"`
	Table      string        `json:"table"`
	TableID    int           `json:"table_id"`
	Name       string        `json:"name"`
	Action     string        `json:"action"`
	Roles      []string      `json:"roles"`
	Command    PolicyCommand `json:"command"`
	Definition string        `json:"definition"`
	Check      *string       `json:"check"`
}

const (
	PolicyCommandSelect PolicyCommand = "SELECT"
	PolicyCommandInsert PolicyCommand = "INSERT"
	PolicyCommandUpdate PolicyCommand = "UPDATE"
	PolicyCommandDelete PolicyCommand = "DELETE"
)

type Policies []Policy

func (p *Policies) FilterByTable(table string) Policies {
	var filteredData Policies
	if p == nil {
		return filteredData
	}

	for _, v := range *p {
		if v.Table == table {
			filteredData = append(filteredData, v)
		}
	}

	return filteredData
}

func (p *Policies) FilterByBucket(bucket Bucket) Policies {
	var filteredData Policies
	if p == nil {
		return filteredData
	}

	for _, v := range *p {
		if v.Schema != "storage" {
			continue
		}

		if strings.Contains(v.Definition, bucket.Name) || (v.Check != nil && strings.Contains(*v.Check, bucket.Name)) {
			filteredData = append(filteredData, v)
		}
	}
	return filteredData
}

type UpdatePolicyType string

const (
	UpdatePolicyName       UpdatePolicyType = "name"
	UpdatePolicyDefinition UpdatePolicyType = "definition"
	UpdatePolicyCheck      UpdatePolicyType = "check"
	UpdatePolicyRoles      UpdatePolicyType = "roles"
)

type UpdatePolicyParam struct {
	Name        string
	ChangeItems []UpdatePolicyType
}
