package service

type newMatchRequest struct{
	GridSize int `json:"gridsize"`
	PlayerWhite string `json:"playerWhite"`
	PlayerBlack string `json:"playerBlack"`
}

func(request newMatchRequest) isValid() (valid bool){
	valid = true
	if request.GridSize !=19 && request.GridSize !=13 && request.GridSize !=9{
		valid = false
	}
	if request.PlayerBlack =="" || request.PlayerWhite == ""{
		valid = false
	}
	return valid
}