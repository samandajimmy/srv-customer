-- Sample Migration Customer
CREATE TABLE public."Customer"
(
    "id"             bigint       NOT NULL,
    "xid"            varchar(64)  NOT NULL,
    "metadata"       JSON         NULL,
    "createdAt"      timestamp without time zone NOT NULL,
    "updatedAt"      timestamp without time zone NOT NULL,
    "modifiedBy"     JSON         NOT NULL,
    "version"        bigint       NOT NULL DEFAULT 1,
    "fullName"       varchar(255) NOT NULL,
    "phone"          varchar(16)  NOT NULL,
    "email"          varchar(16)  NOT NULL,
    "identityType"   SMALLINT     NULL,
    "identityNumber" varchar(64)  NULL,
    "userRefId"      bigint       NULL,
    "photos"         JSON         NULL,
    "profile"        JSON         NULL,
    "cif"            varchar(16)  NULL,
    "sid"            varchar(16)  NULL,
    "referralCode"   varchar(64)  NULL,
    "status"         smallint     NOT NULL,
    PRIMARY KEY ("id")
);

CREATE TABLE public."Address"
(
    "id"             bigint       NOT NULL,
    "xid"            varchar(64)  NOT NULL,
    "metadata"       JSON         NULL,
    "createdAt"      timestamp without time zone NOT NULL,
    "updatedAt"      timestamp without time zone NOT NULL,
    "modifiedBy"     JSON         NOT NULL,
    "version"        bigint       NOT NULL DEFAULT 1,
    "customerId"     bigint       NOT NULL,
    "purpose"        smallint     NOT NULL,
    "provinceId"     bigint       NOT NULL,
    "provinceName"   varchar(255) NOT NULL,
    "cityId"         bigint       NOT NULL,
    "cityName"       varchar(255) NOT NULL,
    "districtId"     bigint       NOT NULL,
    "districtName"   varchar(255) NOT NULL,
    "subDistrictId"  bigint       NOT NULL,
    "subDistrictName" varchar(255) NOT NULL,
    PRIMARY KEY ("id")
);


CREATE TABLE public."Verification"
(
    "id"                              bigint       NOT NULL,
    "xid"                             varchar(64)  NOT NULL,
    "metadata"                        JSON         NULL,
    "createdAt"                       timestamp without time zone NOT NULL,
    "updatedAt"                       timestamp without time zone NOT NULL,
    "modifiedBy"                      JSON         NOT NULL,
    "version"                         bigint       NOT NULL DEFAULT 1,
    "customerId"                      bigint       NOT NULL,
    "kycVerifiedStatus"               smallint      NOT NULL,
    "kycVerifiedAt"                   timestamp without time zone NOT NULL,
    "emailVerifiedStatus"             smallint      NOT NULL,
    "emailVerifiedAt"                 timestamp without time zone NOT NULL,
    "dukcapilVerifiedStatus"          smallint      NOT NULL,
    "dukcapilVerifiedAt"              timestamp without time zone NOT NULL,
    "financialTransactionStatus"      smallint     NOT NULL,
    "financialTransactionActivatedAt" timestamp without time zone NOT NULL,
    PRIMARY KEY ("id")
);

CREATE TABLE public."FinancialData"
(
    "id"                              bigint       NOT NULL,
    "xid"                             varchar(64)  NOT NULL,
    "metadata"                        JSON         NULL,
    "createdAt"                       timestamp without time zone NOT NULL,
    "updatedAt"                       timestamp without time zone NOT NULL,
    "modifiedBy"                      JSON         NOT NULL,
    "version"                         bigint       NOT NULL DEFAULT 1,
    "customerId"                      bigint       NOT NULL,
    "mainAccountNumber"               varchar(255) NOT NULL,
    "accountNumber"                   varchar(255) NOT NULL,
    "goldSavingStatus"                smallint      NOT NULL,
    "goldCardApplicationNumber"       varchar(255) NOT NULL,
    "goldCardAccountNumber"           varchar(255) NOT NULL,
    "balance"                         bigint       NOT NULL,
    PRIMARY KEY ("id")
);

CREATE TABLE public."AccessSession"
(
    "id"                              bigint       NOT NULL,
    "xid"                             varchar(64)  NOT NULL,
    "metadata"                        JSON         NULL,
    "createdAt"                       timestamp without time zone NOT NULL,
    "updatedAt"                       timestamp without time zone NOT NULL,
    "modifiedBy"                      JSON         NOT NULL,
    "version"                         bigint       NOT NULL DEFAULT 1,
    "customerId"                      bigint       NOT NULL,
    "expiredAt"                       timestamp without time zone NOT NULL,
    "notificationToken"               TEXT         NOT NULL,
    "notificationProvider"            smallint     NOT NULL,
    PRIMARY KEY ("id")
);

CREATE TABLE public."Credential"
(
    "id"                              bigint       NOT NULL,
    "xid"                             varchar(64)  NOT NULL,
    "metadata"                        JSON         NULL,
    "createdAt"                       timestamp without time zone NOT NULL,
    "updatedAt"                       timestamp without time zone NOT NULL,
    "modifiedBy"                      JSON         NOT NULL,
    "version"                         bigint       NOT NULL DEFAULT 1,
    "customerId"                      bigint       NOT NULL,
    "password"                        varchar(60)  NOT NULL,
    "nextPasswordResetAt"             timestamp without time zone NOT NULL,
    "pin"                             varchar(32)  NOT NULL,
    "pinCif"                          varchar(200) NOT NULL,
    "pinUpdatedAt"                    timestamp without time zone NOT NULL,
    "pinLastAccessAt"                 timestamp without time zone NOT NULL,
    "pinCounter"                      smallint      NOT NULL,
    "pinBlockedStatus"                smallint      NOT NULL,
    "isLocked"                        smallint      NOT NULL,
    "loginFailCount"                  smallint      NOT NULL,
    "wrongPasswordCount"              smallint      NOT NULL,
    "blockedAt"                       timestamp without time zone NOT NULL,
    "blockedUntilAt"                  timestamp without time zone NOT NULL,
    "biometricLogin"                  smallint      NOT NULL,
    "biometricDeviceId"               bigint       NOT NULL,
    PRIMARY KEY ("id")
);