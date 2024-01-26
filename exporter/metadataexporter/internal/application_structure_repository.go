package internal

import (
	"context"
	"time"

	"github.com/teanoon/opentelemetry-collector-contrib/exporter/metadataexporter/ent"
	"github.com/teanoon/opentelemetry-collector-contrib/exporter/metadataexporter/ent/applicationstructure"
)

type ApplicationStructureRepository interface {
	FindAll(ctx context.Context) ([]*ent.ApplicationStructure, error)
	SaveAll(ctx context.Context, applicationStructures []*ent.ApplicationStructure) error
	DeleteOutdated(ctx context.Context) error
	Shutdown(ctx context.Context) error
}

type ApplicationStructureRepositoryImpl struct {
	client *DatabaseClient
}

func CreateApplicationStructureRepository(client *DatabaseClient) *ApplicationStructureRepositoryImpl {
	return &ApplicationStructureRepositoryImpl{client}
}

func (repository *ApplicationStructureRepositoryImpl) FindAll(ctx context.Context) ([]*ent.ApplicationStructure, error) {
	return repository.client.delegate.ApplicationStructure.Query().
		All(ctx)
}

func (repository *ApplicationStructureRepositoryImpl) SaveAll(ctx context.Context, applicationStructures []*ent.ApplicationStructure) error {
	if len(applicationStructures) == 0 {
		return nil
	}

	// do upserts in a transaction
	tx, err := repository.client.delegate.Tx(ctx)
	if err != nil {
		return err
	}

	// find all by code
	ids := []string{}
	for _, applicationStructure := range applicationStructures {
		ids = append(ids, applicationStructure.ID)
	}
	existed, err := tx.ApplicationStructure.Query().
		Where(applicationstructure.IDIn(ids...)).
		All(ctx)
	if err != nil {
		return rollback(tx, err)
	}
	existedCodes := map[string]*ent.ApplicationStructure{}
	for _, exist := range existed {
		existedCodes[exist.ID] = exist
	}

	// find updates and inserts
	updates := []*ent.ApplicationStructure{}
	inserts := []*ent.ApplicationStructureCreate{}
	for _, applicationStructure := range applicationStructures {
		if exist, ok := existedCodes[applicationStructure.ID]; ok {
			if exist.ParentCode != applicationStructure.ParentCode || exist.Level != applicationStructure.Level {
				updates = append(updates, applicationStructure)
			}
		} else {
			insert := tx.ApplicationStructure.Create().
				SetID(applicationStructure.ID).
				SetParentCode(applicationStructure.ParentCode).
				SetLevel(applicationStructure.Level).
				SetValidDate(applicationStructure.ValidDate).
				SetCreateTime(applicationStructure.CreateTime).
				SetUpdateTime(applicationStructure.UpdateTime)
			inserts = append(inserts, insert)
		}
	}

	// update all updates one by one
	for _, applicationStructure := range updates {
		_, err := tx.ApplicationStructure.UpdateOneID(applicationStructure.ID).
			SetParentCode(applicationStructure.ParentCode).
			SetLevel(applicationStructure.Level).
			SetValidDate(applicationStructure.ValidDate).
			SetCreateTime(applicationStructure.CreateTime).
			Save(ctx)
		if err != nil {
			return rollback(tx, err)
		}
	}

	// insert all inserts in a batch
	if len(inserts) > 0 {
		_, err := tx.ApplicationStructure.
			CreateBulk(inserts...).
			Save(ctx)
		if err != nil {
			return rollback(tx, err)
		}
	}

	return commit(tx)
}

func (repository *ApplicationStructureRepositoryImpl) DeleteOutdated(ctx context.Context) error {
	_, err := repository.client.delegate.ApplicationStructure.Delete().
		Where(applicationstructure.ValidDateLT(time.Now())).
		Exec(ctx)
	return err
}

func (repository *ApplicationStructureRepositoryImpl) Shutdown(ctx context.Context) error {
	return repository.client.delegate.Close()
}
