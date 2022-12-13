package box_type

type CompatibilityList []BoxType

func ExtendFromSlice(cl CompatibilityList, buf []byte) CompatibilityList {
	if len(cl) <= 0 {
		cl = make([]BoxType, 0)
	}
	a := make(BoxType, len(buf))
	a.From(buf)
	cl = append(cl, a)
	return cl
}
