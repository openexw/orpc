package testdata

type Profile struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
	Age  uint8  `json:"age"`
	Sex  uint8  `json:"sex"`
}

var ProfileSingle = &Profile{
	Id:   1,
	Name: "Jack",
	Age:  10,
	Sex:  1,
}
var JsonProfileStr = `{"id":1,"name":"Jack","age":10,"sex":1}`
var JsonZeroStr = `{"id":0,"name":"","age":0,"sex":0}`

type AddProfileArgs struct {
	Name string
	Age  uint8
	Sex  uint8
}

func (p *Profile) AddProfile(args *AddProfileArgs, reply *Profile) error {
	*reply = Profile{
		Id:   1,
		Name: args.Name,
		Age:  args.Age,
		Sex:  args.Sex,
	}
	return nil
}
