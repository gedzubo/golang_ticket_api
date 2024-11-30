package models

type TicketOptionInput struct {
	Name       string `json:"name" binding:"required"`
	Desc       string `json:"desc" binding:"required"`
	Allocation uint64 `json:"allocation" binding:"required"`
}
