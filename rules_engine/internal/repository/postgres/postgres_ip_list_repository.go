package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"net"

	"rules-engine/internal/entity"

	"rules-engine/internal/repository"

	"github.com/google/uuid"
)

type PostgresIPListRepository struct {
	db *sql.DB
}

func NewPostgresIPListRepository(db *sql.DB) repository.IPListRepository {
	return &PostgresIPListRepository{db: db}
}

func (r *PostgresIPListRepository) GetIPLists() ([]entity.IPList, error) {
	rows, err := r.db.Query("SELECT id, ip, list_type, creator_id, created_at FROM ip_lists")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lists []entity.IPList
	for rows.Next() {
		var res entity.IPList
		var cidrStr string
		if err := rows.Scan(&res.ID, &cidrStr, &res.ListType, &res.CreatorID, &res.CreatedAt); err != nil {
			return nil, err
		}

		ip, ipNet, err := net.ParseCIDR(cidrStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse CIDR: %w", err)
		}
		ipNet.IP = ip
		res.IP = *ipNet

		lists = append(lists, res)
	}
	return lists, nil
}

func (r *PostgresIPListRepository) CreateIPList(ipList *entity.IPList) error {
	ipList.ID = uuid.New().String()
	_, err := r.db.Exec("INSERT INTO ip_lists (id, ip, list_type, creator_id) VALUES ($1, $2, $3, $4)",
		ipList.ID, ipList.IP.String(), ipList.ListType, ipList.CreatorID)
	return err
}

func (r *PostgresIPListRepository) UpdateIPList(ipList *entity.IPList) error {
	_, err := r.db.Exec("UPDATE ip_lists SET ip=$1, list_type=$2 WHERE id=$3",
		ipList.IP.String(), ipList.ListType, ipList.ID)
	return err
}

func (r *PostgresIPListRepository) GetIPList(id string) (*entity.IPList, error) {
	query := `SELECT id, ip, list_type, creator_id, created_at FROM ip_lists WHERE id = $1`

	list := &entity.IPList{}
	var cidrStr string
	err := r.db.QueryRow(query, id).Scan(
		&list.ID,
		&cidrStr,
		&list.ListType,
		&list.CreatorID,
		&list.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get IP list: %w", err)
	}

	ip, ipNet, err := net.ParseCIDR(cidrStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse CIDR: %w", err)
	}
	ipNet.IP = ip
	list.IP = *ipNet

	return list, nil
}
