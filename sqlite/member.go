package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ddritzenhoff/dindin"
	"github.com/ddritzenhoff/dindin/sqlite/gen"
)

// Ensure service implements interface.
var _ dindin.MemberService = (*MemberService)(nil)

// MemberService represents a service for managing members.
type MemberService struct {
	query *gen.Queries
	db    *sql.DB
}

// NewMemberService returns a new instance of MemberService.
func NewMemberService(query *gen.Queries, db *sql.DB) *MemberService {
	return &MemberService{query, db}
}

// Retrieves a member by ID
// Returns ErrNotFound if meal does not exist.
func (ms *MemberService) FindMemberByID(id int64) (*dindin.Member, error) {
	m, err := ms.query.FindMemberByID(context.Background(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, dindin.ErrNotFound
		} else {
			return nil, fmt.Errorf("FindMemberByID: %w", err)
		}
	}
	isLeader := m.Leader == 1
	return &dindin.Member{
		ID:          m.ID,
		SlackUID:    m.SlackUid,
		FullName:    m.FullName,
		MealsEaten:  m.MealsEaten,
		MealsCooked: m.MealsCooked,
		Leader:      isLeader,
	}, nil
}

// Retrieves a member by SlackID
// Returns ErrNotFound if meal does not exist.
func (ms *MemberService) FindMemberBySlackUID(slackUID string) (*dindin.Member, error) {
	m, err := ms.query.FindMemberBySlackUID(context.Background(), slackUID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, dindin.ErrNotFound
		} else {
			return nil, fmt.Errorf("FindMemberBySlackUID: %w", err)
		}
	}
	isLeader := m.Leader == 1
	return &dindin.Member{
		ID:          m.ID,
		SlackUID:    m.SlackUid,
		FullName:    m.FullName,
		MealsEaten:  m.MealsEaten,
		MealsCooked: m.MealsCooked,
		Leader:      isLeader,
	}, nil
}

// Retrieves a list of members.
func (ms *MemberService) ListMembers() ([]*dindin.Member, error) {
	mems, err := ms.query.ListMembers(context.Background())
	if err != nil {
		return nil, fmt.Errorf("ListMembers: %w", err)
	}
	var members []*dindin.Member
	for ii := 0; ii < len(mems); ii++ {
		m := &mems[ii]
		isLeader := m.Leader == 1
		members = append(members, &dindin.Member{
			ID:          m.ID,
			SlackUID:    m.SlackUid,
			FullName:    m.FullName,
			MealsEaten:  m.MealsEaten,
			MealsCooked: m.MealsCooked,
			Leader:      isLeader,
		})
	}
	return members, nil
}

// Creates a new member.
func (ms *MemberService) CreateMember(m *dindin.Member) error {
	var isLeader int64
	if m.Leader {
		isLeader = 1
	} else {
		isLeader = 0
	}
	params := gen.CreateMemberParams{
		SlackUid: m.SlackUID,
		FullName: m.FullName,
		Leader:   isLeader,
	}
	_, err := ms.query.CreateMember(context.Background(), params)
	if err != nil {
		return fmt.Errorf("CreateMember: %w", err)
	}
	return nil
}

// Updates a member object.
func (ms *MemberService) UpdateMember(id int64, upd dindin.MemberUpdate) error {
	tx, err := ms.db.Begin()
	if err != nil {
		return fmt.Errorf("UpdateMember db.Begin: %w", err)
	}
	defer tx.Rollback()
	qtx := ms.query.WithTx(tx)
	if upd.Leader != nil {
		var isLeader int64
		if *upd.Leader {
			isLeader = 1
		} else {
			isLeader = 0
		}
		params := gen.UpdateMemberLeaderStatusParams{ID: id, Leader: isLeader}
		err := qtx.UpdateMemberLeaderStatus(context.Background(), params)
		if err != nil {
			return fmt.Errorf("UpdateMember UpdateMemberLeaderStatus: %w", err)
		}
	}
	if upd.MealsCooked != nil {
		params := gen.UpdateMemberMealsCookedParams{ID: id, MealsCooked: *upd.MealsCooked}
		err := qtx.UpdateMemberMealsCooked(context.Background(), params)
		if err != nil {
			return fmt.Errorf("UpdateMember UpdateMemberMealsCooked: %w", err)
		}
	}
	if upd.MealsEaten != nil {
		params := gen.UpdateMemberMealsEatenParams{ID: id, MealsEaten: *upd.MealsEaten}
		err := qtx.UpdateMemberMealsEaten(context.Background(), params)
		if err != nil {
			return fmt.Errorf("UpdateMember UpdateMemberMealsEaten: %w", err)
		}
	}
	return tx.Commit()
}

// Permanently deletes a member.
func (ms *MemberService) DeleteMember(id int64) error {
	err := ms.query.DeleteMember(context.Background(), id)
	if err != nil {
		return fmt.Errorf("DeleteMember: %w", err)
	}
	return nil
}
