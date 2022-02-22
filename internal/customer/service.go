package customer

import (
	"context"
	"github.com/nbs-go/nlogger"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/constant"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer/model"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/dto"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nclient"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/ncore"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nredis"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/ns3"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nval"
)

type Service struct {
	config       *Config
	ctx          context.Context
	repo         *RepositoryContext
	repoExternal *RepositoryContext
	log          nlogger.Logger
	responses    *ncore.ResponseMap
	redis        *nredis.Redis
	minio        *ns3.Minio
	client       *nclient.Nclient
	pdsClient    *nclient.Nclient
}

func (h Handler) NewService(ctx context.Context) *Service {
	svc := Service{
		config:       h.Config,
		client:       h.Client,
		pdsClient:    h.PdsAPIClient,
		redis:        h.Redis,
		minio:        h.Minio,
		responses:    h.Responses,
		repo:         h.Repo.WithContext(ctx),
		repoExternal: h.RepoExternal.WithContext(ctx),
		ctx:          ctx,
		log:          nlogger.NewChild(nlogger.WithNamespace("service"), nlogger.Context(ctx)),
	}

	return &svc
}

func (s *Service) Close() {
	// Close database connection to free pool
	err := s.repo.conn.Close()
	if err != nil {
		s.log.Error("Failed to close connection", nlogger.Error(err))
	}

	// Close database external connection to free pool
	err = s.repoExternal.conn.Close()
	if err != nil {
		s.log.Error("Failed to close external connection", nlogger.Error(err))
	}
}

func (s *Service) composeProfileResponse(customer *model.Customer, address *model.Address, financial *model.FinancialData,
	verification *model.Verification, gs interface{}) dto.ProfileResponse {
	avatarURL := s.AssetGetPublicURL(constant.AssetAvatarProfile, customer.Photos.FileName)
	ktpURL := s.AssetGetPublicURL(constant.AssetKTP, customer.Profile.IdentityPhotoFile)
	npwpURL := s.AssetGetPublicURL(constant.AssetNPWP, customer.Profile.NPWPPhotoFile)
	sidURL := s.AssetGetPublicURL(constant.AssetNPWP, customer.Profile.SidPhotoFile)

	return dto.ProfileResponse{
		CustomerVO: dto.CustomerVO{
			ID:                        customer.UserRefID.String,
			Cif:                       customer.Cif,
			Nama:                      customer.FullName,
			NamaIbu:                   customer.Profile.MaidenName,
			NoKTP:                     customer.IdentityNumber,
			Email:                     customer.Email,
			JenisKelamin:              customer.Profile.Gender,
			TempatLahir:               customer.Profile.PlaceOfBirth,
			TglLahir:                  customer.Profile.DateOfBirth,
			ReferralCode:              customer.ReferralCode,
			NoHP:                      customer.Phone,
			Kewarganegaraan:           customer.Profile.Nationality,
			NoIdentitas:               customer.IdentityNumber,
			TglExpiredIdentitas:       customer.Profile.IdentityExpiredAt,
			StatusKawin:               customer.Profile.MarriageStatus,
			NoNPWP:                    customer.Profile.NPWPNumber,
			NoSid:                     customer.Sid,
			IsKYC:                     nval.ParseStringFallback(verification.KycVerifiedStatus, ""),
			JenisIdentitas:            nval.ParseStringFallback(customer.IdentityType, ""),
			FotoNPWP:                  npwpURL,
			FotoSid:                   sidURL,
			Avatar:                    avatarURL,
			FotoKTP:                   ktpURL,
			Alamat:                    address.Line.String,
			IDProvinsi:                address.ProvinceID.Int64,
			IDKabupaten:               address.CityID.Int64,
			IDKecamatan:               address.DistrictID.Int64,
			IDKelurahan:               address.SubDistrictID.Int64,
			Kelurahan:                 address.SubDistrictName.String,
			Provinsi:                  address.ProvinceName.String,
			Kabupaten:                 address.CityName.String,
			Kecamatan:                 address.DistrictName.String,
			KodePos:                   address.PostalCode.String,
			Norek:                     financial.AccountNumber,
			GoldCardApplicationNumber: financial.GoldCardApplicationNumber,
			GoldCardAccountNumber:     financial.GoldCardAccountNumber,
			Saldo:                     nval.ParseStringFallback(financial.Balance, "0"),
			IsOpenTe:                  nval.ParseStringFallback(financial.GoldSavingStatus, "0"),
			IsEmailVerified:           nval.ParseStringFallback(verification.EmailVerifiedStatus, "0"),
			IsDukcapilVerified:        nval.ParseStringFallback(verification.DukcapilVerifiedStatus, "0"),
			AktifasiTransFinansial:    nval.ParseStringFallback(verification.FinancialTransactionStatus, "0"),
			KodeCabang:                "", // TODO Branch Code
			TabunganEmas:              gs,
		},
	}
}
