-- Sample Migration Customer
CREATE TABLE public."TransactionFavorite"
(
    "id"              BIGSERIAL PRIMARY KEY NOT NULL,
    "xid"             VARCHAR(64)           NOT NULL,
    "metadata"        JSON NULL,
    "createdAt"       TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    "updatedAt"       TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    "modifiedBy"      JSON                  NOT NULL,
    "version"         BIGINT                NOT NULL DEFAULT 1,
    "customerId"      BIGINT                NOT NULL,
    "type"            VARCHAR(32)           NOT NULL,
    "typeTransaction" VARCHAR(16)           NOT NULL,
    "accountName"     VARCHAR(32)           NOT NULL,
    "accountNumber"   VARCHAR(16)           NOT NULL,
    "bankName"        VARCHAR(64) NULL,
    "bankCode"        VARCHAR(5) NULL,
    "groupMpo"        VARCHAR(32) NULL,
    "serviceCodeMpo"  VARCHAR(32) NULL
);

CREATE INDEX "TransactionFavorite__customerId_idx"
    ON public."TransactionFavorite" ("customerId");

CREATE INDEX "TransactionFavorite__type_idx"
    ON public."TransactionFavorite" ("type");
