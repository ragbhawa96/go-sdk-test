package model


type InitRequestModel struct{
	Name  string `json:"name" bson:"name"`
}



type InitResponseModel struct{
	Name  string `json:"NAME" bson:"NAME"`
}




