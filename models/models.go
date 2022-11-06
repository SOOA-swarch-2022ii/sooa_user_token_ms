package models

type User struct {
	ID        string `bson:"id" json:"id"`
	Username  string `bson:"username,omitempty" json:"username,omitempty"`
	Birthdate string `bson:"birthdate,omitempty" json:"birthdate,omitempty"`
	Names     string `bson:"names,omitempty" json:"names,omitempty"`
	Lastnames string `bson:"last_names,omitempty" json:"last_names,omitempty"`
	Password  string `bson:"password,omitempty" json:"password,omitempty"`
	Email     string `bson:"email,omitempty" json:"email,omitempty"`
	Phone     string `bson:"phone,omitempty" json:"phone,omitempty"`
}

type Token struct {
	User     string `bson:"user,omitempty" json:"user,omitempty"`
	Token    string `bson:"token,omitempty" json:"token,omitempty"`
	Creation string `bson:"creation,omitempty" json:"creation,omitempty"`
	Expires  string `bson:"expires,omitempty" json:"expires,omitempty"`
}

type Login struct {
	Username string `bson:"username,omitempty" json:"username,omitempty"`
	Password string `bson:"password,omitempty" json:"password,omitempty"`
}

type LoginResponse struct {
	Token string `bson:"token,omitempty" json:"token,omitempty"`
}
