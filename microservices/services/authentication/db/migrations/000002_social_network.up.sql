CREATE TYPE sn_status AS ENUM ('PENDING', 'ACCEPTED', 'BLOCKED');

CREATE TABLE IF NOT EXISTS "social_network" (
	crux_user_id bigint,
	crux_user_friend_id bigint,
	social_network_status sn_status NOT NULL DEFAULT 'PENDING',
	PRIMARY KEY (crux_user_id, crux_user_friend_id),
	CONSTRAINT fk_social_network_user_id FOREIGN KEY (crux_user_id) REFERENCES crux_user(id) ON UPDATE CASCADE ON DELETE CASCADE,
	CONSTRAINT fk_social_network_crux_user_friend_id FOREIGN KEY (crux_user_id) REFERENCES crux_user(id) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE INDEX idx_social_network_cui_sss ON "social_network" USING btree (crux_user_id, social_network_status);