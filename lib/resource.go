package lib

import "strings"

func GetClient(clientID string) (cli *ConnClientConf) {
    /*
	for _, v := range ConfConnCientMap{
		if v.ID == clientID {
			cli = v
		}
	}*/
	for _,clientConf := range ConfConnCientMap.List {
		if clientConf.ID == clientID{
			cli = clientConf
		}
	}
	return
}

func ScopeJoin(scope []Scope) string {
	var s []string
	for _, sc := range scope {
		s = append(s, sc.ID)
	}
	return strings.Join(s,",")
}

func ScopeFilter(clientID string, scope string) (s []Scope) {
	cli := GetClient(clientID)
	sl := strings.Split(scope, ",")
	for _, str := range sl {
		for _, sc := range cli.Scope {
			if str == sc.ID {
				s = append(s, sc)
			}
		}
	}

	return
}