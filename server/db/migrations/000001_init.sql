-- +goose Up
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE teams (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE players (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    team_id uuid REFERENCES teams(id) ON DELETE SET NULL,
    display_name text NOT NULL,
    avatar_url text,
    auth_token text NOT NULL UNIQUE,
    auth_token_created_at timestamptz NOT NULL DEFAULT now(),
    qr_token text NOT NULL UNIQUE,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX players_team_id_idx ON players (team_id);
CREATE INDEX players_auth_token_idx ON players (auth_token);
CREATE INDEX players_qr_token_idx ON players (qr_token);

CREATE TABLE open_power_ledger (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    player_id uuid NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    amount integer NOT NULL CHECK (amount <> 0),
    reason text NOT NULL CHECK (
        reason IN (
            'mission',
            'quiz',
            'correction',
            'world_boss',
            'crafting',
            'staff_grant'
        )
    ),
    source_type text NOT NULL,
    source_id uuid,
    metadata jsonb NOT NULL DEFAULT '{}'::jsonb CHECK (jsonb_typeof(metadata) = 'object'),
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX open_power_ledger_player_created_at_idx
    ON open_power_ledger (player_id, created_at DESC);
CREATE INDEX open_power_ledger_source_idx
    ON open_power_ledger (source_type, source_id);

CREATE TABLE reward_grants (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    player_id uuid REFERENCES players(id) ON DELETE CASCADE,
    team_id uuid REFERENCES teams(id) ON DELETE CASCADE,
    source_type text NOT NULL CHECK (
        source_type IN (
            'mission',
            'bingo_line',
            'quiz',
            'world_boss_stage',
            'crafting'
        )
    ),
    source_id uuid NOT NULL,
    reward_type text NOT NULL CHECK (
        reward_type IN (
            'open_power',
            'sitone',
            'item',
            'base_theme'
        )
    ),
    reward_ref_id uuid,
    quantity integer NOT NULL CHECK (quantity > 0),
    created_at timestamptz NOT NULL DEFAULT now(),
    CHECK (
        (player_id IS NOT NULL AND team_id IS NULL)
        OR (player_id IS NULL AND team_id IS NOT NULL)
    )
);

CREATE INDEX reward_grants_player_created_at_idx
    ON reward_grants (player_id, created_at DESC)
    WHERE player_id IS NOT NULL;
CREATE INDEX reward_grants_team_created_at_idx
    ON reward_grants (team_id, created_at DESC)
    WHERE team_id IS NOT NULL;
CREATE INDEX reward_grants_source_idx
    ON reward_grants (source_type, source_id);

CREATE TABLE player_sitones (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    player_id uuid NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    sitone_code text NOT NULL,
    source_type text NOT NULL CHECK (
        source_type IN (
            'mission',
            'quiz',
            'world_boss',
            'crafting',
            'staff_grant'
        )
    ),
    source_id uuid,
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX player_sitones_player_created_at_idx
    ON player_sitones (player_id, created_at DESC);
CREATE INDEX player_sitones_player_code_idx
    ON player_sitones (player_id, sitone_code);

CREATE TABLE player_items (
    player_id uuid NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    item_code text NOT NULL,
    quantity integer NOT NULL CHECK (quantity >= 0),
    updated_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (player_id, item_code)
);

CREATE TABLE player_sitone_loadout_slots (
    player_id uuid NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    slot_index integer NOT NULL CHECK (slot_index BETWEEN 1 AND 5),
    player_sitone_id uuid REFERENCES player_sitones(id) ON DELETE SET NULL,
    updated_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (player_id, slot_index),
    UNIQUE (player_sitone_id)
);

CREATE TABLE crafting_records (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    player_id uuid NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    recipe_code text NOT NULL,
    input_player_sitone_id uuid NOT NULL REFERENCES player_sitones(id) ON DELETE RESTRICT,
    input_item_code text NOT NULL,
    output_player_sitone_id uuid REFERENCES player_sitones(id) ON DELETE SET NULL,
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX crafting_records_player_created_at_idx
    ON crafting_records (player_id, created_at DESC);
CREATE INDEX crafting_records_recipe_code_idx
    ON crafting_records (recipe_code);

CREATE TABLE team_base_theme_unlocks (
    team_id uuid NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    base_theme_code text NOT NULL,
    source_type text NOT NULL,
    source_id uuid,
    unlocked_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (team_id, base_theme_code)
);

CREATE TABLE team_base_theme_selection (
    team_id uuid PRIMARY KEY REFERENCES teams(id) ON DELETE CASCADE,
    base_theme_code text NOT NULL,
    updated_at timestamptz NOT NULL DEFAULT now(),
    FOREIGN KEY (team_id, base_theme_code)
        REFERENCES team_base_theme_unlocks (team_id, base_theme_code)
        ON DELETE CASCADE
);

CREATE TABLE player_mission_progress (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    player_id uuid NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    mission_code text NOT NULL,
    progress_count integer NOT NULL DEFAULT 0 CHECK (progress_count >= 0),
    status text NOT NULL CHECK (status IN ('in_progress', 'completed', 'claimed')),
    completed_at timestamptz,
    claimed_at timestamptz,
    updated_at timestamptz NOT NULL DEFAULT now(),
    UNIQUE (player_id, mission_code)
);

CREATE INDEX player_mission_progress_player_status_idx
    ON player_mission_progress (player_id, status);
CREATE INDEX player_mission_progress_mission_code_idx
    ON player_mission_progress (mission_code);

CREATE TABLE mission_progress_events (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    player_id uuid NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    mission_code text NOT NULL,
    event_type text NOT NULL CHECK (
        event_type IN (
            'flag_submit',
            'staff_scan',
            'system_event',
            'counter_increment'
        )
    ),
    amount integer NOT NULL DEFAULT 1 CHECK (amount > 0),
    actor_player_id uuid REFERENCES players(id) ON DELETE SET NULL,
    metadata jsonb NOT NULL DEFAULT '{}'::jsonb CHECK (jsonb_typeof(metadata) = 'object'),
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX mission_progress_events_player_mission_created_at_idx
    ON mission_progress_events (player_id, mission_code, created_at DESC);
CREATE INDEX mission_progress_events_actor_player_id_idx
    ON mission_progress_events (actor_player_id)
    WHERE actor_player_id IS NOT NULL;

CREATE TABLE bingo_boards (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    category text NOT NULL CHECK (category IN ('daily', 'permanent', 'event_limited')),
    title text NOT NULL,
    date date,
    starts_at timestamptz,
    ends_at timestamptz,
    is_active boolean NOT NULL DEFAULT true,
    CHECK (ends_at IS NULL OR starts_at IS NULL OR ends_at > starts_at)
);

CREATE INDEX bingo_boards_active_category_idx
    ON bingo_boards (is_active, category);
CREATE INDEX bingo_boards_date_idx
    ON bingo_boards (date)
    WHERE date IS NOT NULL;

CREATE TABLE bingo_cells (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    board_id uuid NOT NULL REFERENCES bingo_boards(id) ON DELETE CASCADE,
    mission_code text NOT NULL,
    row_index integer NOT NULL CHECK (row_index BETWEEN 0 AND 4),
    col_index integer NOT NULL CHECK (col_index BETWEEN 0 AND 4),
    UNIQUE (board_id, row_index, col_index),
    UNIQUE (board_id, mission_code)
);

CREATE INDEX bingo_cells_board_id_idx ON bingo_cells (board_id);

CREATE TABLE bingo_line_rewards (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    board_id uuid NOT NULL REFERENCES bingo_boards(id) ON DELETE CASCADE,
    claim_order integer NOT NULL CHECK (claim_order > 0),
    reward_type text NOT NULL CHECK (
        reward_type IN (
            'open_power',
            'sitone',
            'item',
            'base_theme'
        )
    ),
    reward_ref_id uuid,
    quantity integer NOT NULL CHECK (quantity > 0),
    UNIQUE (board_id, claim_order)
);

CREATE INDEX bingo_line_rewards_board_id_idx ON bingo_line_rewards (board_id);

CREATE TABLE player_bingo_line_claims (
    player_id uuid NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    bingo_line_reward_id uuid NOT NULL REFERENCES bingo_line_rewards(id) ON DELETE CASCADE,
    line_key text NOT NULL,
    claimed_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (player_id, bingo_line_reward_id)
);

CREATE INDEX player_bingo_line_claims_player_claimed_at_idx
    ON player_bingo_line_claims (player_id, claimed_at DESC);

CREATE TABLE offline_opponent_snapshots (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    source_player_id uuid NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    display_name text NOT NULL,
    loadout_snapshot jsonb NOT NULL CHECK (jsonb_typeof(loadout_snapshot) = 'object'),
    recent_performance jsonb NOT NULL CHECK (jsonb_typeof(recent_performance) = 'object'),
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX offline_opponent_snapshots_source_player_created_at_idx
    ON offline_opponent_snapshots (source_player_id, created_at DESC);

CREATE TABLE world_bosses (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    code text NOT NULL UNIQUE,
    name text NOT NULL,
    description text NOT NULL,
    status text NOT NULL CHECK (status IN ('active', 'defeated', 'archived')),
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX world_bosses_status_idx ON world_bosses (status);

CREATE TABLE world_boss_stages (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    world_boss_id uuid NOT NULL REFERENCES world_bosses(id) ON DELETE CASCADE,
    stage_index integer NOT NULL CHECK (stage_index > 0),
    name text NOT NULL,
    max_hp integer NOT NULL CHECK (max_hp > 0),
    current_hp integer NOT NULL CHECK (current_hp >= 0),
    challenge_limit_per_player integer NOT NULL CHECK (challenge_limit_per_player >= 0),
    topic_weight_config jsonb NOT NULL DEFAULT '{}'::jsonb CHECK (jsonb_typeof(topic_weight_config) = 'object'),
    status text NOT NULL CHECK (status IN ('locked', 'active', 'defeated')),
    unlocked_at timestamptz,
    defeated_at timestamptz,
    UNIQUE (world_boss_id, stage_index),
    CHECK (current_hp <= max_hp)
);

CREATE INDEX world_boss_stages_boss_status_idx
    ON world_boss_stages (world_boss_id, status);

CREATE TABLE matches (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    mode text NOT NULL CHECK (mode IN ('qr_duel', 'offline_duel', 'world_boss')),
    status text NOT NULL CHECK (status IN ('pairing', 'active', 'completed', 'cancelled')),
    initiator_player_id uuid NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    opponent_player_id uuid REFERENCES players(id) ON DELETE SET NULL,
    opponent_snapshot_id uuid REFERENCES offline_opponent_snapshots(id) ON DELETE SET NULL,
    world_boss_stage_id uuid REFERENCES world_boss_stages(id) ON DELETE SET NULL,
    started_at timestamptz,
    completed_at timestamptz,
    created_at timestamptz NOT NULL DEFAULT now(),
    CHECK (completed_at IS NULL OR started_at IS NULL OR completed_at >= started_at)
);

CREATE INDEX matches_initiator_created_at_idx
    ON matches (initiator_player_id, created_at DESC);
CREATE INDEX matches_opponent_created_at_idx
    ON matches (opponent_player_id, created_at DESC)
    WHERE opponent_player_id IS NOT NULL;
CREATE INDEX matches_world_boss_stage_id_idx
    ON matches (world_boss_stage_id)
    WHERE world_boss_stage_id IS NOT NULL;
CREATE INDEX matches_status_idx ON matches (status);

CREATE TABLE match_players (
    match_id uuid NOT NULL REFERENCES matches(id) ON DELETE CASCADE,
    player_id uuid NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    side text NOT NULL CHECK (side IN ('initiator', 'opponent')),
    score integer NOT NULL DEFAULT 0 CHECK (score >= 0),
    correct_count integer NOT NULL DEFAULT 0 CHECK (correct_count >= 0),
    open_power_earned integer NOT NULL DEFAULT 0 CHECK (open_power_earned >= 0),
    result text CHECK (result IS NULL OR result IN ('win', 'lose', 'draw', 'completed')),
    PRIMARY KEY (match_id, player_id),
    UNIQUE (match_id, side)
);

CREATE INDEX match_players_player_id_idx ON match_players (player_id);

CREATE TABLE match_questions (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    match_id uuid NOT NULL REFERENCES matches(id) ON DELETE CASCADE,
    question_code text NOT NULL,
    question_order integer NOT NULL CHECK (question_order > 0),
    question_color text NOT NULL,
    base_score integer NOT NULL CHECK (base_score >= 0),
    color_bonus_score integer NOT NULL CHECK (color_bonus_score >= 0),
    UNIQUE (match_id, question_order)
);

CREATE INDEX match_questions_match_id_idx ON match_questions (match_id);

CREATE TABLE match_answers (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    match_question_id uuid NOT NULL REFERENCES match_questions(id) ON DELETE CASCADE,
    player_id uuid NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    selected_option_key text,
    is_correct boolean NOT NULL,
    score_awarded integer NOT NULL CHECK (score_awarded >= 0),
    answered_at timestamptz NOT NULL DEFAULT now(),
    UNIQUE (match_question_id, player_id)
);

CREATE INDEX match_answers_player_answered_at_idx
    ON match_answers (player_id, answered_at DESC);

CREATE TABLE match_corrections (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    answer_id uuid NOT NULL UNIQUE REFERENCES match_answers(id) ON DELETE CASCADE,
    reason text NOT NULL CHECK (
        reason IN (
            'misread',
            'unfamiliar',
            'confused_by_options',
            'guessed',
            'understood_after_explanation'
        )
    ),
    open_power_earned integer NOT NULL DEFAULT 0 CHECK (open_power_earned >= 0),
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE match_pairings (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    initiator_player_id uuid NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    target_player_id uuid NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    status text NOT NULL CHECK (status IN ('pending', 'accepted', 'cancelled')),
    match_id uuid REFERENCES matches(id) ON DELETE SET NULL,
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX match_pairings_initiator_created_at_idx
    ON match_pairings (initiator_player_id, created_at DESC);
CREATE INDEX match_pairings_target_created_at_idx
    ON match_pairings (target_player_id, created_at DESC);
CREATE INDEX match_pairings_status_idx ON match_pairings (status);

CREATE TABLE world_boss_attempts (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    stage_id uuid NOT NULL REFERENCES world_boss_stages(id) ON DELETE CASCADE,
    player_id uuid NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    match_id uuid NOT NULL UNIQUE REFERENCES matches(id) ON DELETE CASCADE,
    damage_dealt integer NOT NULL DEFAULT 0 CHECK (damage_dealt >= 0),
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX world_boss_attempts_stage_player_idx
    ON world_boss_attempts (stage_id, player_id);
CREATE INDEX world_boss_attempts_player_created_at_idx
    ON world_boss_attempts (player_id, created_at DESC);

CREATE TABLE world_boss_progress_events (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    stage_id uuid NOT NULL REFERENCES world_boss_stages(id) ON DELETE CASCADE,
    player_id uuid REFERENCES players(id) ON DELETE SET NULL,
    team_id uuid REFERENCES teams(id) ON DELETE SET NULL,
    source_type text NOT NULL CHECK (
        source_type IN (
            'match_answer',
            'correction',
            'mission',
            'activity'
        )
    ),
    source_id uuid,
    damage integer NOT NULL CHECK (damage >= 0),
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX world_boss_progress_events_stage_created_at_idx
    ON world_boss_progress_events (stage_id, created_at DESC);
CREATE INDEX world_boss_progress_events_player_created_at_idx
    ON world_boss_progress_events (player_id, created_at DESC)
    WHERE player_id IS NOT NULL;
CREATE INDEX world_boss_progress_events_source_idx
    ON world_boss_progress_events (source_type, source_id);

-- +goose Down
DROP TABLE IF EXISTS world_boss_progress_events;
DROP TABLE IF EXISTS world_boss_attempts;
DROP TABLE IF EXISTS match_pairings;
DROP TABLE IF EXISTS match_corrections;
DROP TABLE IF EXISTS match_answers;
DROP TABLE IF EXISTS match_questions;
DROP TABLE IF EXISTS match_players;
DROP TABLE IF EXISTS matches;
DROP TABLE IF EXISTS world_boss_stages;
DROP TABLE IF EXISTS world_bosses;
DROP TABLE IF EXISTS offline_opponent_snapshots;
DROP TABLE IF EXISTS player_bingo_line_claims;
DROP TABLE IF EXISTS bingo_line_rewards;
DROP TABLE IF EXISTS bingo_cells;
DROP TABLE IF EXISTS bingo_boards;
DROP TABLE IF EXISTS mission_progress_events;
DROP TABLE IF EXISTS player_mission_progress;
DROP TABLE IF EXISTS team_base_theme_selection;
DROP TABLE IF EXISTS team_base_theme_unlocks;
DROP TABLE IF EXISTS crafting_records;
DROP TABLE IF EXISTS player_sitone_loadout_slots;
DROP TABLE IF EXISTS player_items;
DROP TABLE IF EXISTS player_sitones;
DROP TABLE IF EXISTS reward_grants;
DROP TABLE IF EXISTS open_power_ledger;
DROP TABLE IF EXISTS players;
DROP TABLE IF EXISTS teams;
DROP EXTENSION IF EXISTS pgcrypto;
