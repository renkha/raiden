package roles

import (
	"github.com/sev-2/raiden"
	"github.com/sev-2/raiden/pkg/resource/migrator"
	"github.com/sev-2/raiden/pkg/state"
	"github.com/sev-2/raiden/pkg/supabase"
	"github.com/sev-2/raiden/pkg/supabase/objects"
)

type MigrateItem = migrator.MigrateItem[objects.Role, objects.UpdateRoleParam]
type MigrateActionFunc = migrator.MigrateActionFunc[objects.Role, objects.UpdateRoleParam]

var ActionFunc = MigrateActionFunc{
	CreateFunc: supabase.CreateRole, UpdateFunc: supabase.UpdateRole, DeleteFunc: supabase.DeleteRole,
}

func BuildMigrateData(extractedLocalData state.ExtractRoleResult, supabaseData []objects.Role) (migrateData []MigrateItem, err error) {
	Logger.Info("start build migrate role data")
	// compare and bind existing table to migrate data
	mapSpRole := make(map[int]bool)
	for i := range supabaseData {
		t := supabaseData[i]
		mapSpRole[t.ID] = true
	}

	// filter existing table need compare or move to create new
	Logger.Debug("filter extracted data for update new local role data")
	var compareRoles []objects.Role
	for i := range extractedLocalData.Existing {
		et := extractedLocalData.Existing[i]
		if _, isExist := mapSpRole[et.ID]; isExist {
			compareRoles = append(compareRoles, et)
		} else {
			extractedLocalData.New = append(extractedLocalData.New, et)
		}
	}

	if rs, err := BuildMigrateItem(supabaseData, compareRoles); err != nil {
		return migrateData, err
	} else {
		migrateData = append(migrateData, rs...)
	}

	// bind new table to migrated data
	Logger.Debug("filter new role data")
	if len(extractedLocalData.New) > 0 {
		for i := range extractedLocalData.New {
			t := extractedLocalData.New[i]
			migrateData = append(migrateData, MigrateItem{
				Type:    migrator.MigrateTypeCreate,
				NewData: t,
			})
		}
	}

	Logger.Debug("filter delete role data")
	if len(extractedLocalData.Delete) > 0 {
		for i := range extractedLocalData.Delete {
			t := extractedLocalData.Delete[i]
			isExist := false
			for i := range supabaseData {
				tt := supabaseData[i]
				if tt.Name == t.Name {
					isExist = true
					break
				}
			}

			if isExist {
				migrateData = append(migrateData, MigrateItem{
					Type:    migrator.MigrateTypeDelete,
					OldData: t,
				})
			}
		}
	}
	Logger.Info("finish build migrate role data")
	return
}

func BuildMigrateItem(supabaseData []objects.Role, localData []objects.Role) (migrateData []MigrateItem, err error) {
	Logger.Info("compare supabase and local resource for existing role data")
	result, e := CompareList(localData, supabaseData)
	if e != nil {
		err = e
		return
	}

	for i := range result {
		r := result[i]

		migrateType := migrator.MigrateTypeIgnore
		if r.IsConflict {
			migrateType = migrator.MigrateTypeUpdate
		}

		r.DiffItems.OldData = r.TargetResource
		migrateData = append(migrateData, MigrateItem{
			Type:           migrateType,
			NewData:        r.SourceResource,
			OldData:        r.TargetResource,
			MigrationItems: r.DiffItems,
		})
	}

	return
}

func Migrate(config *raiden.Config, roles []MigrateItem, stateChan chan any, actions MigrateActionFunc) []error {
	return migrator.MigrateResource(config, roles, stateChan, actions, migrator.DefaultMigrator)
}
