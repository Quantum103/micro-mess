package database

import (
	"context"
	"fmt"
)

type UserStatus struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	City         string `json:"city,omitempty"`
	FriendStatus string `json:"friend_status"`
}

func (userRepo *UserRepository) ListAllUsers(ctx context.Context, currUserID int, limit, offset int) ([]UserStatus, error) {
	query := `SELECT 
    u.id,
    u.username,
    COALESCE(u.location, '') AS location,

    CASE
        -- уже друзья (в любую сторону)
        WHEN EXISTS (
            SELECT 1 FROM friends f
            WHERE (
                (f.user_id = ? AND f.friend_id = u.id)
                OR
                (f.user_id = u.id AND f.friend_id = ?)
            )
            AND f.status = 'accepted'
        ) THEN 'accepted'

        -- я отправил заявку
        WHEN EXISTS (
            SELECT 1 FROM friends f
            WHERE f.user_id = ? 
              AND f.friend_id = u.id 
              AND f.status = 'pending'
        ) THEN 'pending'

        -- мне отправили заявку
        WHEN EXISTS (
            SELECT 1 FROM friends f
            WHERE f.user_id = u.id 
              AND f.friend_id = ? 
              AND f.status = 'pending'
        ) THEN 'incoming'

        ELSE 'none'
    END AS friend_status

FROM users u
WHERE u.id != ?
ORDER BY u.username
LIMIT ? OFFSET ?`
	rows, err := userRepo.db.QueryContext(ctx, query,
		currUserID,
		currUserID,
		currUserID,
		currUserID,
		currUserID,
		limit,
		offset,
	)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения списка пользователей %w", err)
	}
	defer rows.Close()

	var users []UserStatus
	for rows.Next() {
		var u UserStatus
		err := rows.Scan(&u.ID, &u.Username, &u.City, &u.FriendStatus)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования списка пользователей %w", err)
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

func (userRepo *UserRepository) GetFrineds(ctx context.Context, userID int, limit, offset int) ([]UserStatus, error) {
	query := `SELECT 
    u.id,
    u.username,
    COALESCE(u.location, '') AS location,

    'accepted' AS friend_status

FROM users u
INNER JOIN friends f
    ON (
        (f.user_id = ? AND f.friend_id = u.id)
        OR
        (f.user_id = u.id AND f.friend_id = ?)
    )
WHERE f.status = 'accepted'
  AND u.id != ?

ORDER BY u.username
LIMIT ? OFFSET ?`

	rows, err := userRepo.db.QueryContext(ctx, query,
		userID,
		userID,
		userID,
		limit,
		offset,
	)
	if err != nil {
		return nil, fmt.Errorf("ошибка сканирования списка друзей %w", err)
	}
	defer rows.Close()

	var friends []UserStatus
	for rows.Next() {
		var f UserStatus
		err := rows.Scan(&f.ID, &f.Username, &f.City, &f.FriendStatus)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования списка друзей %w", err)
		}
		friends = append(friends, f)
	}
	return friends, rows.Err()
}
func (userRepo *UserRepository) AddFriend(ctx context.Context, userID, friendID int) error {
	if userID == friendID {
		return fmt.Errorf("нельзя добавить себя в друзья")
	}

	var status string
	var fromUser int

	err := userRepo.db.QueryRowContext(ctx, `
		SELECT user_id, status 
		FROM friends 
		WHERE (user_id = ? AND friend_id = ?) 
		   OR (user_id = ? AND friend_id = ?)
	`,
		userID, friendID,
		friendID, userID,
	).Scan(&fromUser, &status)

	if err == nil {
		// если есть входящая заявка → принимаем
		if status == "pending" && fromUser == friendID {
			_, err := userRepo.db.ExecContext(ctx, `
				UPDATE friends 
				SET status = 'accepted'
				WHERE user_id = ? AND friend_id = ?
			`, friendID, userID)
			return err
		}

		return fmt.Errorf("уже есть связь")
	}

	// ✅ ВАЖНО: создаём новую заявку
	_, err = userRepo.db.ExecContext(ctx, `
		INSERT INTO friends (user_id, friend_id, status)
		VALUES (?, ?, 'pending')
	`, userID, friendID)

	if err != nil {
		return fmt.Errorf("ошибка создания заявки: %w", err)
	}

	return nil
}

func (userRepo *UserRepository) AcceptFriend(ctx context.Context, userID, friendID int) error {
	res, err := userRepo.db.ExecContext(ctx, `
        UPDATE friends 
        SET status = 'accepted' 
        WHERE user_id = ? AND friend_id = ? AND status = 'pending'
    `, friendID, userID)

	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("заявка не найдена")
	}

	return nil
}
