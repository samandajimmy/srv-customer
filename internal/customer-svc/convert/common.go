package convert

import (
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/dto"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer-svc/model"
)

func ModifierModelToDTO(model model.Modifier) dto.Modifier {
	return dto.Modifier{
		ID:       model.ID,
		Role:     model.Role,
		FullName: model.FullName,
	}
}

func ModifierDTOToModel(dto dto.Modifier) model.Modifier {
	return model.Modifier{
		ID:       dto.ID,
		Role:     dto.Role,
		FullName: dto.FullName,
	}
}

func ItemMetadataModelToResponse(model model.ItemMetadata) dto.ItemMetadataResponse {
	return dto.ItemMetadataResponse{
		UpdatedAt:  model.UpdatedAt.Unix(),
		CreatedAt:  model.CreatedAt.Unix(),
		ModifiedBy: ModifierModelToDTO(*model.ModifiedBy),
		Version:    model.Version,
	}
}
