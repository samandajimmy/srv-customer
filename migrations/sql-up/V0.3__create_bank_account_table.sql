CREATE TABLE public."BankAccount"
(
    "id"            BIGSERIAL PRIMARY KEY NOT NULL,
    "xid"           VARCHAR(64)           NOT NULL,
    "metadata"      JSON NULL,
    "createdAt"     TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    "updatedAt"     TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    "modifiedBy"    JSON                  NOT NULL,
    "version"       BIGINT                NOT NULL DEFAULT 1,
    "customerId"    BIGINT                NOT NULL,
    "accountName"   VARCHAR               NOT NULL,
    "accountNumber" VARCHAR               NOT NULL,
    "bank"          JSON                  NOT NULL,
    "status"        SMALLINT              NOT NULL
);

CREATE INDEX BankAccount_customerId_idx
    ON public."BankAccount" ("customerId");
