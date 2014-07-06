package maize

type Util struct{}

func (u *Util) MergeMap(a *map[interface{}]interface{}, b *map[interface{}]interface{}) {
	if a == nil {
		a = b
	} /*else {
		for key, value := range b {
			a[key] = value
		}
	}*/
}
