package etop

type (
	Status3 int
	Status4 int
	Status5 int
)

const (
	S3Negative Status3 = -1
	S3Zero     Status3 = 0
	S3Positive Status3 = 1

	S4Negative Status4 = -1
	S4Zero     Status4 = 0
	S4SuperPos Status4 = 2
	S4Positive Status4 = 1

	S5NegSuper Status5 = -2
	S5Negative Status5 = -1
	S5Zero     Status5 = 0
	S5Positive Status5 = 1
	S5SuperPos Status5 = 2
)

func Status3FromInt(s int) Status3 {
	return Status3(s)
}

func Status4FromInt(s int) Status4 {
	return Status4(s)
}

func Status5FromInt(s int) Status5 {
	return Status5(s)
}
