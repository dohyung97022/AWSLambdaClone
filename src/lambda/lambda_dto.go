package lambda

type Lambda struct {
	Id       string `bson:"_id,omitempty" json:"_id,omitempty"`
	Title    string `bson:"title,omitempty" json:"title"`
	Runtime  string `bson:"runtime,omitempty" json:"runtime"`
	Version  string `bson:"version,omitempty" json:"version"`
	Endpoint string `bson:"endpoint,omitempty" json:"endpoint,omitempty"`
	Code     string `bson:"code,omitempty" json:"code"`
}
