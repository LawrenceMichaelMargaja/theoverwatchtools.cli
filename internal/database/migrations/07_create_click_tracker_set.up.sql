SET
FOREIGN_KEY_CHECKS = 0;

CREATE TABLE `click_tracker_set`
(
    `id`                        int(11) NOT NULL AUTO_INCREMENT,
    `name`                      varchar(255) NOT NULL,
    `url_name`                  LONGTEXT,
    `analytics_number_of_links` INTEGER      NOT NULL DEFAULT 0,
    `analytics_last_updated_at` INTEGER      NOT NULL DEFAULT 0,
    `last_impression_at`        timestamp NULL DEFAULT NULL,
    `organization_id`           INTEGER      NOT NULL,

    -- Audit fields
    `created_by`                int(11) DEFAULT NULL,
    `last_updated_by`           int(11) DEFAULT NULL,
    `created_at`                timestamp    NOT NULL DEFAULT current_timestamp,
    `last_updated_at`           timestamp NULL DEFAULT NULL ON UPDATE current_timestamp,
    `is_active`                 bool         NOT NULL DEFAULT TRUE,

    CONSTRAINT `click_tracker_set_created_by_ref_id_fk` FOREIGN KEY (`created_by`) REFERENCES `user` (`id`),
    CONSTRAINT `click_tracker_set_last_updated_by_ref_id_fk` FOREIGN KEY (`last_updated_by`) REFERENCES `user` (`id`),

    PRIMARY KEY (`id`),
    UNIQUE KEY `name` (`name`),
    UNIQUE KEY click_tracker_set_unique_name_organization (name, organization_id)
);

INSERT INTO click_tracker_set (
    id, name, url_name, analytics_number_of_links, analytics_last_updated_at, last_impression_at, organization_id, created_by, last_updated_by, created_at, last_updated_at, is_active
) VALUES (
             1, 'Example Set', 'https://example.com', 10, 1693774800, '2023-09-03 10:00:00', 1, 1, 1, '2023-09-01 08:00:00', NULL, TRUE
         );

SET
FOREIGN_KEY_CHECKS = 1;
