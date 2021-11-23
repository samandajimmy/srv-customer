-- Sample Migration User
CREATE TABLE public."User"
(
    "id"         bigint       NOT NULL,
    "username"   varchar(255) NOT NULL,
    "password"   varchar(255) NOT NULL,
    "metadata"   JSON NULL,
    "createdAt"  timestamp without time zone NOT NULL,
    "updatedAt"  timestamp without time zone NOT NULL,
    "modifiedBy" JSON         NOT NULL,
    "version"    bigint       NOT NULL DEFAULT 1,
    PRIMARY KEY ("id")
);