package sqlstore

import (
	"context"
	"time"

	"o.o/api/fabo/fbmessagetemplate"
	"o.o/backend/com/fabo/main/fbmessagetemplate/convert"
	"o.o/backend/com/fabo/main/fbmessagetemplate/model"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sqlstore"
	"o.o/capi/dot"
)

type FbMessageTemplateStore struct {
	ft FbMessageTemplateFilters

	query cmsql.QueryFactory
	preds []interface{}
}

type FbMessageTemplateStoreFactory func(context.Context) *FbMessageTemplateStore

func NewFbMessageTemplateStoreFactory(db *cmsql.Database) FbMessageTemplateStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *FbMessageTemplateStore {
		return &FbMessageTemplateStore{
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

func (s *FbMessageTemplateStore) ID(id dot.ID) *FbMessageTemplateStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *FbMessageTemplateStore) ShopID(id dot.ID) *FbMessageTemplateStore {
	s.preds = append(s.preds, s.ft.ByShopID(id))
	return s
}

func (s *FbMessageTemplateStore) CreateFbMessageTemplateDB(template *fbmessagetemplate.FbMessageTemplate) (*model.FbMessageTemplate, error) {
	sqlstore.MustNoPreds(s.preds)
	templateDB := &model.FbMessageTemplate{}
	templateDB = convert.Convert_fbmessagetemplate_FbMessageTemplate_fbmessagetemplatemodel_FbMessageTemplate(template, templateDB)
	templateDB.ID = cm.NewID()
	templateDB.CreatedAt = time.Now()
	templateDB.UpdatedAt = time.Now()
	_, err := s.query().Insert(templateDB)
	return templateDB, err
}

func (s *FbMessageTemplateStore) CreateFbMessageTemplate(createTemplate *fbmessagetemplate.CreateFbMessageTemplate) (*fbmessagetemplate.FbMessageTemplate, error) {
	template := &fbmessagetemplate.FbMessageTemplate{}
	template = convert.Apply_fbmessagetemplate_CreateFbMessageTemplate_fbmessagetemplate_FbMessageTemplate(createTemplate, template)
	templateDB, err := s.CreateFbMessageTemplateDB(template)
	if err != nil {
		return nil, err
	}

	return convert.Convert_fbmessagetemplatemodel_FbMessageTemplate_fbmessagetemplate_FbMessageTemplate(templateDB, template), nil
}

func (s *FbMessageTemplateStore) UpdateMessageTemplate(args *fbmessagetemplate.UpdateFbMessageTemplate) error {
	template := &model.FbMessageTemplate{}
	if args.Template.Valid {
		template.Template = args.Template.String
	}

	if args.ShortCode.Valid {
		template.ShortCode = args.ShortCode.String
	}

	s.ID(args.ID).ShopID(args.ShopID)
	err := s.query().Where(s.preds).ShouldUpdate(template)
	return err
}

func (s *FbMessageTemplateStore) DeleteMessageTemplate() error {
	return s.query().Where(s.preds).ShouldDelete(&model.FbMessageTemplate{})
}

func (s *FbMessageTemplateStore) GetFbMessageTemplatesDB() (model.FbMessageTemplates, error) {
	var result model.FbMessageTemplates
	err := s.query().Where(s.preds).OrderBy("created_at desc").Find(&result)
	return result, err
}

func (s *FbMessageTemplateStore) GetFbMessageTemplates() ([]*fbmessagetemplate.FbMessageTemplate, error) {
	templateDBs, err := s.GetFbMessageTemplatesDB()
	if err != nil {
		return nil, err
	}

	return convert.Convert_fbmessagetemplatemodel_FbMessageTemplates_fbmessagetemplate_FbMessageTemplates(templateDBs), nil
}
