package repositories

import (
	"database/sql"

	"github.com/noir143/noir_chat/src/database/entities"
)

type FriendRequestRepository struct {
	db *sql.DB
}

func FriendRequestRepositoryConstructor(db *sql.DB) *FriendRequestRepository {
	return &FriendRequestRepository{db: db}
}

func (r *FriendRequestRepository) Create(friendRequest entities.FriendRequest) (entities.FriendRequest, error) {
	query := `
		INSERT INTO friend_requests (sender_id, receiver_id, status)
		VALUES ($1, $2, $3)
		RETURNING id, sender_id, receiver_id, status, created_at, updated_at
	`

	var createdFriendRequest entities.FriendRequest
	err := r.db.QueryRow(query, friendRequest.SenderId, friendRequest.ReceiverId, friendRequest.Status).Scan(
		&createdFriendRequest.ID,
		&createdFriendRequest.SenderId,
		&createdFriendRequest.ReceiverId,
		&createdFriendRequest.Status,
		&createdFriendRequest.CreatedAt,
		&createdFriendRequest.UpdatedAt,
	)

	if err != nil {
		return entities.FriendRequest{}, err
	}

	return createdFriendRequest, nil
}

func (r *FriendRequestRepository) Approve(id int) (entities.FriendRequestStatus, error) {
	query := `
		UPDATE friend_requests
		SET status = $1
		WHERE id = $2
	`

	_, err := r.db.Exec(query, entities.Accepted.String(), id)

	if err != nil {
		return entities.Accepted, err
	}

	return entities.Accepted, nil
}

func (r *FriendRequestRepository) GetFriendIds(id int) ([]int, error) {
	query := `
		SELECT sender_id, receiver_id
		FROM friend_request
		WHERE (receiver_id = $1 OR sender_id = $1) AND status = $2
	`

	rows, err := r.db.Query(query, id, entities.Accepted)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var friendIds []int
	for rows.Next() {
		var senderId, receiverId int
		if err := rows.Scan(&senderId, &receiverId); err != nil {
			return nil, err
		}

		if senderId != id {
			friendIds = append(friendIds, senderId)
		}
		if receiverId != id {
			friendIds = append(friendIds, receiverId)
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return friendIds, nil
}
