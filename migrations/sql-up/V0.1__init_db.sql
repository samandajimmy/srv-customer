-- Sample Migration Customer
CREATE TABLE public."Customer"
(
    "id"             BIGSERIAL PRIMARY KEY NOT NULL,
    "xid"            varchar(64)  NOT NULL,
    "metadata"       JSON         NULL,
    "createdAt"      timestamp without time zone NOT NULL,
    "updatedAt"      timestamp without time zone NOT NULL,
    "modifiedBy"     JSON         NOT NULL,
    "version"        bigint       NOT NULL DEFAULT 1,
    "fullName"       varchar(255) NOT NULL,
    "phone"          varchar(16)  NOT NULL,
    "email"          varchar(64)  NOT NULL,
    "identityType"   SMALLINT     NULL,
    "identityNumber" varchar(64)  NULL,
    "userRefId"      bigint       NULL,
    "photos"         JSON         NULL,
    "profile"        JSON         NULL,
    "cif"            varchar(16)  NULL,
    "sid"            varchar(16)  NULL,
    "referralCode"   varchar(64)  NULL,
    "status"         smallint     NOT NULL

);

CREATE TABLE public."VerificationOTP"(
    "id"             BIGSERIAL PRIMARY KEY NOT NULL,
    "createdAt"      timestamp without time zone NOT NULL,
    "registrationId" varchar      NOT NULL,
    "phone"          varchar(16)  NOT NULL


);

CREATE TABLE public."Address"
(
    "id"             BIGSERIAL PRIMARY KEY NOT NULL,
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
    "subDistrictName" varchar(255) NOT NULL

);


CREATE TABLE public."Verification"
(
    "id"                              BIGSERIAL PRIMARY KEY       NOT NULL,
    "xid"                             varchar(64)                 NOT NULL,
    "metadata"                        JSON                        NULL,
    "createdAt"                       timestamp without time zone NOT NULL,
    "updatedAt"                       timestamp without time zone NOT NULL,
    "modifiedBy"                      JSON                        NOT NULL,
    "version"                         bigint                      NOT NULL DEFAULT 1,
    "customerId"                      bigint                      NOT NULL,
    "kycVerifiedStatus"               smallint                    NOT NULL,
    "kycVerifiedAt"                   timestamp without time zone NOT NULL,
    "emailVerificationToken"              VARCHAR(128),
    "emailVerifiedStatus"             smallint                    NOT NULL,
    "emailVerifiedAt"                 timestamp without time zone NOT NULL,
    "dukcapilVerifiedStatus"          smallint                    NOT NULL,
    "dukcapilVerifiedAt"              timestamp without time zone NOT NULL,
    "financialTransactionStatus"      smallint                    NOT NULL,
    "financialTransactionActivatedAt" timestamp without time zone NOT NULL
);

CREATE TABLE public."FinancialData"
(
    "id"                              BIGSERIAL PRIMARY KEY NOT NULL,
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
    "balance"                         bigint       NOT NULL
);

CREATE TABLE public."AccessSession"
(
    "id"                              BIGSERIAL PRIMARY KEY NOT NULL,
    "xid"                             varchar(64)  NOT NULL,
    "metadata"                        JSON         NULL,
    "createdAt"                       timestamp without time zone NOT NULL,
    "updatedAt"                       timestamp without time zone NOT NULL,
    "modifiedBy"                      JSON         NOT NULL,
    "version"                         bigint       NOT NULL DEFAULT 1,
    "customerId"                      bigint       NOT NULL,
    "expiredAt"                       timestamp without time zone NOT NULL,
    "notificationToken"               TEXT         NOT NULL,
    "notificationProvider"            smallint     NOT NULL
);

CREATE TABLE public."Credential"
(
    "id"                              BIGSERIAL PRIMARY KEY NOT NULL,
    "xid"                             varchar(64)  NOT NULL,
    "metadata"                        JSON         NULL,
    "createdAt"                       timestamp without time zone NOT NULL,
    "updatedAt"                       timestamp without time zone NOT NULL,
    "modifiedBy"                      JSON         NOT NULL,
    "version"                         bigint       NOT NULL DEFAULT 1,
    "customerId"                      bigint       NOT NULL,
    "password"                        varchar(60)  NOT NULL,
    "nextPasswordResetAt"             timestamp without time zone,
    "pin"                             varchar(32)  NOT NULL,
    "pinCif"                          varchar(200) NOT NULL,
    "pinUpdatedAt"                    timestamp without time zone,
    "pinLastAccessAt"                 timestamp without time zone,
    "pinCounter"                      smallint      NOT NULL,
    "pinBlockedStatus"                smallint      NOT NULL,
    "isLocked"                        smallint      NOT NULL,
    "loginFailCount"                  smallint      NOT NULL,
    "wrongPasswordCount"              smallint      NOT NULL,
    "blockedAt"                       timestamp without time zone,
    "blockedUntilAt"                  timestamp without time zone,
    "biometricLogin"                  smallint      NOT NULL,
    "biometricDeviceId"               varchar(64)   NOT NULL
);

-- TODO: ADD FOREIGN KEY each customerId field --
CREATE TABLE public."OTP"
(
    "id"                              BIGSERIAL PRIMARY KEY           NOT NULL,
    "updatedAt"                       timestamp without time zone     NOT NULL,
    "customerId"                      BIGINT                          NOT NULL,
    "content"                         varchar(6)                      NULL,
    "type"                            VARCHAR(32)                     NULL,
    "data"                            VARCHAR(128)                    NULL,
    "status"                          VARCHAR(16)                     NULL
);
