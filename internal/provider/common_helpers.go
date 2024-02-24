package provider

import (
	"context"

	"aembit.io/aembit"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func newTagsModel(ctx context.Context, tags []aembit.TagDTO) types.Map {
	respMap := make(map[string]string)

	if len(tags) > 0 {
		tflog.Info(ctx, "newTagsModel: tags found.")
		for _, tagEntry := range tags {
			respMap[tagEntry.Key] = tagEntry.Value
		}
		tagsMap, _ := types.MapValueFrom(ctx, types.StringType, respMap)
		return tagsMap
	}

	return types.MapNull(types.StringType)
}
