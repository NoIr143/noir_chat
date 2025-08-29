package entities

import "time"

type FriendRequestStatus int

const (
	Waiting FriendRequestStatus = iota
	Accepted
)

func (s FriendRequestStatus) String() string {
	switch s {
	case Waiting:
		return "waiting"
	case Accepted:
		return "accepted"
	default:
		return "unknown"
	}
}

type FriendRequest struct {
	ID         int                 `json:"id"`
	SenderId   int                 `json:"sender_id"`
	ReceiverId int                 `json:"receiver_id"`
	Status     FriendRequestStatus `json:"status"`
	CreatedAt  time.Time           `json:"created_at"`
	UpdatedAt  time.Time           `json:"updated_at"`
}
