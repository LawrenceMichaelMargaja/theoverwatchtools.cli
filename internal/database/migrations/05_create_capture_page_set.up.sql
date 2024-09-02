SET
FOREIGN_KEY_CHECKS = 0;

CREATE TABLE `capture_page_set`
(
    `id`                        int(11) NOT NULL AUTO_INCREMENT,
    `name`                      varchar(255) NOT NULL,
    `url_name`                  LONGTEXT,
    `switch_duration`           INTEGER      NOT NULL,
    `organization_ref_id`       INTEGER,

    -- Updated by workers
    `analytics_number_of_forms` INTEGER      NOT NULL DEFAULT 0,
    `analytics_impressions`     INTEGER      NOT NULL DEFAULT 0,
    `analytics_submissions`     INTEGER      NOT NULL DEFAULT 0,
    `analytics_last_updated_at` INTEGER      NOT NULL DEFAULT 0,

    -- Audit fields
    `created_by`                int(11) DEFAULT NULL,
    `last_updated_by`           int(11) DEFAULT NULL,
    `created_at`                timestamp    NOT NULL DEFAULT current_timestamp,
    `last_updated_at`           timestamp NULL DEFAULT NULL ON UPDATE current_timestamp,
    `is_active`                 bool         NOT NULL DEFAULT TRUE,

    CONSTRAINT `capture_page_set_created_by_ref_id_fk` FOREIGN KEY (`created_by`) REFERENCES `user` (`id`),
    CONSTRAINT `capture_page_set_last_updated_by_ref_id_fk` FOREIGN KEY (`last_updated_by`) REFERENCES `user` (`id`),
    CONSTRAINT `capture_page_set_organization_id_fk` FOREIGN KEY (`organization_ref_id`) REFERENCES `organization` (`id`),

    PRIMARY KEY (`id`),
    UNIQUE KEY `name` (`name`),
    UNIQUE KEY capture_page_set_unique_name_organization (name, organization_ref_id)
);

INSERT INTO `capture_page_set` (
    `name`,
    `url_name`,
    `switch_duration`,
    `organization_ref_id`,
    `analytics_number_of_forms`,
    `analytics_impressions`,
    `analytics_submissions`,
    `analytics_last_updated_at`,
    `created_by`,
    `last_updated_by`,
    `created_at`,
    `last_updated_at`,
    `is_active`
) VALUES
    (
        'Landing Page Set 1',
        'landing-page-set-1',
        10,
        1,
        50,
        100,
        25,
        UNIX_TIMESTAMP(),
        2,
        2,
        NOW(),
        NOW(),
        TRUE
    );

INSERT INTO `capture_page_set` (
    `name`,
    `url_name`,
    `switch_duration`,
    `organization_ref_id`,
    `analytics_number_of_forms`,
    `analytics_impressions`,
    `analytics_submissions`,
    `analytics_last_updated_at`,
    `created_by`,
    `last_updated_by`,
    `created_at`,
    `last_updated_at`,
    `is_active`
) VALUES
    (
        'Product Launch Set',
        'product-launch-set',
        15,
        2,
        75,
        200,
        40,
        UNIX_TIMESTAMP(),
        3,
        3,
        NOW(),
        NOW(),
        TRUE
    );

INSERT INTO `capture_page_set` (
    `name`,
    `url_name`,
    `switch_duration`,
    `organization_ref_id`,
    `analytics_number_of_forms`,
    `analytics_impressions`,
    `analytics_submissions`,
    `analytics_last_updated_at`,
    `created_by`,
    `last_updated_by`,
    `created_at`,
    `last_updated_at`,
    `is_active`
) VALUES
    (
        'Event Promotion Set',
        'event-promotion-set',
        5,
        3,
        30,
        120,
        15,
        UNIX_TIMESTAMP(),
        4,
        4,
        NOW(),
        NOW(),
        FALSE
    );
SET
FOREIGN_KEY_CHECKS = 1;

