package main

type Players []*Player

func (pp Players) Playing() (ret Players) {
	for _, p := range pp {
		if !p.folded && p.ChipAmount > 0 {
			ret = append(ret, p)
		}
	}
	return
}
