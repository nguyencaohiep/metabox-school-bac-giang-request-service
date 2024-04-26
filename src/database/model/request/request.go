package model_account

import (
	"time"

	graph_model "metabox-school-bac-giang-request-service/src/graph/generated/model"
)

type Request struct {
	ID string `json:"id" bson:"_id"`

	Username string `json:"username" bson:"username"`
	Phone    string `json:"phone" bson:"phone"`
	Note     string `json:"note" bson:"note"`

	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
	Deleted   bool      `json:"deleted" bson:"deleted"`
}

func (r *Request) ConvertToModelGraph() *graph_model.Request {
	data := graph_model.Request{
		ID: r.ID,

		Note:  r.Note,
		Phone: r.Phone,

		Account: &graph_model.Account{
			UserName: r.Username,
		},

		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}

	return &data
}
